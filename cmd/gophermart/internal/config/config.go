package config

import (
	"flag"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURL          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	RequestInterval      int    `env:"REQUEST_INTERVAL"`
	RateLimit            int    `env:"RATE_LIMIT"`
	QueryTimeout         int    `env:"QUERY_TIMEOUT"`
	ConnectTimeout       int    `env:"CONNECT_TIMEOUT"`
}

func New() (*Config, error) {
	var c Config
	flag.StringVar(&c.RunAddress, "a", ":8081", "address and port to run service")
	flag.StringVar(&c.DatabaseURL, "d", "postgres://postgres:postgres@localhost:5432/gofermart?sslmode=disable", "url for database connection")
	flag.StringVar(&c.AccrualSystemAddress, "r", ":8080", "accrual system address")
	flag.IntVar(&c.RequestInterval, "p", 1, "frequency of data requests to accrual service")
	flag.IntVar(&c.QueryTimeout, "q", 5, "database query timeout")
	flag.IntVar(&c.ConnectTimeout, "c", 10, "database connection timeout")
	flag.Parse()

	err := env.Parse(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
