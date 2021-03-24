package gmysql

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"rest/pkg/tools"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	userName = "root"
	password = "root"
	ip = "127.0.0.1"
	port = "3306"
	dbName = "test"
)
var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			userName,
			password,
			ip+":"+port,
			dbName))
	if err != nil {
		panic(err)
	}
	DB.LogMode(true)
}



//StringMapByQuery 查询sql后返回 []map[string]interface{} 字符串map数组 (columnsType == string 二进制值转字符串)
func MapByQuery(db *sql.DB, sqlInfo string, columnsType string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Query(sqlInfo, args...)
	if err != nil {
		return nil, err
	}
	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength) //临时存储每行数据
	for index, _ := range cache {              //为每一列初始化一个指针
		var a interface{}
		cache[index] = &a
	}
	var list []map[string]interface{} //返回的切片
	for rows.Next() {
		_ = rows.Scan(cache...)
		item := make(map[string]interface{})
		for i, data := range cache {
			//值处理
			if columnsType == "string" {
				x := *data.(*interface{})
				//item[columns[i]] = string(x.([]byte))//转换字符串类型
				item[columns[i]] = tools.Strval(x)
			} else {
				item[columns[i]] = *data.(*interface{}) //取实际类型
			}
		}
		list = append(list, item)
	}
	_ = rows.Close()
	return list, nil
}

//ReplaceInsertManyCheck 数据批量插入拆分(避免一次插入太多sql过长)
func ReplaceInsertManySplit(db *sql.DB, table string, arr []map[string]interface{}) error {
	quantity := int64(5000) //拆分批量插入，最高5000条
	max := int64(len(arr))
	var segmens = make([][]map[string]interface{}, 0)
	num := (max / quantity) + 1
	end := int64(0)
	for i := int64(1); i <= num; i++ {
		qu := i * quantity
		if i != num {
			segmens = append(segmens, arr[i-1+end:qu])
		} else {
			segmens = append(segmens, arr[i-1+end:])
		}
		end = qu - i
	}
	for _, v := range segmens {
		err := ReplaceInsertMany(db, table, v)
		if err != nil {
			return err
		}
	}
	return nil
}

//ReplaceInsertMany 批量 插入 替换 方法 (如果需要替换功能，一定在表中添加主键或者唯一索引，当数据重叠则替换 已有的)
func ReplaceInsertMany(db *sql.DB, table string, args []map[string]interface{}) error {
	//获取表结构信息
	var columns []string
	columnsSet, _ := MapByQuery(db, "select column_name,data_type from information_schema.columns  where table_name='"+table+"' ", "string")
	for _, vSet := range columnsSet {
		columns = append(columns, tools.Strval(vSet["column_name"]))
	}

	OutColumnMap := make(map[string]string, 0)
	oneOutMap := make(map[string]string, 0)

	for _, v := range columns {
		for _, vSet := range columnsSet {
			if vSet["column_name"] == v {
				switch tools.Strval(vSet["data_type"]) {
				case "decimal", "tinyint", "int", "smallint", "mediumint", "bigint", "float", "double":
					oneOutMap[v] = "0"
				default:
					oneOutMap[v] = ""
				}
				OutColumnMap[v] = tools.Strval(vSet["data_type"])
			}
		}
	}

	insertWhere := ""
	for _, v := range columns {
		insertWhere += "`" + v + "`,"
	}
	insertWhere = strings.TrimRight(insertWhere, ",")

	insert := " REPLACE INTO " + table + " ( " + insertWhere + ") Value "
	var buffer bytes.Buffer
	buffer.WriteString(insert)

	for _, v := range args {
		for kr, vr := range v {
			if typeKey, ok := OutColumnMap[kr]; ok {
				switch typeKey {
				case "decimal", "tinyint", "int", "smallint", "mediumint", "bigint", "float", "double":
					value := tools.Strval(vr)
					if value == "" {
						value = "0"
					}
					oneOutMap[kr] = value
				default:
					oneOutMap[kr] = tools.Strval(vr)
				}
			} else if typeKey, ok := OutColumnMap[tools.CamelToSnakeString(kr)]; ok {
				switch typeKey {
				case "decimal", "tinyint", "int", "smallint", "mediumint", "bigint", "float", "double":
					value := tools.Strval(vr)
					if value == "" {
						value = "0"
					}
					oneOutMap[tools.CamelToSnakeString(kr)] = value
				default:
					oneOutMap[tools.CamelToSnakeString(kr)] = tools.Strval(vr)
				}
			}

		}

		//拼接
		insertOne := "("
		for _, vo := range columns {
			if val, ok := oneOutMap[vo]; ok {
				insertOne += "'" + val + "',"
			}
		}
		insertOne = strings.TrimRight(insertOne, ",")
		insertOne += "),"
		buffer.WriteString(insertOne)
	}

	insertSql := buffer.String()
	insertSql = strings.TrimRight(insertSql, ",")

	fmt.Println(insertSql)
	_, insertErr := db.Exec(insertSql)

	if insertErr != nil {
		return insertErr
	}

	return nil
}
