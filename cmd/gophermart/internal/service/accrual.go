package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

	for t := range requestTicker.C {
		numbers, err := a.storage.GetProcessingOrders(context.Background(), a.cfg.RateLimit)
		if err != nil {
			logrus.Error("processing order collection error ", err, t)
			continue
		}
		for _, number := range numbers {
			go a.sendRequest(number)
		}
	}
}

func (a *Accrual) sendRequest(number string) {
	url := fmt.Sprintf("%s/api/orders/%s", a.cfg.AccrualSystemAddress, number)
	resp, err := a.client.R().Get(url)
	if err != nil {
		logrus.Error("error accrual requesting order status ", err)
		return
	}
	if resp.StatusCode() != http.StatusOK {
		logrus.Info("received response from accrual with status code ", resp.StatusCode())
		return
	}
	var order models.AccrualResponse
	err = json.Unmarshal(resp.Body(), &order)
	if err != nil {
		logrus.Error("invalid accrual response body ", err)
		return
	}
	if order.Status == models.Registered {
		order.Status = models.Processing
	}
	err = a.storage.Update(context.Background(), order)
	if err != nil {
		logrus.Error(err)
	}
}
