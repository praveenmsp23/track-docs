package migrations

import (
	"errors"
	"sort"

	"github.com/praveenmsp23/trackdocs/pkg/config"
	"github.com/zerogate/gormigrate/v2"
)

type MigrationProvider interface {
	GetMigration(cfg *config.Config) *gormigrate.Migration
}

var providers = make(map[string]MigrationProvider)

func MigrationRegister(id string, provider MigrationProvider) error {
	if provider == nil {
		return errors.New("migration: provider is nil")
	}
	if _, ok := providers[id]; ok {
		return errors.New("migration: provider already exists for " + id)
	}
	providers[id] = provider
	return nil
}

func GetMigrations(cfg *config.Config) ([]*gormigrate.Migration, error) {
	keys := make([]string, 0, len(providers))
	for k := range providers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ret := []*gormigrate.Migration{}
	for _, k := range keys {
		migration := providers[k].GetMigration(cfg)
		ret = append(ret, migration)
	}
	return ret, nil
}
