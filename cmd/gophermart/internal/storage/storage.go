package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
)

//go:generate mockgen -source=storage.go -destination=mocks/mock.go

type Transaction interface {
	GetCtxTxOrDefault(ctx context.Context) (pgx.Tx, error)
}

type Authorization interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GetUser(ctx context.Context, username, password string) (*models.User, error)
}

type Order interface {
	Add(ctx context.Context, userID int, orderNumber string) error
	GetAll(ctx context.Context, userID int) ([]models.Order, error)
	Update(ctx context.Context, order models.AccrualResponse) error
	GetProcessingOrders(ctx context.Context) ([]string, error)
	GetUserIDByNumber(ctx context.Context, orderNumber string) (*int, error)
}

type Withdraw interface {
	GetUserBalance(ctx context.Context, userID int) (current, withdrawn int64, err error)
	WithdrawReward(ctx context.Context, userID int, req models.WithdrawRequest) error
	GetAllWithdrawals(ctx context.Context, userID int) ([]models.Withdraw, error)
}

type Storage struct {
	Authorization
	Order
	Withdraw
}

func NewStorage(db *pgxpool.Pool, tm Transaction, timeoutSec int) *Storage {
	return &Storage{
		Authorization: NewAuthPostgres(db, timeoutSec),
		Order:         NewOrderPostgres(db, tm, timeoutSec),
		Withdraw:      NewWithdrawPostgres(db, tm, timeoutSec),
	}
}
