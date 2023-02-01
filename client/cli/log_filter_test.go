package cli_test

import (
	"bytes"
	"testing"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	tmlog "github.com/tendermint/tendermint/libs/log"

	"github.com/functionx/fx-core/v3/client/cli"
)

func TestFxZeroLogWrapper(t *testing.T) {
	testcases := []struct {
		name     string
		logTypes []string
		level    zerolog.Level
		call     func(tmlog.Logger)
		output   string
	}{
		{
			name:     "filter msg",
			logTypes: []string{"msg key value"},
			level:    zerolog.InfoLevel,
			call: func(logger tmlog.Logger) {
				logger.Info("msg")
				logger.Info("msg key value")
			},
			output: "{\"level\":\"info\",\"message\":\"msg\"}\n",
		},
		{
			name:     "filter module with",
			logTypes: []string{"p2p:info"},
			level:    zerolog.InfoLevel,
			call: func(logger tmlog.Logger) {
				logger.With("module", "p2p").Debug("msg")
				logger.With("module", "p2p").Info("msg")
			},
			output: "{\"level\":\"info\",\"module\":\"p2p\",\"message\":\"msg\"}\n",
		},
		{
			name:     "filter module kv",
			logTypes: []string{"p2p:info"},
			level:    zerolog.InfoLevel,
			call: func(logger tmlog.Logger) {
				logger.Debug("msg", "module", "p2p")
				logger.Info("msg", "module", "p2p")
			},
			output: "{\"level\":\"info\",\"module\":\"p2p\",\"message\":\"msg\"}\n",
		},
		{
			name:     "filter module msg",
			logTypes: []string{"p2p:error"},
			level:    zerolog.DebugLevel,
			call: func(logger tmlog.Logger) {
				logger.Info("msg")
			},
			output: "{\"level\":\"info\",\"message\":\"msg\"}\n",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := server.ZeroLogWrapper{
				Logger: zerolog.New(buf).Level(testcase.level).With().Logger(),
			}
			logWrapper := cli.NewFxZeroLogWrapper(logger, testcase.logTypes)
			testcase.call(logWrapper)
			assert.Equal(t, testcase.output, buf.String())
		})
	}
}
