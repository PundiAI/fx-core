package server

import (
	"strings"

	"github.com/rs/zerolog"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

const (
	FlagLogFilter = "log_filter"
)

func NewFxZeroLogWrapper(logger zerolog.Logger, logTypes []string) FxZeroLogWrapper {
	filterMsg := make(map[string]bool)
	filterModule := make(map[string]zerolog.Level)
	for _, logType := range logTypes {
		if kv := strings.Split(logType, ":"); len(kv) == 2 {
			filterModule[kv[0]], _ = zerolog.ParseLevel(kv[1])
		} else {
			filterMsg[logType] = true
		}
	}
	return FxZeroLogWrapper{Logger: logger, filterMsg: filterMsg, filterModule: filterModule}
}

var _ tmlog.Logger = (*FxZeroLogWrapper)(nil)

type FxZeroLogWrapper struct {
	zerolog.Logger
	filterMsg    map[string]bool
	filterModule map[string]zerolog.Level
}

var _ tmlog.Logger = (*FxZeroLogWrapper)(nil)

func (z FxZeroLogWrapper) Debug(msg string, keyVals ...interface{}) {
	if exists, ok := z.filterMsg[msg]; exists || ok {
		return
	}
	fields, level := z.getLogFields(keyVals...)
	if level > zerolog.DebugLevel {
		return
	}
	z.Logger.Debug().Fields(fields).Msg(msg)
}

func (z FxZeroLogWrapper) Info(msg string, keyVals ...interface{}) {
	if exists, ok := z.filterMsg[msg]; exists || ok {
		return
	}
	fields, level := z.getLogFields(keyVals...)
	if level > zerolog.InfoLevel {
		return
	}
	z.Logger.Info().Fields(fields).Msg(msg)
}

func (z FxZeroLogWrapper) Error(msg string, keyVals ...interface{}) {
	if exists, ok := z.filterMsg[msg]; exists || ok {
		return
	}
	fields, level := z.getLogFields(keyVals...)
	if level > zerolog.ErrorLevel {
		return
	}
	z.Logger.Error().Fields(fields).Msg(msg)
}

func (z FxZeroLogWrapper) With(keyVals ...interface{}) tmlog.Logger {
	fields, level := z.getLogFields(keyVals...)
	logger := z.Logger.Level(level).With().Fields(fields).Logger()
	return FxZeroLogWrapper{logger, z.filterMsg, z.filterModule}
}

func (z FxZeroLogWrapper) getLogFields(keyVals ...interface{}) (fields map[string]interface{}, logLevel zerolog.Level) {
	logLevel = z.GetLevel()
	if len(keyVals)%2 != 0 {
		return nil, logLevel
	}

	fields = make(map[string]interface{})
	for i := 0; i < len(keyVals); i += 2 {
		key, ok := keyVals[i].(string)
		if !ok {
			continue
		}
		if key == "module" {
			v, _ := keyVals[i+1].(string)
			if level, ok := z.filterModule[v]; ok && level != zerolog.NoLevel {
				logLevel = level
			}
		}
		fields[key] = keyVals[i+1]
	}

	return fields, logLevel
}
