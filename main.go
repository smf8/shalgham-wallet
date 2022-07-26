package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/smf8/shalgham-voucher/pkg/database"
	"github.com/smf8/shalgham-voucher/pkg/log"
	"github.com/smf8/shalgham-voucher/pkg/router"
	"github.com/smf8/shalgham-wallet/internal/app/config"
	"github.com/smf8/shalgham-wallet/internal/app/handler"
	"github.com/smf8/shalgham-wallet/internal/app/model"
)

func main() {
	cfg := config.New()

	log.SetupLogger(cfg.LogLevel)

	app := router.New(cfg.Server)

	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		logrus.Fatalf("database failed: %s", err.Error())
	}

	profileRepo := &model.SQLProfileRepo{DB: db}

	transactionHandler := handler.Transaction{ProfileRepo: profileRepo}
	profileHandler := handler.Profile{ProfileRepo: profileRepo}

	app.Get("/healthz", handler.CheckHealth)
	api := app.Group("/api")
	api.Post("/transactions", transactionHandler.ApplyTransaction)
	api.Post("/profiles", profileHandler.Create)
	api.Get("/profiles/:phone", profileHandler.Get)

	go func() {
		if err := app.Listen(cfg.Server.Port); err != nil {
			logrus.Fatalf("http server failed: %s", err.Error())
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	s := <-sig
	logrus.Infof("signal %s received\n", s)

	if err = app.Shutdown(); err != nil {
		logrus.Errorf("failed to shutdown server: %s", err.Error())
	}
}
