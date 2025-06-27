package main

import (
	"os"

	pkgHelpers "entain/internal/helpers"
	pkgDatabase "entain/pkg/database/postgres"
	"entain/schema/postgresql"

	kitzapadapter "github.com/go-kit/kit/log/zap"
	"github.com/go-kit/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	_ "github.com/lib/pq"
)

// AppVersion - application version
var AppVersion = "unversioned"

func main() {
	// Create a single logger, which we'll use and give to other components
	//
	zapLogger, _ := zap.NewProduction()
	defer func() {
		_ = zapLogger.Sync()
	}()

	var logger log.Logger
	logger = kitzapadapter.NewZapSugarLogger(zapLogger, zapcore.InfoLevel)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	logger = log.With(logger, "ver", AppVersion)
	// Logging helper function
	logFatal := func(err error) {
		_ = logger.Log("err", err)
		os.Exit(1)
	}

	// Reads configuration from the environment variables
	//
	cfg, err := pkgHelpers.LoadConfig()
	if err != nil {
		logFatal(err)
	}

	// Setup database connection
	//
	db, err := pkgDatabase.NewConnection(cfg.DSN)
	if err != nil {
		logFatal(err)
	}

	if err := postgresql.MigrateUp(db); err != nil {
		logFatal(err)
	}
}
