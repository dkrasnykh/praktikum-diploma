package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type TxManager struct {
	db *pgxpool.Pool
}

func NewManager(db *pgxpool.Pool) *TxManager {
	return &TxManager{
		db: db,
	}
}

func (m *TxManager) InitTx(ctx context.Context) (context.Context, error) {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, "tx", tx), nil
}

func (m *TxManager) CompleteTx(ctx context.Context) {
	value := ctx.Value("tx")
	tx, ok := value.(pgx.Tx)
	if !ok {
		logrus.Error("context value with key 'tx' not found")
	}
	err := tx.Commit(ctx)
	if err != nil {
		logrus.Error("commit transaction ", err)
		err = tx.Rollback(ctx)
		if err != nil {
			logrus.Error("rollback transaction ", err)
		}
	}
}

func (m *TxManager) GetCtxTxOrDefault(ctx context.Context) (pgx.Tx, error) {
	value := ctx.Value("tx")
	tx, ok := value.(pgx.Tx)
	var err error
	if !ok {
		tx, err = m.db.Begin(ctx)
		if err != nil {
			return nil, err
		}
	}
	return tx, nil
}
