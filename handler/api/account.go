package api

import (
	"github.com/gin-gonic/gin"
	"github.com/praveenmsp23/trackdocs/pkg/models"
	"github.com/praveenmsp23/trackdocs/pkg/models/dto"
	"github.com/praveenmsp23/trackdocs/pkg/server"
	"github.com/praveenmsp23/trackdocs/pkg/store"
	"github.com/praveenmsp23/trackdocs/pkg/token"
	"net/http"
)

func HandleGetAccount() gin.HandlerFunc {
	return server.HandleFunc(func(c *models.TrackDocsContext) {
		if c.Account == nil {
			c.Error(models.ErrTokenExpired)
			return
		}
		account := c.Account
		c.JSON(http.StatusOK, models.NewSuccessResponse(account))
	})
}

func HandleAccountUpdate(repo *store.Store) gin.HandlerFunc {
	return server.HandleFunc(func(c *models.TrackDocsContext) {
		if c.Account == nil {
			c.Error(models.ErrTokenExpired)
			return
		}
		account := c.Account
		var json dto.AccountUpdateRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.Error(err).SetType(gin.ErrorTypeBind)
			return
		}
		account, err := repo.AccountStore.UpdateAccountFromRequest(account.Id, &json)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(account))
	})
}

func HandleAccountLogout(token *token.Manager) gin.HandlerFunc {
	return server.HandleFunc(func(c *models.TrackDocsContext) {
		token.TokenDestroy(c.Context)
		c.JSON(http.StatusOK, models.NewSuccessResponse(gin.H{}))
	})
}
