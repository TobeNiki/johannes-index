package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Database struct {
	DB *gorm.DB
}

func New() *Database {
	config := LoadConfig()
	DBMS := config.DBMS
	USER := config.USER
	PASS := config.PASS
	PROTOCOL := config.PROTOCOL
	DBNAME := config.DBNAME
	PARAME := config.PARAME
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?" + PARAME
	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		panic(err.Error())
	}
	db.SingularTable(true)

	return &Database{
		DB: db,
	}
}
