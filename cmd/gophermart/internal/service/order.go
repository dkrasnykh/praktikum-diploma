package service

import (
	"context"
	"unicode"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/storage"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/errs"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
)

type OrderService struct {
	storage storage.Order
}

func NewOrderService(storage storage.Order) *OrderService {
	return &OrderService{storage: storage}
}

func (o *OrderService) Add(ctx context.Context, userID int, orderNumber string) error {
	var id *int
	var err error
	if err = validateOrderNumber(orderNumber); err != nil {
		return err
	}
	if id, err = o.storage.GetUserIDByNumber(ctx, orderNumber); err != nil {
		return err
	}
	if id != nil && *id == userID {
		return errs.ErrOrderExist
	} else if id != nil {
		return errs.ErrUnreachableOrder
	}
	return o.storage.Add(ctx, userID, orderNumber)
}

func (o *OrderService) GetAll(ctx context.Context, userID int) ([]models.Order, error) {
	return o.storage.GetAll(ctx, userID)
}

func validateOrderNumber(number string) error {
	for _, v := range number {
		if !unicode.IsDigit(v) {
			return errs.ErrInvalidOrderNumber
		}
	}
	return nil
}
