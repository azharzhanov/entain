// Package main API.
//
// The purpose of this application is to provide a complete API for Entain technical task.
//
//	Schemes: http, https
//	Host: localhost
//	BasePath: /
//	Version: 0.127.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pkgAccount "entain/internal/account"
	pkgDatabase "entain/internal/database/postgres"
	pkgHelpers "entain/internal/helpers"
	pkgTransaction "entain/internal/transactions"
	pkgPostgres "entain/pkg/database/postgres"

	kitzapadapter "github.com/go-kit/kit/log/zap"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Create a single logger, which we'll use and give to other components.
	zapLogger, _ := zap.NewProduction()
	defer func() {
		_ = zapLogger.Sync()
	}()

	// Create a logger instance using go-kit zap adapter.
	var logger = kitzapadapter.NewZapSugarLogger(zapLogger, zapcore.InfoLevel)
	logger = log.With(logger,
		"ts", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)

	// Logging helper function
	logFatal := func(err error) {
		_ = logger.Log("err", err)
		os.Exit(1)
	}

	// Parse command line flags.
	fs := flag.NewFlagSet("", flag.ExitOnError)
	httpAddr := fs.String("http-addr", ":8080", "HTTP listen address")
	fs.Usage = pkgHelpers.UsageFor(fs, os.Args[0]+" [flags]")
	err := fs.Parse(os.Args[1:])
	if err != nil {
		logFatal(err)
	}

	// Reads configuration from the environment variables.
	cfg, err := pkgHelpers.LoadConfig()
	if err != nil {
		logFatal(err)
	}

	// Validate configurations
	if err := pkgHelpers.ValidateConfig(cfg); err != nil {
		logFatal(err)
	}

	// Setup database connection.
	db, err := pkgPostgres.NewConnection(cfg.DSN)
	if err != nil {
		logFatal(err)
	}
	defer db.Close()

	// Setup database.
	sqlcRepository, err := pkgDatabase.NewSqlcRepository(db)
	if err != nil {
		logFatal(err)
	}
	defer sqlcRepository.Close()

	// Repository layer
	var (
		transactionRepository = pkgDatabase.NewTransactionRepository(sqlcRepository)
		accountRepository     = pkgDatabase.NewAccountRepository(sqlcRepository)
	)

	// Service layer
	var (
		transactionService = pkgTransaction.NewService(
			transactionRepository,
			accountRepository,
			sqlcRepository,
		)
		accountService = pkgAccount.NewService(accountRepository)
	)

	// Endpoints layer
	var (
		transactionEndpoint = pkgTransaction.NewEndpoints(transactionService)
		accountndpoint      = pkgAccount.NewEndpoints(accountService)
	)

	// Transport layer
	r := mux.NewRouter()
	{
		pkgTransaction.RegisterRoutersV1(
			r,
			transactionEndpoint,
			logger,
		)
		pkgAccount.RegisterRoutersV1(
			r,
			accountndpoint,
			logger,
		)
	}

	http.Handle("/", r)

	// This function just sits and waits for ctrl-C.
	// Implementing this way allows us to gracefully shutdown the server
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	go func() {
		_ = logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	_ = logger.Log("exit", <-errs)
}
