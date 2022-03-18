package common

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dbDSN    = "host=127.0.0.1,user=postgres,password=realibox2021_postgres,dbname=new_application,port=54321,sslmode=disable,TimeZone=Asia/Shanghai"
	dbDriver = "postgres"
)

var db *gorm.DB

// ConnectDB return orm.DB
func ConnectDB() (*gorm.DB, error) {
	var err error
	if db == nil {
		db, err = gorm.Open(postgres.Open(dbDSN), &gorm.Config{})
		if err != nil {
			return nil, err
		}
	}

	return db, err
}

// DisconnectDB ...
func DisconnectDB() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		err = sqlDB.Close()
		db = nil
		return err
	}
	return nil
}

// NewDBEngine create db engine
func NewDBEngine() *gorm.DB {
	db, err := ConnectDB()

	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()
	return db
}
