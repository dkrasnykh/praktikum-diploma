package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
)

type OrderPostgres struct {
	db           *pgxpool.Pool
	queryTimeout time.Duration
	*TxManager
}

func NewOrderPostgres(db *pgxpool.Pool, timeoutSec int) *OrderPostgres {
	return &OrderPostgres{
		db:           db,
		queryTimeout: time.Duration(timeoutSec) * time.Second,
		TxManager:    NewManager(db),
	}
}

func (o *OrderPostgres) Add(ctx context.Context, userID int, orderNumber string) error {
	newCtx, cancel := context.WithTimeout(ctx, o.queryTimeout)
	defer cancel()

	tx, err := o.TxManager.GetCtxTxOrDefault(ctx)
	if err != nil {
		return err
	}

	rows, err := tx.Query(newCtx, "INSERT INTO rewards (user_id, order_number, status) values ($1, $2, $3)",
		userID, orderNumber, models.New)
	rows.Close()
	return err
}

func (o *OrderPostgres) GetAll(ctx context.Context, userID int) ([]models.Order, error) {
	newCtx, cancel := context.WithTimeout(ctx, o.queryTimeout)
	defer cancel()

	rows, err := o.db.Query(newCtx, "SELECT (order_number, status, accrual, updated_at) "+
		"FROM rewards WHERE user_id=$1 AND withdraw IS NULL ORDER BY updated_at", userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("postgres get all orders query error: %w", err)
	}
	listResult, err := pgx.CollectRows(rows, pgx.RowTo[OrderResult])
	return ordersList(listResult), err
}

func (o *OrderPostgres) Update(ctx context.Context, order models.AccrualResponse) error {
	newCtx, cancel := context.WithTimeout(ctx, o.queryTimeout)
	defer cancel()

	if order.Accrual != nil {
		accrual := int64(*order.Accrual * 100)
		_, err := o.db.Exec(newCtx, "UPDATE rewards SET status=$1, accrual=$2, updated_at=$3 WHERE order_number=$4",
			order.Status, accrual, time.Now(), order.Order)
		return err
	}
	_, err := o.db.Exec(newCtx, "UPDATE rewards SET status=$1, updated_at=$2 WHERE order_number=$3",
		order.Status, time.Now(), order.Order)
	return err
}

func (o *OrderPostgres) GetProcessingOrders(ctx context.Context) ([]string, error) {
	newCtx, cancel := context.WithTimeout(ctx, o.queryTimeout)
	defer cancel()

	rows, err := o.db.Query(newCtx,
		"SELECT order_number FROM rewards WHERE status NOT IN ($1, $2)",
		models.Processed, models.Invalid)
	if err != nil {
		return nil, err
	}
	listResult, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, fmt.Errorf("postgres processing orders collect rows error: %w", err)
	}
	return listResult, nil
}

func (o *OrderPostgres) GetUserIDByNumber(ctx context.Context, orderNumber string) (*int, error) {
	newCtx, cancel := context.WithTimeout(ctx, o.queryTimeout)
	defer cancel()
	var err error

	tx, err := o.TxManager.GetCtxTxOrDefault(ctx)
	if err != nil {
		return nil, err
	}

	var id int
	row := tx.QueryRow(newCtx, "SELECT user_id FROM rewards WHERE order_number=$1 LIMIT 1", orderNumber)

	if err = row.Scan(&id); err != nil && errors.Is(err, pgx.ErrNoRows) {
		logrus.Error(err)
		return nil, nil
	}
	return &id, err
}
