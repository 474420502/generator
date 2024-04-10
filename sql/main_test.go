package sql_generator

import (
	"log"
	"testing"

	// sql_generator "github.com/474420502/generator/sql"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

func TestCase1(t *testing.T) {

	genModel := NewGenModel().WithOpenSqlDriver("mysql", "php:aFk3i4Dj#76!4sd@tcp(47.243.100.6:3306)/zunxinfinance?charset=UTF8MB4&timeout=10s")
	genModel.GenWithLogics()

	log.Println()
}
