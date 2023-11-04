package orm

import (
	"database/sql"
	"log"
	"os"

	user "github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewConnectionDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
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
	gormDB.AutoMigrate(&user.Articles{})
	gormDB.AutoMigrate(&user.Classifications{})
	gormDB.AutoMigrate(&user.Lineups{})
	gormDB.AutoMigrate(&user.Matches{})
	gormDB.AutoMigrate(&user.Players{})
	gormDB.AutoMigrate(&user.Roles{})
	gormDB.AutoMigrate(&user.Seasons{})
	gormDB.AutoMigrate(&user.TeamsStats{})
	gormDB.AutoMigrate(&user.Teams{})
	gormDB.AutoMigrate(&user.Users{})

	return gormDB, nil
}
