package migrate

import (
	"gorm.io/gorm"

	"userservice/config"
)

// Migrate Migrate all table
func Migrate(db *gorm.DB, conf *config.Config) {
	MigrateUserTable(db, conf) // Migrate user table
}
