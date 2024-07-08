package migrations

import (
	"database/sql"
	"github.com/praveenmsp23/trackdocs/pkg/config"
	"github.com/zerogate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	MigrationRegister("001", &InitialMigrationProvider{})
}

type Account struct {
	Base
	Name        string `gorm:"size:256;not null;"`
	Email       string `gorm:"size:256;not null;unique"`
	Status      string `gorm:"size:20;default:'active';index"`
	LastLoginAt sql.NullTime
}

type InitialMigrationProvider struct{}

func (m InitialMigrationProvider) GetMigration(cfg *config.Config) *gormigrate.Migration {
	return &gormigrate.Migration{
		ID:       "001",
		Migrate:  m.Migrate,
		Rollback: m.Rollback,
	}
}

func (m InitialMigrationProvider) Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&Account{}); err != nil {
		return err
	}
	return nil
}

func (m InitialMigrationProvider) Rollback(tx *gorm.DB) error {
	if err := tx.Migrator().DropTable(&Account{}); err != nil {
		return err
	}
	return nil
}
