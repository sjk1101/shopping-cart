package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbInstance *gorm.DB

func Init() error {
	if err := IntMySql(); err != nil {
		return err
	}
	return nil
}

func IntMySql() error {

	address := "127.0.0.1:4000"
	userName := "root"
	password := "abc123"
	database := "shopping_cart"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=UTC",
		userName,
		password,
		address,
		database)

	conn, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return fmt.Errorf("connect to database failed,err: %v", err)
	}

	dbInstance = conn
	return nil
}

func Session() *gorm.DB {
	return dbInstance.Session(&gorm.Session{})
}
