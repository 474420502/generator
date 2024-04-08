package {{.PackageName}}

import (
    "github.com/jmoiron/sqlx"
)

type {{.TableNameCamel}}Model struct {
    // fields ...
    db *sqlx.DB
}