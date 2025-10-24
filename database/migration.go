package database

import (
	"github.com/ranggakrisnaa/go-fiber-starter/database/entities"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&entities.User{},
		&entities.RefreshToken{},
		&entities.Role{},
		&entities.UserRole{},
	); err != nil {
		return err
	}

	return nil
}
