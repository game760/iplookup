package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL驱动
	"ip-query-system/internal/config"
)

// DB 数据库连接实例
type DB struct {
	*sql.DB
}

// Init 初始化数据库连接（仅支持MySQL）
func Init(cfg *config.Config) (*DB, error) {
	// 构建MySQL连接字符串
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("连接MySQL失败: %v", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("MySQL ping失败: %v", err)
	}

	return &DB{db}, nil
}

// Close 关闭数据库连接
func Close(db *DB) error {
	if db == nil {
		return nil
	}
	return db.Close()
}