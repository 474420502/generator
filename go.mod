module github.com/474420502/generator

go 1.22.1

require (
	golang.org/x/text v0.14.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.5.6
	gorm.io/gorm v1.25.9
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
)

replace github.com/go-gorm/gorm v1.25.9 => gorm.io/gorm v1.25.9
