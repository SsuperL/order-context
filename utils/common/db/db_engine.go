package db

import (
	"log"
	"order-context/acl/adapters/pl"
	"order-context/utils/common"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db       *gorm.DB
	dbDriver string
)

// ConnectDB return orm.DB
func ConnectDB() (*gorm.DB, error) {
	var err error
	config := common.FileConfig
	if err != nil {
		return nil, err
	}

	if driver := os.Getenv("DRIVER"); driver == "" {
		dbDriver = config.DB.Driver
	} else {
		dbDriver = driver
	}

	if db == nil {
		if dbDriver == "postgres" {
			db, err = gorm.Open(postgres.Open(config.DB.DSN), &gorm.Config{})
		} else if dbDriver == "sqlite" {
			db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
		}

		if err != nil {
			return nil, err
		}
		log.Println("Connect to db successfully.")
	}

	return db, err
}

// InitTables init tables
func InitTables(db *gorm.DB) {
	order := pl.Order{}
	invoice := pl.Invoice{}
	tables := map[string]interface{}{
		order.TableName():   order,
		invoice.TableName(): invoice,
	}
	for k, v := range tables {
		if !db.Migrator().HasTable(k) {
			err := db.Migrator().CreateTable(v)
			if err != nil {
				log.Fatalf("init table %s error: %v", k, err)
			}
		}
	}
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
	return db
}
