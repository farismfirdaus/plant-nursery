package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/farismfirdaus/plant-nursery/auth"
	apperr "github.com/farismfirdaus/plant-nursery/errors"
	"github.com/farismfirdaus/plant-nursery/utils/response"
)

var client auth.Auth

func New(r *gin.Engine, c auth.Auth) {
	r.Use(Logger())
	r.Use(gin.Recovery())

	client = c
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process Request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		statusCode := c.Writer.Status()

		logFields := map[string]interface{}{
			"status_code":   statusCode,
			"path":          c.Request.URL.Path,
			"method":        c.Request.Method,
			"start_time":    start.Format("2006/01/02 - 15:04:05"),
			"response_time": latency,
		}

		if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {
			log.Info().Fields(logFields).Send()
		} else {
			// populate errors
			errs := []error{}
			for _, e := range c.Errors.Errors() {
				errs = append(errs, errors.New(e))
			}

			// log error
			log.Error().Fields(logFields).Errs("errors", errs).Send()
		}
	}
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		token := c.Request.Header.Get("Authorization")

		if !strings.HasPrefix(token, "Bearer ") {
			response.BuildErrors(c, apperr.Unauthorized)
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		customerID, err := client.Verify(ctx, token)
		if err != nil {
			response.BuildErrors(c, apperr.Unauthorized)
			c.Abort()
			return
		}

		c.Set(gin.AuthUserKey, customerID)
		c.Next()
	}
}
