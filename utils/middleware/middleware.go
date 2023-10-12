package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func New(r *gin.Engine) {
	r.Use(Logger())
	r.Use(gin.Recovery())
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
