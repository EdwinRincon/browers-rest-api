package orm

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/EdwinRincon/browersfc-api/config"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
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
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening SQL connection: %w", err)
	}
	return sqlDB, nil
}

func openGormConnection(sqlDB *sql.DB) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger:                                   NewContextAwareGormLogger(),
		DisableForeignKeyConstraintWhenMigrating: true, // Disable FK constraints during migration
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), gormConfig)
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
	// Migrate tables in dependency order to avoid foreign key constraint errors
	// Migrate each table individually to avoid GORM processing relationships prematurely

	// Step 1: Base tables with no foreign key dependencies
	if err := db.AutoMigrate(&model.Role{}); err != nil {
		return fmt.Errorf("error migrating role table: %w", err)
	}
	if err := db.AutoMigrate(&model.Season{}); err != nil {
		return fmt.Errorf("error migrating season table: %w", err)
	}
	if err := db.AutoMigrate(&model.Team{}); err != nil {
		return fmt.Errorf("error migrating team table: %w", err)
	}

	// Step 2: Tables that depend on base tables
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return fmt.Errorf("error migrating user table: %w", err)
	}
	if err := db.AutoMigrate(&model.Article{}); err != nil {
		return fmt.Errorf("error migrating article table: %w", err)
	}

	// Step 3: Player (depends on User, but UserID is nullable)
	if err := db.AutoMigrate(&model.Player{}); err != nil {
		return fmt.Errorf("error migrating player table: %w", err)
	}

	// Step 4: Tables that depend on Team, Season, and Player
	if err := db.AutoMigrate(&model.Match{}); err != nil {
		return fmt.Errorf("error migrating match table: %w", err)
	}
	if err := db.AutoMigrate(&model.PlayerTeam{}); err != nil {
		return fmt.Errorf("error migrating player_team table: %w", err)
	}
	if err := db.AutoMigrate(&model.TeamStat{}); err != nil {
		return fmt.Errorf("error migrating team_stat table: %w", err)
	}

	// Step 5: Tables that depend on Match
	if err := db.AutoMigrate(&model.Lineup{}); err != nil {
		return fmt.Errorf("error migrating lineup table: %w", err)
	}
	if err := db.AutoMigrate(&model.PlayerStat{}); err != nil {
		return fmt.Errorf("error migrating player_stat table: %w", err)
	}

	return nil
}
