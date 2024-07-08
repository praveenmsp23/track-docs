package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/praveenmsp23/trackdocs/pkg/crypto"
	"gorm.io/gorm"
	"time"
)

const (
	IdSize    = 16
	TokenSize = 32
)

type Base struct {
	Id        string         `gorm:"primaryKey" json:"id"`
	Updated   int64          `gorm:"autoUpdateTime:milli" json:"created"`
	Created   int64          `gorm:"autoCreateTime:milli" json:"updated"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type AuditBase struct {
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
}

type TypeInterface interface {
	GetType() string
}

// H is a shortcut for map[string]interface{}
type H map[string]any

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.Id = crypto.GenerateId("ukn", IdSize)

	return nil
}

func NewSqlNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

type Jsonb map[string]interface{}

func (j Jsonb) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *Jsonb) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}
