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
	router.POST("/api/user/orders", h.userIdentity, h.Add)
	router.GET("/api/user/orders", h.userIdentity, h.getAll)
	router.GET("/api/user/balance", h.userIdentity, h.getBalance)
	router.POST("/api/user/balance/withdraw", h.userIdentity, h.withdrawReward)
	router.GET("/api/user/withdrawals", h.userIdentity, h.getAllWithdrawals)

	return router
}
