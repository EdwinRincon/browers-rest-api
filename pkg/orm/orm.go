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
)

// GetDBInstance returns a singleton *gorm.DB instance. It initializes the connection only once.
// Returns an error if the connection or migration fails.
func GetDBInstance() (*gorm.DB, error) {
	onceErr := func() error {
		var initErr error
		once.Do(func() {
			dsn, errorDBURL := config.GetDBURL()
			if errorDBURL != nil {
				initErr = fmt.Errorf("error getting database URL: %w", errorDBURL)
				return
			}
			sqlDB, openErr := sql.Open("mysql", dsn)
			if openErr != nil {
				initErr = fmt.Errorf("error initializing database connection: %w", openErr)
				return
			}

			// ensure MySQL session timezone is UTC:
			if _, err := sqlDB.Exec("SET time_zone = '+00:00'"); err != nil {
				initErr = fmt.Errorf("error setting MySQL timezone: %w", err)
				return
			}

			// Initialize GORM with custom logger
			gormConfig := &gorm.Config{
				Logger: NewContextAwareGormLogger(),
			}

			instance, initErr = gorm.Open(mysql.New(mysql.Config{
				Conn: sqlDB,
			}), gormConfig)
			if initErr != nil {
				initErr = fmt.Errorf("error initializing database connection gorm: %w", initErr)
				return
			}

			// Set connection pool parameters
			sqlDB.SetMaxOpenConns(10)               // Maximum number of open connections
			sqlDB.SetMaxIdleConns(5)                // Maximum number of idle connections
			sqlDB.SetConnMaxLifetime(1 * time.Hour) // Maximum connection lifetime

			// Phase 1: Disable foreign key checks temporarily
			if err := instance.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
				initErr = fmt.Errorf("error disabling foreign key checks: %w", err)
				return
			}

			// First, create all tables without foreign key constraints
			migrateErr := instance.AutoMigrate(
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
			)
			if migrateErr != nil {
				initErr = fmt.Errorf("error running migrations: %w", migrateErr)
				return
			}

			// Phase 2: Re-enable foreign key checks
			if err := instance.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
				initErr = fmt.Errorf("error re-enabling foreign key checks: %w", err)
				return
			}
		})
		return initErr
	}()

	return instance, onceErr
}
