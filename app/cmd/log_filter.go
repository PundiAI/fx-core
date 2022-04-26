package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/server"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

const (
	FlagLogFilter = "log_filter"
)

func AddCmdLogWrapFilterLogType(cmd *cobra.Command) error {
	filterLogTypes, err := cmd.Flags().GetStringSlice(FlagLogFilter)
	if err != nil {
		return err
	}
	if len(filterLogTypes) <= 0 {
		return nil
	}
	serverCtx := server.GetServerContextFromCmd(cmd)
	if zeroLog, ok := serverCtx.Logger.(server.ZeroLogWrapper); ok {
		serverCtx.Logger = NewFxZeroLogWrapper(zeroLog, filterLogTypes)
	}
	return server.SetCmdServerContext(cmd, serverCtx)
}

func NewFxZeroLogWrapper(logger server.ZeroLogWrapper, logTypes []string) FxZeroLogWrapper {
	filterLogMap := make(map[string]bool, len(logTypes))
	for _, logType := range logTypes {
		filterLogMap[logType] = true
	}
	fmt.Printf(" ---------- filter log wrapper -------------")
	fmt.Printf("filter log type:%v\n", strings.Join(logTypes, ","))
	return FxZeroLogWrapper{ZeroLogWrapper: logger, filterLogs: filterLogMap}
}

type FxZeroLogWrapper struct {
	server.ZeroLogWrapper
	filterLogs map[string]bool
}

var _ tmlog.Logger = (*FxZeroLogWrapper)(nil)

func (z FxZeroLogWrapper) Info(msg string, keyVals ...interface{}) {
	if exists, ok := z.filterLogs[msg]; exists || ok {
		return
	}
	z.Logger.Info().Fields(getLogFields(keyVals...)).Msg(msg)
}

func getLogFields(keyVals ...interface{}) map[string]interface{} {
	if len(keyVals)%2 != 0 {
		return nil
	}

	fields := make(map[string]interface{})
	for i := 0; i < len(keyVals); i += 2 {
		fields[keyVals[i].(string)] = keyVals[i+1]
	}

	return fields
}

// With returns a new wrapped logger with additional context provided by a set
// of key/value tuples. The number of tuples must be even and the key of the
// tuple must be a string.
func (z FxZeroLogWrapper) With(keyVals ...interface{}) tmlog.Logger {
	return FxZeroLogWrapper{
		server.ZeroLogWrapper{Logger: z.ZeroLogWrapper.Logger.With().Fields(getLogFields(keyVals...)).Logger()},
		z.filterLogs,
	}
}
