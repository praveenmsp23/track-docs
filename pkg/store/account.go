package store

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/praveenmsp23/trackdocs/pkg/cache"
	"github.com/praveenmsp23/trackdocs/pkg/config"
	"github.com/praveenmsp23/trackdocs/pkg/logger"
	"github.com/praveenmsp23/trackdocs/pkg/models"
	"github.com/praveenmsp23/trackdocs/pkg/models/dto"
	"gorm.io/gorm"
)

type accountStore struct {
	db    *gorm.DB
	cfg   *config.Config
	cache *cache.Cache
	repo  *Store
}

const MaxVerifyAttempts = 10

const (
	AccountCachePrefix      = "account_v1::"
	AccountEmailCachePrefix = "account_email_v1::"
)

func getAccountCacheKey(accountId string) string {
	return fmt.Sprintf("%s%s", AccountCachePrefix, accountId)
}

func getAccountEmailCacheKey(email string) string {
	return fmt.Sprintf("%s%s", AccountEmailCachePrefix, sha256.Sum256([]byte(email)))
}

func newAccountStore(conn *gorm.DB, cache *cache.Cache, cfg *config.Config) *accountStore {
	return &accountStore{db: conn, cache: cache, cfg: cfg}
}

func (u *accountStore) NewAccountFromRequest(req *dto.AccountCreateRequest) (*models.Account, error) {
	return u.NewAccount(req.Name, req.Email)
}

func (u *accountStore) UpdateAccountFromRequest(accountId string, req *dto.AccountUpdateRequest) (*models.Account, error) {
	return u.UpdateAccount(accountId, req.Name)
}

func (u *accountStore) FindAccountById(uid string) (*models.Account, error) {
	var err error
	account := &models.Account{}
	err = u.cache.Get(getAccountCacheKey(uid), account)
	if err == nil && account != nil && account.Id != "" {
		return account, nil
	}
	err = u.db.Model(models.Account{}).Where("id = ?", uid).Take(&account).Error
	if err != nil {
		return &models.Account{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &models.Account{}, models.ErrAccountNotFound
	}
	err = u.cache.Set(getAccountCacheKey(uid), account)
	if err != nil {
		logger.Errorf("FindAccountById error while setting cache:%s for key %s", err.Error(), getAccountCacheKey(uid))
	}
	return account, err
}

func (u *accountStore) FindAccountByEmail(email string) (*models.Account, error) {
	var err error
	account := &models.Account{}
	var id string
	id, err = u.cache.GetString(getAccountEmailCacheKey(email))
	if err == nil && id != "" {
		return u.FindAccountById(id)
	}
	err = u.db.Model(models.Account{}).Where("email = ?", email).Take(account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &models.Account{}, models.ErrAccountNotFound
	} else if err != nil {
		return &models.Account{}, err
	}
	err = u.cache.SetString(getAccountEmailCacheKey(email), account.Id)
	if err != nil {
		logger.Errorf("FindAccountByEmail error while setting cache:%s for key %s", err.Error(), getAccountEmailCacheKey(email))
	}
	return account, err
}

func (u *accountStore) NewAccount(name, email string) (*models.Account, error) {
	account, err := u.FindAccountByEmail(email)
	if errors.Is(err, models.ErrAccountNotFound) {
		t := models.NewAccount(name, email)
		t, err := t.Create(u.db)
		if err != nil {
			return t, err
		}
		return t, nil
	} else if err != nil {
		return account, err
	}
	return account, models.ErrAccountExists
}

func (u *accountStore) UpdateAccount(accountId, name string) (*models.Account, error) {
	account, err := u.FindAccountById(accountId)
	if err != nil {
		return account, err
	}
	account.Name = name
	account, err = account.Update(u.db)
	if err != nil {
		return account, err
	}
	err = u.cache.Del(getAccountCacheKey(accountId))
	if err != nil {
		logger.Errorf("UpdateAccount error while deleting cache:%s for key %s", err.Error(), getAccountCacheKey(accountId))
	}
	return account, nil
}

func (u *accountStore) Update(account *models.Account) (*models.Account, error) {
	account, err := account.Update(u.db)
	if err != nil {
		return account, err
	}
	err = u.cache.Del(getAccountCacheKey(account.Id))
	if err != nil {
		logger.Errorf("Update account error while deleting cache:%s for key %s", err.Error(), getAccountCacheKey(account.Id))
	}
	return account, nil
}
