package config

import (
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"github.com/sirupsen/logrus"
)

const Attempts uint = 100

func DelayType(n uint, _ error, config *retry.Config) time.Duration {
	return 1 * time.Second
}

func OnRetry(n uint, err error) {
	logrus.Error(fmt.Sprintf(`%d %s`, n, err.Error()))
}
