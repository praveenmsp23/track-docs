package api

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/praveenmsp23/trackdocs/pkg/cache"
	"github.com/praveenmsp23/trackdocs/pkg/config"
	"github.com/praveenmsp23/trackdocs/pkg/logger"
	"github.com/praveenmsp23/trackdocs/pkg/models"
	"github.com/praveenmsp23/trackdocs/pkg/store"
	"github.com/praveenmsp23/trackdocs/pkg/token"
)

func AuthMiddleware(s *store.Store, manager *token.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := models.NewTrackDocsContext(c)
		t := manager.TokenGet(p.Context)
		if t == nil || t.TokenID() == "" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, models.ErrUnauthorized))
			c.Abort()
			return
		}
		accountId, isExists := t.Get("account_id")
		if !isExists || accountId == "" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, models.ErrTokenExpired))
			c.Abort()
			return
		}
		account, err := s.AccountStore.FindAccountById(accountId)
		if err != nil {
			logger.Error(err)
			c.Error(err)
			c.Abort()
			return
		}
		c.Set("account", account)
		c.Next()
	}
}

// RateLimitMiddleware is a Gin middleware that limits the rate of API authentication requests based on API key
func RateLimitMiddleware(limit int, cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := models.NewTrackDocsContext(c)
		account := p.Account
		if account == nil {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, models.ErrUnauthorized))
			c.Abort()
			return
		}
		// Create a Redis key for the rate limiter
		limitKey := fmt.Sprintf("ratelimit::%s", account.Id)

		res, err := cache.Allow(limitKey, limit)
		if err != nil && err != redis.Nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.Writer.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Writer.Header().Set("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))
		c.Writer.Header().Set("X-RateLimit-Policy", fmt.Sprintf("%d;w=60", limit))
		if res.Allowed == 0 {
			retryAfter := strconv.Itoa(int(res.RetryAfter / time.Second))
			c.Writer.Header().Set("X-RateLimit-Reset", retryAfter)
			c.Writer.Header().Set("Retry-After", retryAfter)
			c.JSON(http.StatusTooManyRequests, models.NewErrorResponse(http.StatusTooManyRequests, fmt.Errorf("exceeds rate limit, retry in %s second(s", retryAfter)))
			c.Abort()
			return
		}

		// Call the next handler in the chain
		c.Next()
	}
}

func UcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func LcFirst(str string) string {
	return strings.ToLower(str)
}

func Split(src string) string {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return src
	}
	var entries []string
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}

	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}

	for index, word := range entries {
		if index == 0 {
			entries[index] = UcFirst(word)
		} else {
			entries[index] = LcFirst(word)
		}
	}
	justString := strings.Join(entries, " ")
	return justString
}

// ValidationErrorToText Check https://github.com/go-playground/validator
func ValidationErrorToText(e validator.FieldError) string {
	word := Split(e.Field())
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", word)
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", word, e.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s", word, e.Param())
	case "email":
		return "invalid email format"
	case "len":
		return fmt.Sprintf("%s must be %s characters long", word, e.Param())
	}
	return fmt.Sprintf("%s is not valid", word)
}

// Errors This method collects all errors and submits them to Rollbar
func Errors(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// Only run if there are some errors to handle
		if len(c.Errors) > 0 {
			errorMessage := ""

			for _, e := range c.Errors {
				logger.Errorf("Error in the request %v", e)
				errorMessage += fmt.Sprintf("Error: %s\n", e.Error())
				// Find out what type of error it is
				switch e.Type {
				case gin.ErrorTypePublic:
					// Only output public errors if nothing has been written yet
					if !c.Writer.Written() {
						c.JSON(models.ErrorStatusCode(e), models.NewErrorResponse(models.ErrorStatusCode(e), e))
						//c.JSON(c.Writer.Status(), gin.H{"Error": e.Error()})
						break
					}
				case gin.ErrorTypeBind:
					errs := e.Err.(validator.ValidationErrors)
					ret := ""
					for _, err := range errs {
						ret = fmt.Sprintf("%s%s: %s \n", ret, err.Field(), ValidationErrorToText(err))
					}

					c.JSON(http.StatusBadRequest, models.NewErrorsResponse(http.StatusBadRequest, ret))
					break

				default:
					if models.IsErrorCustom(e) {
						c.JSON(models.ErrorStatusCode(e), models.NewErrorResponse(models.ErrorStatusCode(e), e))
						break
					}
				}

			}
			if !c.Writer.Written() {
				if cfg.Env == config.ApplicationEnvLocal {
					c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, fmt.Errorf(errorMessage)))
				} else {
					c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, models.ErrInternalServer))
				}
			}
		}
	}
}
