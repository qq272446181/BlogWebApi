package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	cfg := GetConfig()

	var err error

	switch cfg.Database.Driver {
	case "sqlite":
		// 使用 modernc.org/sqlite 驱动
		DB, err = gorm.Open(sqlite.Open(cfg.Database.Name+"?_pragma=foreign_keys(1)"), &gorm.Config{})
	default:
		panic("不支持的数据库驱动")
	}

	if err != nil {
		panic("failed to connect database:" + err.Error())
	}
}
