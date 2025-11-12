package database

import (
	"database/sql"
	"fmt"
	"iplookup/iplookup_go/internal/config"

	_ "github.com/go-sql-driver/mysql" // MySQL驱动
)

// DB 数据库连接包装
type DB struct {
	*sql.DB
}

// Init 初始化数据库连接
func Init(cfg *config.Config) (*DB, error) {
	// 构建DSN
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	// 打开连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	db.SetMaxOpenConns(100)    // 最大打开连接数
	db.SetMaxIdleConns(20)     // 最大空闲连接数
	db.SetConnMaxLifetime(300) // 连接最大存活时间（秒）

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// Close 关闭数据库连接
func Close(db *DB) error {
	return db.Close()
}