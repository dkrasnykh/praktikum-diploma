package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/config"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	cfg, err := config.New()
	if err != nil {
		logrus.Fatal(err)
	}
	srv := new(internal.Server)
	go func() {
		if err := srv.Run(cfg); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	logrus.Print("TodoApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("TodoApp Shutting Down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
}
