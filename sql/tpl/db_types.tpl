package {{.PackageName}}

import (
	"time"
)

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