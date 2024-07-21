package orm

import (
	"database/sql"
	"log"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewConnectionDB() (*gorm.DB, error) {
	dsn, errorDBURL := config.GetDBURL()
	if errorDBURL != nil {
		log.Fatal("error getting database URL")
		return nil, errorDBURL
	}
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("error initializing database connection")
		return nil, err
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("error initializing database connection gorm")
		return nil, err
	}

	// Migración de esquema (creación de tablas)
	gormDB.AutoMigrate(&model.Articles{})
	gormDB.AutoMigrate(&model.Classifications{})
	gormDB.AutoMigrate(&model.Lineups{})
	gormDB.AutoMigrate(&model.Matches{})
	gormDB.AutoMigrate(&model.Players{})
	gormDB.AutoMigrate(&model.Roles{})
	gormDB.AutoMigrate(&model.Seasons{})
	gormDB.AutoMigrate(&model.TeamsStats{})
	gormDB.AutoMigrate(&model.Teams{})
	gormDB.AutoMigrate(&model.Users{})

	return gormDB, nil
}
