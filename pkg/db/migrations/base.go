package migrations

import "gorm.io/gorm"

const (
	IdSize = 16
)

type Base struct {
	Id        string         `gorm:"size:36;primaryKey;" json:"id"`
	Updated   int64          `gorm:"autoUpdateTime:milli" json:"created"`
	Created   int64          `gorm:"autoCreateTime:milli" json:"updated"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type AuditBase struct {
	CreatedBy  string `gorm:"size:36;not null;"`
	ModifiedBy string `gorm:"size:36;"`
}
