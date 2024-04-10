package sql_generator

import (
	"log"
	"testing"

	// sql_generator "github.com/474420502/generator/sql"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

func TestCase1(t *testing.T) {

	genModel := NewGenModel().WithOpenSqlDriver("mysql", " ")
	genModel.GenWithLogics()

	log.Println()
}
