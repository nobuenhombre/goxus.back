package logfile

import (
	"log"
	"time"

	"goxus/src/internal/app/goxus/cli"
	configapp "goxus/src/internal/app/goxus/config"

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
// When log.quiet is set in config, returns a no-op logger that discards SQL logs.
func ProvideSQLLogger(appConfig configapp.Service) (types.SQLLoggerFunc, func(), error) {
	cleanup := func() {}

	if appConfig.Get().Log.Quiet {
		logger := types.SQLLoggerFunc(func(sql string, du time.Duration, sqlParams ...any) {})
		return logger, cleanup, nil
	}

	logger := types.SQLLoggerFunc(func(sql string, du time.Duration, sqlParams ...any) {
		log.Printf("[SQL] %s [%v] %v", sql, du, sqlParams)
	})

	return logger, cleanup, nil
}
