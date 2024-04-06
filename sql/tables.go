package sql_generator

import (
	"database/sql"
)

// TableType 表示一个表
type TableType struct {
	TableName      string
	TableNameCamel string
	TableComment   string
	TableFields    []TableTypeField
}

// Field 表示一个字段
type TableTypeField struct {
	FieldName      string
	FieldNameCamel string
	FieldType      string
	FieldTag       string
	FieldComment   string
}

type TableNameComment struct {
	Name    string
	GoName  string
	Comment string
}

type TMCS []TableNameComment

func (u TMCS) Len() int {
	return len(u)
}

func (u TMCS) Less(i, j int) bool {
	return u[i].Name < u[j].Name
}

func (u TMCS) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func GetAllTableNames(db *sql.DB) []string {

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		panic(err)
	}

	var tableNames []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			panic(err)
		}
		tableNames = append(tableNames, tableName)
	}

	return tableNames
}
