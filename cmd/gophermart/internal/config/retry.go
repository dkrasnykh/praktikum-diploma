package config

import (
	"errors"
	"time"

	"github.com/avast/retry-go"
	"github.com/sirupsen/logrus"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/errs"
)

type Retry struct {
	RetryAfterMillisec int
	Attempts           uint
}

func RetryNew() *Retry {
	return &Retry{RetryAfterMillisec: 3600, Attempts: 1}
}

func (c *Retry) DelayType(n uint, _ error, config *retry.Config) time.Duration {
	return time.Duration(c.RetryAfterMillisec) * time.Millisecond
}

func (c *Retry) OnRetry(n uint, err error) {
	logrus.Error(err, n)
}

func (c *Retry) IfStatusTooManyRequests(err error) bool {
	if errors.Is(err, errs.ErrTooManyRequests) {
		return true
	}
	return false
}
