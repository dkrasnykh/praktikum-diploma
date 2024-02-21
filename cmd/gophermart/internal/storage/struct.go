package storage

import (
	"database/sql"
	"time"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
)

type OrderResult struct {
	Number     string
	Status     string
	Accure     sql.NullInt64
	UploadedAt time.Time
}

type WithdrawResult struct {
	Number      string
	Withdraw    int64
	ProcessedAt time.Time
}

func ordersList(orderResult []OrderResult) []models.Order {
	orders := make([]models.Order, 0, len(orderResult))
	for _, v := range orderResult {
		o := models.Order{Number: v.Number, Status: v.Status, UploadedAt: v.UploadedAt.Format(time.RFC3339)}
		if v.Accure.Valid {
			value := float64(v.Accure.Int64) / 100
			o.Accrual = &value
		}
		orders = append(orders, o)
	}
	return orders
}

func withdrawList(withdawResult []WithdrawResult) []models.Withdraw {
	withdrawals := make([]models.Withdraw, 0, len(withdawResult))
	for _, v := range withdawResult {
		w := models.Withdraw{Order: v.Number, Sum: float64(v.Withdraw) / 100, ProcessedAt: v.ProcessedAt.Format(time.RFC3339)}
		withdrawals = append(withdrawals, w)
	}
	return withdrawals
}
