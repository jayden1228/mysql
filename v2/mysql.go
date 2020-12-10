package mysql

import (
	"database/sql"
	"os"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"log"
	"time"

	"gorm.io/gorm"
	// 引用数据库驱动初始化
	"gorm.io/driver/mysql"
)

var container = make(map[string]*gorm.DB)

const defaultKey = "default"

// GetDB get default db
func GetDB() *gorm.DB {
	return GetDBByKey(defaultKey)
}

// GetDBByKey get by with key
func GetDBByKey(key string) *gorm.DB {
	return container[key]
}

// Close closes current db connection
func Close() {
	for id := range container {
		sqlDB, _ := container[id].DB()
		if err := sqlDB.Close(); err != nil {
			log.Print(err)
		}
	}
}

// Register register default examples
func Register(config *Config) {
	RegisterByKey(config, defaultKey)
}

// RegisterByKey register examples by key
func RegisterByKey(config *Config, key string) {

	UnregisterByKey(key)

	connectionString := config.GenConnectionString()

	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      config.LogLevel,
			Colorful:      false,
		},
	)

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger,
	})

	if err != nil {
		log.Println("connect to database fail, ", connectionString, err)
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Println("get database pool fail,  ", err)
		panic(err)
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime * time.Second)

	container[key] = db
}

// Unregister unregister default examples
func Unregister() {
	UnregisterByKey(defaultKey)
}

// UnregisterByKey unregister examples by key
func UnregisterByKey(key string) {
	db := container[key]
	if db == nil {
		return
	}

	sqlDB, _ := db.DB()

	if err := sqlDB.Close(); err != nil {
		log.Print(err)
	}
	delete(container, key)
}

// MockDB mock default DB
func MockDB() (*sql.DB, sqlmock.Sqlmock) {
	return MockDBByKey(defaultKey)
}

// MockDBByKey mock DB by key
func MockDBByKey(key string) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	dialer := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	gdb, err := gorm.Open(dialer, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Println("connect to database fail, ", err)
		panic(err)
	}

	container[key] = gdb
	return db, mock
}

