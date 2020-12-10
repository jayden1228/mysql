package mysql

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"

	"log"
	"time"

	"github.com/jinzhu/gorm"
	// 引用数据库驱动初始化
	_ "github.com/jinzhu/gorm/dialects/mysql"
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
		if err := container[id].Close(); err != nil {
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
	// unregister existed one
	UnregisterByKey(key)

	connectionString := config.GenConnectionString()

	db, err := gorm.Open("mysql", connectionString)
	if err != nil {
		log.Println("connect to examples fail, ", connectionString, err)
		panic(err)
	}

	db.LogMode(config.EnableLog)

	db.DB().SetConnMaxLifetime(config.ConnMaxLifetime * time.Second)
	db.DB().SetMaxOpenConns(config.MaxOpenConns)
	db.DB().SetMaxIdleConns(config.MaxIdleConns)

	// 注册创建/更新回调函数，自动插入时间戳
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	// 禁止update/delete传空对象
	db.BlockGlobalUpdate(true)

	// retain the db
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
	if err := db.Close(); err != nil {
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

	gdb, err := gorm.Open("mysql", db)
	if err != nil {
		panic(err)
	}
	gdb.LogMode(true)
	gdb.SingularTable(true)

	container[key] = gdb
	return db, mock
}

func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if scope.HasError() {
		return
	}

	now := time.Now().UTC().Unix()

	if createdAtField, ok := scope.FieldByName("CreatedAt"); ok && createdAtField.IsBlank {
		if err := createdAtField.Set(now); err != nil {
			log.Print(err)
		}
	}

	if updatedAtField, ok := scope.FieldByName("UpdatedAt"); ok && updatedAtField.IsBlank {
		if err := updatedAtField.Set(now); err != nil {
			log.Print(err)
		}
	}
}

func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		now := time.Now().UTC().Unix()
		if err := scope.SetColumn("UpdatedAt", now); err != nil {
			log.Print(err)
		}
	}
}
