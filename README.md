# mysql

> mysql package

## Install

```shell
go get gitlab.com/makeblock-go/mysql
```

## Usage

### DB

See more at [mysql_test.go](./examples/mysql/mysql_test.go) and [gorm](https://gorm.io/docs/).

```go
package examples

import (
	"fmt"

	"gitlab.com/makeblock-go/mysql"
)

type Product struct {
	// mysql base model
	mysql.Model
	Code  string
	Price uint
}

func ExampleRegister() {
	cnf := mysql.NewConfig("root",
		"123456",
		"127.0.0.1",
		"3306",
		"test",
		"utf8mb4",
		true,
		mysql.WithMaxOpenConns(1500))
	// register default db
	mysql.Register(cnf)

	cnf2 := mysql.NewConfig("root",
		"123456",
		"127.0.0.1",
		"3306",
		"test2",
		"utf8mb4",
		true)
	// register another db with key "test2"
	mysql.RegisterByKey(cnf2, "test2")

	defer mysql.Close()

	// default db

	db := mysql.GetDB()

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "L1212", Price: 1000})

	// Read
	var product Product
	db.First(&product, 1)                   // find product with id 1
	db.First(&product, "code = ?", "L1212") // find product with code l1212

	// Update - update product's price to 2000
	db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	db.Delete(&product)

	fmt.Println("default db finish")

	// db with key "test2"

	db2 := mysql.GetDBByKey("test2")

	// Migrate the schema
	db2.AutoMigrate(&Product{})

	// Create
	db2.Create(&Product{Code: "L1212", Price: 1000})

	// Read
	db2.First(&product, 1)                   // find product with id 1
	db2.First(&product, "code = ?", "L1212") // find product with code l1212

	// Update - update product's price to 2000
	db2.Model(&product).Update("Price", 2000)

	// Delete - delete product
	db2.Delete(&product)

	fmt.Println("db2 finish")

	// Output: default db finish
	// db2 finish
}
```

### MockDB

mock usage see more at <https://github.com/DATA-DOG/go-sqlmock/blob/master/examples/orders/orders_test.go>

```go
_, mock := mysql.MockDB()
```