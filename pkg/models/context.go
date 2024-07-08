package models

import (
	"github.com/gin-gonic/gin"
)

type TrackDocsContext struct {
	*gin.Context
	Account *Account
}

func NewTrackDocsContext(c *gin.Context) *TrackDocsContext {
	z := &TrackDocsContext{Context: c}
	if obj, ok := z.Get("account"); ok && obj != nil {
		z.Account = obj.(*Account)
	}
	return z
}
