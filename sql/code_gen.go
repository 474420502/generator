package sql_generator

import (
	"database/sql"
	"fmt"
	"go/format"
	"os"

	"github.com/474420502/generator/sql/tpl"
)

// genDir+"/types_gen.go"
func GenModel(genFilePath string, db *sql.DB) {
	var tableTypes []*TableType

	tablenames := GetAllTableNames(db)
	for _, testName := range tablenames {
		cols, tname, tcomment := GetColsFromTable(testName, db)
		tableTypes = append(tableTypes, getTableType(cols, tname, tcomment))
	}

	err := createDirIfNotExists(genFilePath)
	if err != nil {
		panic(err)
	}

	typesGenFile, err := os.OpenFile(genFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	// log.Println(typesGenFile)

	typesGenData, err := tpl.ExecuteTemplateBytes("db_types.tpl", map[string]any{
		"AllTables": tableTypes,
	})
	if err != nil {
		panic(err)
	}
	typesGenData, err = format.Source(typesGenData)
	if err != nil {
		panic(err)
	}
	typesGenFile.Write(typesGenData)
}

var typeForMysqlWithNotNull = map[string]string{
	// 整数
	"int":       "int64",
	"integer":   "int64",
	"tinyint":   "int64",
	"smallint":  "int64",
	"mediumint": "int64",
	"bigint":    "int64",
	"year":      "int64",
	"bit":       "int64",

	"int unsigned":       "uint64",
	"integer unsigned":   "uint64",
	"tinyint unsigned":   "uint64",
	"smallint unsigned":  "uint64",
	"mediumint unsigned": "uint64",
	"bigint unsigned":    "uint64",

	// 布尔类型
	"bool": "bool",

	// 字符串
	"enum":       "string",
	"set":        "string",
	"varchar":    "string",
	"char":       "string",
	"tinytext":   "string",
	"mediumtext": "string",
	"text":       "string",
	"longtext":   "string",

	// 二进制
	"binary":     "[]byte",
	"varbinary":  "[]byte",
	"blob":       "[]byte",
	"tinyblob":   "[]byte",
	"mediumblob": "[]byte",
	"longblob":   "[]byte",

	// 日期时间
	"date":      "time.Time",
	"datetime":  "time.Time",
	"timestamp": "time.Time",
	"time":      "time.Time",

	// 浮点数
	"float":   "float64",
	"double":  "float64",
	"decimal": "float64",

	"float unsigned":   "float64",
	"double unsigned":  "float64",
	"decimal unsigned": "float64",
}

func getTableType(cols []*Column, tableName string, tableComment string) *TableType {
	tableNameCamel := toPascalCase(tableName)

	var tfields []TableTypeField
	tt := &TableType{
		TableName:      tableName,
		TableNameCamel: tableNameCamel,
		TableComment:   tableComment,
	}

	fieldstr := ""
	for _, col := range cols {
		ttypeField := TableTypeField{}
		fieldName := toPascalCase(col.Name)
		typeName := getSqlToGoStruct(col)

		tagstr := getTagString(col)

		ttypeField.FieldName = col.Name
		ttypeField.FieldNameCamel = fieldName
		ttypeField.FieldType = typeName
		ttypeField.FieldTag = tagstr
		ttypeField.FieldComment = col.Comment

		fieldColStr := fmt.Sprintf("\n%s %s %s// %s", fieldName, typeName, tagstr, col.Comment)

		fieldstr += fieldColStr

		tfields = append(tfields, ttypeField)
	}

	// err := createFileWithTplIfNotExists(fmt.Sprintf("%s/%s_logic.go", mdir, tableName), func(f io.Writer) error {
	// 	return tpl.ExecuteTemplate(f, "db_logic.tpl", map[string]any{
	// 		"TableName": tableName,
	// 	})
	// })

	// if err != nil {
	// 	panic(err)
	// }

	tt.TableFields = tfields

	return tt
}

func getTagString(col *Column) string {
	return fmt.Sprintf("`db:\"%s\"`", col.Name)
}

func getSqlToGoStruct(col *Column) string {

	if v, ok := typeForMysqlWithNotNull[col.GetType()]; ok {
		if col.NotNull {
			return v
		}
		return "*" + v
	}

	panic(col)
}
