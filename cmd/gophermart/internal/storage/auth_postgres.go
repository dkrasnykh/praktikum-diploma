package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/errs"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
)

type AuthPostgres struct {
	db           *pgxpool.Pool
	queryTimeout time.Duration
}

func NewAuthPostgres(db *pgxpool.Pool, timeoutSec int) *AuthPostgres {
	return &AuthPostgres{
		db:           db,
		queryTimeout: time.Duration(timeoutSec) * time.Second,
	}
}

func (r *AuthPostgres) CreateUser(ctx context.Context, user models.User) (int, error) {
	newCtx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	var id int
	var err error
	row := r.db.QueryRow(newCtx, "INSERT INTO users (login, password_hash) values ($1, $2) RETURNING id",
		user.Login, user.Password)
	if err = row.Scan(&id); err != nil && isLoginExistError(err) {
		return 0, errs.ErrLoginAlreadyExist
	}
	return id, err
}

func (r *AuthPostgres) GetUser(ctx context.Context, username, password string) (*models.User, error) {
	newCtx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	var id int
	var err error
	row := r.db.QueryRow(newCtx, "select id from users where login = $1 and password_hash=$2", username, password)
	if err = row.Scan(&id); err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, errs.ErrInvalidLoginOrPassword
	}
	return &models.User{ID: id, Login: username, Password: password}, err
}

func isLoginExistError(err error) bool {
	pgxErr, ok := err.(*pgconn.PgError)
	if ok && pgxErr.Code == "23505" {
		return true
	}
	return false
}
