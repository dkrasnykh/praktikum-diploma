package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/errs"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
)

type OrderPostgres struct {
	db           *pgxpool.Pool
	queryTimeout time.Duration
}

func NewOrderPostgres(db *pgxpool.Pool, timeoutSec int) *OrderPostgres {
	return &OrderPostgres{
		db:           db,
		queryTimeout: time.Duration(timeoutSec) * time.Second,
	}
}

func (o *OrderPostgres) Add(ctx context.Context, userID int, orderNumber string) error {
	newCtx, cancel := context.WithTimeout(ctx, o.queryTimeout)
	defer cancel()

	tx, err := o.db.Begin(newCtx)
	if err != nil {
		return err
	}
	var id sql.NullInt32
	row := tx.QueryRow(newCtx, "SELECT user_id FROM rewards WHERE order_number=$1 LIMIT 1", orderNumber)
	if err = row.Scan(&id); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		exErr := fmt.Errorf("select user_id by number exec %w", err)
		err = tx.Rollback(newCtx)
		if err != nil {
			return errors.Join(exErr, err)
		}
		return exErr
	}

	if id.Valid && int(id.Int32) == userID {
		return errs.ErrOrderExist
	} else if id.Valid {
		return errs.ErrUnreachableOrder
	}

	rows, err := tx.Query(newCtx, "INSERT INTO rewards (user_id, order_number, status) values ($1, $2, $3)",
		userID, orderNumber, models.New)
	if err != nil {
		exErr := fmt.Errorf("insert new order into db exec %w", err)
		err = tx.Rollback(newCtx)
		if err != nil {
			return errors.Join(exErr, err)
		}
		return exErr
	}
	rows.Close()
	return tx.Commit(newCtx)
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

func (o *OrderPostgres) GetProcessingOrders(ctx context.Context, limit int) ([]string, error) {
	newCtx, cancel := context.WithTimeout(ctx, o.queryTimeout)
	defer cancel()

	rows, err := o.db.Query(newCtx,
		"SELECT order_number FROM rewards WHERE status NOT IN ($1, $2) ORDER BY updated_at LIMIT $3",
		models.Processed, models.Invalid, limit)
	if err != nil {
		return nil, err
	}
	listResult, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, fmt.Errorf("postgres processing orders collect rows error: %w", err)
	}
	return listResult, nil
}
