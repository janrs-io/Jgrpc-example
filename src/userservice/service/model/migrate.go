package model

import (
	"gorm.io/gorm"
)

// Migrate 迁移表格
func Migrate(mysqlDB *gorm.DB) {
	MigrateUserTable(mysqlDB) // Migrate user table
}
