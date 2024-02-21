package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/config"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/storage"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
)

type Accrual struct {
	storage storage.Order
	cfg     *config.Config
	client  *resty.Client
}

func NewAccrual(s storage.Order, c *config.Config) *Accrual {
	return &Accrual{
		storage: s,
		cfg:     c,
		client:  resty.New(),
	}
}

func (a *Accrual) Run() {
	requestTicker := time.NewTicker(time.Duration(a.cfg.RequestInterval) * time.Second)
	defer requestTicker.Stop()

	for _ = range requestTicker.C {
		numbers, err := a.storage.GetProcessingOrders(context.Background(), a.cfg.RateLimit)
		if err != nil {
			logrus.Error(err)
		}
		for _, number := range numbers {
			go a.sendRequest(number)
		}
	}
}

func (a *Accrual) sendRequest(number string) {
	url := fmt.Sprintf("http://%s/api/orders/%s", a.cfg.AccrualSystemAddress, number)
	resp, err := a.client.R().Get(url)
	if resp.StatusCode() != 200 {
		return
	}
	var order models.AccrualResponse
	err = json.Unmarshal(resp.Body(), &order)
	if err != nil {
		logrus.Error(err)
	}
	err = a.storage.Update(context.Background(), order)
	if err != nil {
		logrus.Error(err)
	}
}
