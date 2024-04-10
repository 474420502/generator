package {{.PackageName}}

import (
    "gorm.io/gorm"
)

type {{.TableNameCamel}}Model struct {
    // fields ...
    db *gorm.DB
    TableName string // 表名
}