package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
)

type WithdrawPostrges struct {
	db           *pgxpool.Pool
	queryTimeout time.Duration
}

func NewWithdrawPostgres(db *pgxpool.Pool, timeoutSec int) *WithdrawPostrges {
	return &WithdrawPostrges{
		db:           db,
		queryTimeout: time.Duration(timeoutSec) * time.Second,
	}
}

func (s *WithdrawPostrges) GetUserBalance(ctx context.Context, userID int) (*models.UserBalance, error) {
	newCtx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	var debit, credit int64
	row := s.db.QueryRow(newCtx, "SELECT SUM(accrual) as total, SUM(withdraw) as withdraw "+
		"FROM (SELECT COALESCE(accrual, 0) as accrual, COALESCE(withdraw, 0) as withdraw FROM rewards WHERE user_id=$1) as t1", userID)
	if err := row.Scan(&debit, &credit); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	return &models.UserBalance{Current: float32(debit-credit) / 100, Withdrawn: float32(credit) / 100}, nil
}

func (s *WithdrawPostrges) WithdrawReward(ctx context.Context, userID int, req models.WithdrawRequest) error {
	newCtx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	_, err := s.db.Query(newCtx, "INSERT INTO rewards (user_id, order_number, status, withdraw) values ($1, $2, $3, $4)",
		userID, req.Order, models.Processed, int64(req.Sum*100))
	return err
}

func (s *WithdrawPostrges) GetAllWithdrawals(ctx context.Context, userID int) ([]models.Withdraw, error) {
	newCtx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	rows, err := s.db.Query(newCtx,
		"SELECT (order_number, withdraw, updated_at) FROM rewards WHERE withdraw IS NOT NULL AND user_id=$1 ORDER BY updated_at", userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("postgres all withdrawal query error: %w", err)
	}
	listResult, err := pgx.CollectRows(rows, pgx.RowTo[WithdrawResult])
	return withdrawList(listResult), err
}
