package orm

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

func GetDBInstance() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		dsn, errorDBURL := config.GetDBURL()
		if errorDBURL != nil {
			log.Fatal("error getting database URL")
			return
		}
		sqlDB, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal("error initializing database connection")
			return
		}
		instance, err = gorm.Open(mysql.New(mysql.Config{
			Conn: sqlDB,
		}), &gorm.Config{})
		if err != nil {
			log.Fatal("error initializing database connection gorm")
			return
		}

		log.Println("Database connection established")

		// Set connection pool parameters
		sqlDB.SetMaxOpenConns(10)               // Maximum number of open connections
		sqlDB.SetMaxIdleConns(5)                // Maximum number of idle connections
		sqlDB.SetConnMaxLifetime(1 * time.Hour) // Maximum connection lifetime

		// Migración de esquema (creación de tablas)
		instance.AutoMigrate(&model.Article{})
		instance.AutoMigrate(&model.Lineup{})
		instance.AutoMigrate(&model.Match{})
		instance.AutoMigrate(&model.Player{})
		instance.AutoMigrate(&model.Role{})
		instance.AutoMigrate(&model.Season{})
		instance.AutoMigrate(&model.TeamStat{})
		instance.AutoMigrate(&model.Team{})
		instance.AutoMigrate(&model.User{})

	})

	return instance, err
}
