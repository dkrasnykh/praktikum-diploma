package service

import (
	"context"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/go-resty/resty/v2"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/config"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/storage"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/errs"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
)

type OrderService struct {
	storage storage.Order
	client  *resty.Client
	cfg     *config.Config
}

func NewOrderService(storage storage.Order, cfg *config.Config) *OrderService {
	return &OrderService{
		storage: storage,
		client:  resty.New(),
		cfg:     cfg,
	}
}

func (o *OrderService) Add(ctx context.Context, userID int, orderNumber string) error {
	if err := goluhn.Validate(orderNumber); err != nil {
		return errs.ErrInvalidOrderNumber
	}
	return o.storage.Add(ctx, userID, orderNumber)
}

func (o *OrderService) GetAll(ctx context.Context, userID int) ([]models.Order, error) {
	return o.storage.GetAll(ctx, userID)
}
