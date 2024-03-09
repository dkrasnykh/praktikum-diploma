package api

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
)

type Authorization interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GenerateToken(ctx context.Context, login, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Order interface {
	Add(ctx context.Context, userID int, orderNumber string) error
	GetAll(ctx context.Context, userID int) ([]models.Order, error)
}

type Withdraw interface {
	GetUserBalance(ctx context.Context, userID int) (*models.UserBalance, error)
	WithdrawReward(ctx context.Context, userID int, req models.WithdrawRequest) error
	GetAllWithdrawals(ctx context.Context, userID int) ([]models.Withdraw, error)
}

type Service interface {
	Authorization
	Order
	Withdraw
}

type Handler struct {
	service Service
}

func New(s Service) *Handler {
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
