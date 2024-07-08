package db

import (
	"github.com/praveenmsp23/trackdocs/pkg/config"
	"github.com/praveenmsp23/trackdocs/pkg/lock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

const (
	DefaultMaxLifetime        = time.Hour * 1
	DefaultMaxIdleConnections = 10
	DefaultMaxOpenConnections = 30
)

func NewDB(cfg *config.Config, lock *lock.RedisLock) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  cfg.Datasource,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(DefaultMaxIdleConnections)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(DefaultMaxOpenConnections)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(DefaultMaxLifetime)
	if cfg.Env == config.ApplicationEnvLocal {
		db.Logger = logger.Default.LogMode(logger.Info)
	} else {
		db.Logger = logger.Default.LogMode(logger.Silent)
	}
	err = Migrate(cfg, lock, db)
	return db, err
}
