package model

import (
	"gorm.io/gorm"
)

// Migrate 迁移表格
func Migrate(db *gorm.DB) {
	MigrateUserTable(db) // Migrate user table
}
