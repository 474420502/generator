package sql_generator

import "database/sql"

type ColumnFieldType int

const (
	COLUMN_UNKNOWN ColumnFieldType = iota

	COLUMN_Name
	COLUMN_Type
	COLUMN_DefaultValue
	COLUMN_Length
	COLUMN_Decimal
	COLUMN_Unsigned
	COLUMN_NotNull
	COLUMN_AutoIncrement
	COLUMN_Comment

	COLUMN_IndexType
)

type Column struct {
	Name          string
	Type          string
	DefaultValue  *string
	Length        int
	Decimal       int
	Unsigned      bool
	NotNull       bool
	AutoIncrement bool
	Comment       string

	IsUndefine bool

	IndexType string
}

func (col *Column) GetType() string {
	content := col.Type
	if col.Unsigned {
		return content + " unsigned"
	}
	return content
}

func GetColsFromTable(tname string, db *sql.DB) (result []*Column, tableName, tableComment string) {

	var a, ddl string
	err := db.QueryRow("SHOW CREATE TABLE "+tname).Scan(&a, &ddl)
	// log.Println(ddl)
	if err != nil {
		panic(err)
	}

	return ParserDDL(ddl)
}
