package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/avast/retry-go"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/config"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/storage"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/errs"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
)

type Accrual struct {
	storage  storage.Order
	cfg      *config.Config
	client   *resty.Client
	retryCfg *config.Retry
}

func NewAccrual(s storage.Order, c *config.Config) *Accrual {
	return &Accrual{
		storage:  s,
		cfg:      c,
		client:   resty.New(),
		retryCfg: config.RetryNew(),
	}
}

func (a *Accrual) Run(ctx context.Context) {
	requestTicker := time.NewTicker(time.Duration(a.cfg.RequestInterval) * time.Second)
	defer requestTicker.Stop()

	for t := range requestTicker.C {
		numbers, err := a.storage.GetProcessingOrders(ctx, a.cfg.RateLimit)
		if err != nil {
			logrus.Error("processing order collection error ", err, t)
			continue
		}
		for _, number := range numbers {
			order, err := a.getOrderStatus(number)
			if err != nil {
				logrus.Error(err)
				continue
			}
			err = a.updateOrderStatus(ctx, *order)
			logrus.Error(err)
		}
	}
}

func (a *Accrual) getOrderStatus(number string) (*models.AccrualResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", a.cfg.AccrualSystemAddress, number)
	var resp *resty.Response
	err := retry.Do(
		func() error {
			var err error
			resp, err = a.client.R().Get(url)
			if err != nil {
				return fmt.Errorf("error accrual requesting order status : %w", err)
			}
			if resp.StatusCode() == http.StatusTooManyRequests {
				retryAfterHeader := resp.Header().Get("Retry-After")
				if retryAfterHeader != "" {
					a.retryCfg.RetryAfterMillisec, err = strconv.Atoi(retryAfterHeader)
					logrus.Error(err)
				}
				return errs.ErrTooManyRequests
			}
			return nil
		},
		retry.RetryIf(a.retryCfg.IfStatusTooManyRequests),
		retry.Attempts(a.retryCfg.Attempts),
		retry.DelayType(a.retryCfg.DelayType),
	)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("received response from accrual with status code : %d", resp.StatusCode())
	}
	var order models.AccrualResponse
	err = json.Unmarshal(resp.Body(), &order)
	if err != nil {
		return nil, fmt.Errorf("invalid accrual response body : %w", err)
	}
	return &order, nil
}

func (a *Accrual) updateOrderStatus(ctx context.Context, o models.AccrualResponse) error {
	if o.Status == models.Registered {
		o.Status = models.Processing
	}
	return a.storage.Update(ctx, o)
}
