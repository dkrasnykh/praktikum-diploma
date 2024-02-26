package api

import (
	"github.com/gin-gonic/gin"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/service"
)

type Handler struct {
	service *service.Service
}

func New(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.POST("/api/user/register", h.signUp)
	router.POST("/api/user/login", h.signIn)

	api := router.Group("/api/user", h.userIdentity)
	{
		api.POST("/orders", h.Add)
		api.GET("/orders", h.getAll)
		api.GET("/balance", h.getBalance)
		api.POST("/balance/withdraw", h.withdrawReward)
		api.GET("/withdrawals", h.getAllWithdrawals)
	}
	return router
}
