package models

import (
	"database/sql"
	"github.com/praveenmsp23/trackdocs/pkg/crypto"
	"gorm.io/gorm"
)

type UserStatus string

const (
	UserStatusActive      UserStatus = "active"
	UserStatusLocked      UserStatus = "locked"
	UserStatusSuspended   UserStatus = "suspended"
	UserStatusDeActivated UserStatus = "deactivated"
)

type Account struct {
	Base
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	Status      UserStatus   `json:"status"`
	LastLoginAt sql.NullTime `json:"last_login_at"`
}

func NewAccount(name, email string) *Account {
	return &Account{
		Name:   name,
		Email:  email,
		Status: UserStatusActive,
	}
}

func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	a.Id = crypto.GenerateId("acc", IdSize)
	return nil
}

func (a *Account) Create(db *gorm.DB) (*Account, error) {
	err := db.Create(&a).Error
	if err != nil {
		return &Account{}, err
	}
	return a, nil
}

func (a *Account) Update(db *gorm.DB) (*Account, error) {
	db = db.Model(&Account{}).Where("id = ?", a.Id).UpdateColumns(
		map[string]interface{}{
			"email":         a.Email,
			"name":          a.Name,
			"last_login_at": a.LastLoginAt,
			"status":        a.Status,
		},
	)
	if db.Error != nil {
		return &Account{}, db.Error
	}
	// This is the display the updated user
	err := db.Model(&Account{}).Where("id = ?", a.Id).Take(&a).Error
	if err != nil {
		return &Account{}, err
	}
	return a, nil
}

func (a *Account) Delete(db *gorm.DB) (int64, error) {
	db = db.Model(&Account{}).Where("id = ?", a.Id).Take(&Account{}).Delete(&Account{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
