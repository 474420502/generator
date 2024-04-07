package sql_generator

import (
	"database/sql"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/474420502/generator/sql/tpl"
)

type GenModel struct {
	GenFileDir          string
	TableStructFileName string
	PackageName         string
	LogicDir            *string
	LogicPackageName    *string
	db                  *sql.DB
}

func NewGenModel() *GenModel {
	return &GenModel{
		GenFileDir:          "model",
		TableStructFileName: "types_gen.go",
		PackageName:         "model",
	}
}

func (gen *GenModel) WithGenFileDir(GenFileDir string) *GenModel {
	absPath, err := filepath.Abs(GenFileDir)
	if err != nil {
		panic(err)
	}
	gen.GenFileDir = absPath
	return gen
}

func (gen *GenModel) WithTableStructFileName(TableStructFileName string) *GenModel {
	gen.TableStructFileName = TableStructFileName
	return gen
}

func (gen *GenModel) WithPackageName(PackageName string) *GenModel {
	gen.PackageName = PackageName
	return gen
}

func (gen *GenModel) WithLoigcDir(LoigcDir string) *GenModel {
	absPath, err := filepath.Abs(LoigcDir)
	if err != nil {
		panic(err)
	}
	gen.LogicDir = &absPath
	return gen
}

func (gen *GenModel) WithLoigcPackName(LoigcPackName string) *GenModel {
	gen.LogicPackageName = &LoigcPackName
	return gen
}

func (gen *GenModel) WithOpenSqlDriver(driverName string, dataSourceName string) *GenModel {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
	gen.db = db
	return gen
}

func (gen *GenModel) WithSqlDb(db *sql.DB) *GenModel {
	gen.db = db
	return gen
}

func (gen *GenModel) Gen() {
	if gen.db == nil {
		panic("WithSqlDb | WithOpenSqlDriver must be called at least once")
	}

	var tableTypes []*TableType

	tablenames := GetAllTableNames(gen.db)
	for _, testName := range tablenames {
		cols, tname, tcomment := GetColsFromTable(testName, gen.db)
		tableTypes = append(tableTypes, getTableType(cols, tname, tcomment))
	}
	genFilePath := strings.TrimRight(gen.GenFileDir, "/") + "/" + gen.TableStructFileName
	tpl.ExecuteTemplateWithCreateGoFile(genFilePath, "db_types.tpl", map[string]any{
		"AllTables":   tableTypes,
		"PackageName": gen.PackageName,
	})
}

func (gen *GenModel) GenWithLogics() {
	if gen.db == nil {
		panic("WithSqlDb | WithOpenSqlDriver must be called at least once")
	}

	var tableTypes []*TableType

	tablenames := GetAllTableNames(gen.db)
	for _, testName := range tablenames {
		cols, tname, tcomment := GetColsFromTable(testName, gen.db)
		tableTypes = append(tableTypes, getTableType(cols, tname, tcomment))
	}

	genFilePath := strings.TrimRight(gen.GenFileDir, "/") + "/" + gen.TableStructFileName
	tpl.ExecuteTemplateWithCreateGoFile(genFilePath, "db_types.tpl", map[string]any{
		"AllTables":   tableTypes,
		"PackageName": gen.PackageName,
	})
	var logicDir, logicPackageName string
	if gen.LogicDir != nil {
		logicDir = *gen.LogicDir
	} else {
		logicDir = gen.GenFileDir
	}

	if gen.LogicPackageName != nil {
		logicPackageName = *gen.LogicPackageName
	} else {
		logicPackageName = gen.PackageName
	}

	for _, tt := range tableTypes {
		err := createFileWithTplIfNotExists(fmt.Sprintf("%s/%s_logic.go", logicDir, tt.TableName), func(f io.Writer) error {
			return tpl.ExecuteTemplate(f, "db_logic.tpl", map[string]any{
				"TableNameCamel": tt.TableNameCamel,
				"PackageName":    logicPackageName,
			})
		})

		if err != nil {
			panic(err)
		}
	}
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
