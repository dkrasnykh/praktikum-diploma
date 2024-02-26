package internal

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/api"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/config"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/service"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/storage"
)

type Server struct {
	httpServer *http.Server
	db         *pgxpool.Pool
}

func (s *Server) Run(cfg *config.Config) error {
	var err error
	s.db, err = storage.New(cfg)
	if err != nil {
		return err
	}
	if err != nil {
		logrus.Error(err)
	}
	r := storage.NewStorage(s.db, cfg.QueryTimeout)
	accrual := service.NewAccrual(r, cfg)
	go accrual.Run()
	services := service.New(r, cfg)
	handlers := api.New(services)

	s.httpServer = &http.Server{
		Addr:           cfg.RunAddress,
		Handler:        handlers.InitRoutes(),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.db.Close()
	return s.httpServer.Shutdown(ctx)
}
