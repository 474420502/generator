package {{.}}

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
{{end}}