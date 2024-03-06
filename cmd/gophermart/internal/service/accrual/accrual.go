package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	ch      *WorkerTimeoutChan
}

func New(s storage.Order, c *config.Config) *Accrual {
	return &Accrual{
		storage: s,
		cfg:     c,
		client:  resty.New(),
		ch:      NewWorkerChanMap(c.RateLimit),
	}
}

func (a *Accrual) Run(ctx context.Context) {
	requestTicker := time.NewTicker(time.Duration(a.cfg.RequestInterval) * time.Second)
	defer requestTicker.Stop()

	tasks := make(chan string, 10)

	for i := 0; i < a.cfg.RateLimit; i++ {
		a.ch.Insert(i, make(chan int))
		go a.worker(ctx, tasks, i)
	}

	for t := range requestTicker.C {
		numbers, err := a.storage.GetProcessingOrders(ctx)
		if err != nil {
			logrus.Error("processing order collection error ", err, t)
			continue
		}

		for _, number := range numbers {
			tasks <- number
		}
	}
}

func (a *Accrual) worker(ctx context.Context, tasks <-chan string, id int) {
	for {
		select {
		case number, opened := <-tasks:
			if !opened {
				return //канал закрыт -> завершение работы воркера
			}
			a.processOrder(ctx, number)
		case timeout := <-a.ch.Get(id):
			time.Sleep(time.Duration(timeout) * time.Millisecond)
		}
	}
}

func (a *Accrual) processOrder(ctx context.Context, number string) {
	order, err := a.getOrderStatus(number)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = a.updateOrderStatus(ctx, *order)
	logrus.Error(err)
}

func (a *Accrual) getOrderStatus(number string) (*models.AccrualResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", a.cfg.AccrualSystemAddress, number)
	resp, err := a.client.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("error accrual requesting order status : %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		var order models.AccrualResponse
		err = json.Unmarshal(resp.Body(), &order)
		if err != nil {
			return nil, fmt.Errorf("invalid accrual response body : %w", err)
		}
		return &order, nil
	case http.StatusTooManyRequests:
		retryAfterHeader := resp.Header().Get("Retry-After")
		if retryAfterHeader != "" {
			timeout, err := strconv.Atoi(retryAfterHeader)
			if err != nil {
				a.ch.Broadcast(2000)
				logrus.Error(err)
			}
			a.ch.Broadcast(timeout)
		}
		fallthrough
	default:
		return nil, fmt.Errorf("received response from accrual with status code : %d", resp.StatusCode())
	}
}

func (a *Accrual) updateOrderStatus(ctx context.Context, o models.AccrualResponse) error {
	if o.Status == models.Registered {
		o.Status = models.Processing
	}
	return a.storage.Update(ctx, o)
}
