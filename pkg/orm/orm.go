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

			// Schema migration (table creation)
			migrateErr := instance.AutoMigrate(
				&model.Article{},
				&model.Lineup{},
				&model.Match{},
				&model.Player{},
				&model.Role{},
				&model.Season{},
				&model.TeamStat{},
				&model.Team{},
				&model.User{},
				&model.PlayerTeam{},
				&model.PlayerStat{},
			)
			if migrateErr != nil {
				initErr = fmt.Errorf("error running migrations: %w", migrateErr)
				return
			}
		})
		return initErr
	}()

	return instance, onceErr
}
