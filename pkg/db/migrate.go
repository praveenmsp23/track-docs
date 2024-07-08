package db

import (
	"errors"
	"github.com/zerogate/gormigrate/v2"
	"time"

	"github.com/praveenmsp23/trackdocs/pkg/config"
	"github.com/praveenmsp23/trackdocs/pkg/db/migrations"
	"github.com/praveenmsp23/trackdocs/pkg/lock"
	"github.com/praveenmsp23/trackdocs/pkg/logger"
	"gorm.io/gorm"
)

func Migrate(cfg *config.Config, redisLock *lock.RedisLock, db *gorm.DB) error {
	mig, err := migrations.GetMigrations(cfg)
	if err != nil {
		logger.Fatalf("could not migrate: %v", err)
		return err
	}
	logger.Info("acquiring lock for db_migrate")
	mutex := redisLock.NewMutex("db_migrate", lock.WithExpiry(time.Minute*10), lock.WithRetryDelay(time.Second*5), lock.WithRetryCount(1000))
	ok, err := mutex.Lock()
	if err != nil {
		logger.Fatalf("could not migrate: Error while get lock on db_migrate %v", err)
		return err
	}
	if !ok {
		logger.Fatalf("could not migrate: Unable get lock on db_migrate %v", ok)
		return errors.New("unable get lock on db_migrate")
	}
	defer mutex.Unlock()
	logger.Info("lock acquired successfully for db_migrate")

	m := gormigrate.New(db, gormigrate.DefaultOptions, mig)
	if err = m.Migrate(); err != nil {
		logger.Fatalf("could not migrate: %v", err)
		return err
	}
	logger.Infof("migration did run successfully.")
	return nil
}
