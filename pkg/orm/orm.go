package orm

import (
	"database/sql"
	"fmt"
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
	initErr  error
)

// GetDBInstance returns a singleton *gorm.DB instance.
func GetDBInstance() (*gorm.DB, error) {
	once.Do(func() {
		initErr = initializeDatabase()
	})
	return instance, initErr
}

func initializeDatabase() error {
	dsn, err := config.GetDBURL()
	if err != nil {
		return fmt.Errorf("error getting database URL: %w", err)
	}

	sqlDB, err := openSQLConnection(dsn)
	if err != nil {
		return err
	}

	if err := setSessionTimezone(sqlDB); err != nil {
		return err
	}

	gormDB, err := openGormConnection(sqlDB)
	if err != nil {
		return err
	}
	instance = gormDB

	configureConnectionPool(sqlDB)

	if err := runMigrations(gormDB); err != nil {
		return err
	}

	return nil
}

func openSQLConnection(dsn string) (*sql.DB, error) {
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening SQL connection: %w", err)
	}
	return sqlDB, nil
}

func setSessionTimezone(sqlDB *sql.DB) error {
	if _, err := sqlDB.Exec("SET time_zone = '+00:00'"); err != nil {
		return fmt.Errorf("error setting MySQL timezone: %w", err)
	}
	return nil
}

func openGormConnection(sqlDB *sql.DB) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger: NewContextAwareGormLogger(),
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("error opening GORM connection: %w", err)
	}
	return gormDB, nil
}

func configureConnectionPool(sqlDB *sql.DB) {
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func runMigrations(db *gorm.DB) error {
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return fmt.Errorf("error disabling foreign key checks: %w", err)
	}

	if err := db.AutoMigrate(
		&model.Article{},
		&model.Role{},
		&model.Season{},
		&model.User{},
		&model.Player{},
		&model.Team{},
		&model.Match{},
		&model.Lineup{},
		&model.TeamStat{},
		&model.PlayerTeam{},
		&model.PlayerStat{},
	); err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		return fmt.Errorf("error re-enabling foreign key checks: %w", err)
	}

	return nil
}
