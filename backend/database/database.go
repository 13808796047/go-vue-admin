package database

import (
	"os"
	"path/filepath"

	"go-vue-admin/config"
	"go-vue-admin/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg *config.Config) error {
	if dir := filepath.Dir(cfg.DBPath); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return err
	}
	DB = db

	// 自动迁移
	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Menu{},
		&model.API{},
		&model.UserRoles{},
		&model.RoleMenus{},
		&model.RoleApis{},
	); err != nil {
		return err
	}

	seed(db)
	return nil
}
