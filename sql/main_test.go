package sql_generator_test

import (
	"testing"

	sql_generator "github.com/474420502/generator/sql"
	// _ "github.com/go-sql-driver/mysql" // MySQL driver
)

func TestCase1(t *testing.T) {
	genModel := sql_generator.NewGenModel().WithOpenSqlDriver("mysql", "")
	genModel.GenWithLogics()
}
