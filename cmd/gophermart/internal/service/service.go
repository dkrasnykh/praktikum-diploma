package service

import (
	"context"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/storage"
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

type Service struct {
	Authorization
	Order
	Withdraw
}

func New(storage *storage.Storage) *Service {
	return &Service{
		Authorization: NewAuthService(storage),
		Order:         NewOrderService(storage),
		Withdraw:      NewWithdrawService(storage),
	}
}
