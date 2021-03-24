package model

import (
	"time"
)

//{{.Name}}
type {{.ModelName}} struct {
    {{range .Columns}}//{{.ColumnComment}}
    {{.ColumnNameBigCamel}} {{.StuctTypeName}}  `gorm:"column:{{.ColumnName}}"`
    {{end}}
}

//{{.Name}}表名
func ({{.ModelName}}) TableName() string {
	return "{{.TableName}}"
}
