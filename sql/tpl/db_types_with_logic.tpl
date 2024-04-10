package {{.PackageName}}

import (
	"time"

	"gorm.io/gorm"
)

var Models *LogicModels = &LogicModels{}
 
type LogicModels struct {
{{- range .AllTables}}
	{{.TableNameCamel}}Model *{{.TableNameCamel}}Model
{{- end}} 
} 

{{range .AllTables}}
// {{.TableName}} {{.TableComment}} 
type {{.TableNameCamel}} struct { 
{{- range .TableFields}} 
	{{.FieldNameCamel}} {{.FieldType}} {{.FieldTag}} // {{.FieldComment}}
{{- end}}
}

func (tstru *{{.TableNameCamel}}) TableName() string {
	return "{{.TableName}}"
}
{{end}}


func (models *LogicModels) SetGormDriver(db *gorm.DB) { 
{{- range .AllTables}}
	models.{{.TableNameCamel}}Model = &{{.TableNameCamel}}Model{db: db, TableName: "{{.TableName}}"}
{{- end}} 
}