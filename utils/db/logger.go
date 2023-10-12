package db

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

type GormLogger struct{}

func NewGormLogger() *GormLogger {
	return &GormLogger{}
}

const (
	ErrorStr = "ERROR"
)

func (g *GormLogger) Printf(format string, values ...interface{}) {
	if strings.Contains(fmt.Sprintf(format, values...), ErrorStr) {
		log.Error().Msgf(format, values...)
		return
	}
	log.Info().Msgf(format, values...)
}
