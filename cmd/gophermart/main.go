package main

import (
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
	if err := srv.Run(cfg); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}
