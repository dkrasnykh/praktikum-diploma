package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"unicode"

	"github.com/avast/retry-go"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

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
	err = o.storage.Add(ctx, userID, orderNumber)
	go o.accrualRequest(orderNumber)
	return err
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

func (a *OrderService) accrualRequest(number string) {
	url := fmt.Sprintf("http://%s/api/orders/%s", a.cfg.AccrualSystemAddress, number)

	err := retry.Do(
		func() error {
			resp, err := a.client.R().Get(url)
			if err != nil {
				return err
			}
			if resp.StatusCode() != http.StatusOK {
				return errs.ErrStatusIsNotFinal
			}
			var order models.AccrualResponse
			err = json.Unmarshal(resp.Body(), &order)
			if err != nil {
				return err
			}
			err = a.storage.Update(context.Background(), order)
			if err != nil {
				return err
			}
			if order.Status != "PROCESSED" && order.Status != "INVALID" {
				return errs.ErrStatusIsNotFinal
			}
			return err
		},
		retry.Attempts(config.Attempts),
		retry.DelayType(config.DelayType),
		retry.OnRetry(config.OnRetry),
	)
	logrus.Error(err)
}
