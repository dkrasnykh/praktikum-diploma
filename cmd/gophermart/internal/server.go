package internal

import (
	"context"
	"net/http"
	"time"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/api"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/config"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/service"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/service/accrual"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/storage"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg *config.Config) error {
	db, err := storage.New(cfg)
	if err != nil {
		return err
	}
	tm := storage.NewManager(db)
	r := storage.NewStorage(db, tm, cfg.QueryTimeout)
	accrualService := accrual.New(r, cfg)
	go accrualService.Run(context.Background())
	services := service.New(r, tm, cfg)
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
