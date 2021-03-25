package gcode

import (

	"os"

	"strings"
	"rest/pkg/tools"
	"rest/pkg/gmysql"
	"text/template"
)

type MakeCodingRequest struct {
	ServerName           string
	ModuleName           string
	TableName            string
	Name                 string
	DatabaseName         string

}


//MakeCoding 代码编写生成，暂时放在这里
func  MakeCoding( in MakeCodingRequest)  {


	//sweaters := Inventory{"BasInfo", "Item", "BasInfoClassItemT","basInfo", "item"}

	sweaters := Inventory{in.Name,in.ServerName, in.ModuleName,"" ,"",in.TableName,nil,""}

	moduleNameSnake := tools.SnakeString(sweaters.ModuleName)
	serverNameSnake := tools.SnakeString(sweaters.ServerName)

	sweaters.ModelName = tools.SnakeToBigCamel(sweaters.TableName)
	sweaters.Url = strings.Replace(serverNameSnake+"/"+moduleNameSnake,"_","/",-1)
	//SELECT COLUMN_NAME as column_name, column_comment FROM INFORMATION_SCHEMA. COLUMNS WHERE table_name = 'bas_mq_dlx' AND table_schema = 'wos_common'

	var result []SqlResult
	// Raw SQL
	gmysql.DB.Raw("SELECT COLUMN_NAME as column_name,data_type, column_comment FROM INFORMATION_SCHEMA. COLUMNS WHERE table_name = ? AND table_schema = ?", sweaters.TableName,in.DatabaseName).Scan(&result)
	i := 1
	for k,v:= range result {
		result[k].ColumnNameCamel = tools.SnakeToCamel(v.ColumnName)
		result[k].ColumnNameBigCamel = tools.SnakeToBigCamel(v.ColumnName)
		result[k].Id = i
		i++

		switch v.DataType {
		case "int","tinyint","smallint","mediumint":
			result[k].StuctTypeName = "int"
			result[k].TypeName = "int32"
		case "varchar","char","text","tinytext","mediumtext","longtext":
			result[k].StuctTypeName = "string"
			result[k].TypeName = "string"
		case "double":
			result[k].StuctTypeName = "float32"
			result[k].TypeName = "float"
		case "bool":
			result[k].StuctTypeName = "bool"
			result[k].TypeName = "bool"
		case "float":
			result[k].StuctTypeName = "float64"
			result[k].TypeName = "float"
		case "timestamp","date","datetime","time":
			result[k].StuctTypeName = "time.Time"
			result[k].TypeName = "string"
		default:
			result[k].StuctTypeName = "string"
			result[k].TypeName = "string"
			//
		}
	}

	sweaters.Columns = result
	sweaters.DbName  = tools.SnakeToBigCamel(strings.Replace(in.DatabaseName,"wos_","",-1))
	serverNameCamel := tools.SnakeToCamel(serverNameSnake)

	MakeFile(sweaters,"./pkg/gcode/model.go.tpl","./model/"+sweaters.TableName+".go-")
	MakeFile(sweaters,"./pkg/gcode/proto.go.tpl","./proto/"+serverNameCamel+".proto-")
	MakeFile(sweaters,"./pkg/gcode/server.go.tpl","./service/"+serverNameCamel+"Service.go-")

}



type SqlResult struct {
	Id int
	DataType string
	ColumnName string
	ColumnComment string
	TypeName string
	ColumnNameCamel string
	ColumnNameBigCamel string
	StuctTypeName string

}

type Inventory struct {
	Name string
	ServerName string
	ModuleName string
	ModelName string
	Url string
	TableName string
	Columns []SqlResult
	DbName string
}

//MakeFile 参数 模板路径 输出文件
func MakeFile(sweaters Inventory,tmplFile string,outFile string) {

	tmpl, _ := template.ParseFiles(tmplFile)

	_ = tmpl.Execute(os.Stdout, sweaters)

	f, _ := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()

	// 渲染并写入文件
	_ = tmpl.Execute(f, sweaters)

}
