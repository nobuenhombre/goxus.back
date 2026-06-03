package logfile

import (
	"log"
	"time"

	"goxus/src/internal/app/goxus/cli"

	"github.com/google/wire"
	"github.com/nobuenhombre/suikat/pkg/db/types"
)

// ProviderSet exports Wire providers for the logfile package.
var ProviderSet = wire.NewSet(
	ProvideLogFile,
	ProvideSQLLogger,
)

// ProvideLogFile opens the log file if a path was provided via CLI flags.
func ProvideLogFile(cliConfig cli.Service) (ILogFile, func(), error) {
	cleanup := func() {
		log.Println("Log File cleanup")
	}

	logFile := &LogFile{}

	if len(cliConfig.(*cli.Config).LogFile) != 0 {
		logFile.Open(cliConfig.(*cli.Config).LogFile)
		cleanup = logFile.Close
	}

	return logFile, cleanup, nil
}

// ProvideSQLLogger provides a SQL query logger for xo-generated repositories.
func ProvideSQLLogger() (types.SQLLoggerFunc, func(), error) {
	cleanup := func() {}

	logger := types.SQLLoggerFunc(func(sql string, du time.Duration, sqlParams ...any) {
		log.Printf("[SQL] %s [%v] %v", sql, du, sqlParams)
	})

	return logger, cleanup, nil
}
