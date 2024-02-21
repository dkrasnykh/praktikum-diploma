package service

import (
	"context"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/storage"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/errs"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
)

type WithdrawService struct {
	storage storage.Withdraw
}

func NewWithdrawService(s storage.Withdraw) *WithdrawService {
	return &WithdrawService{storage: s}
}

func (s *WithdrawService) GetUserBalance(ctx context.Context, userId int) (*models.UserBalance, error) {
	return s.storage.GetUserBalance(ctx, userId)
}

func (s *WithdrawService) WithdrawReward(ctx context.Context, userId int, req models.WithdrawRequest) error {
	if err := validateOrderNumber(req.Order); err != nil {
		return err
	}
	balance, err := s.GetUserBalance(ctx, userId)
	if err != nil {
		return err
	}
	if int64(balance.Current*100) < int64(req.Sum*100) {
		return errs.ErrNoReward
	}
	return s.storage.WithdrawReward(ctx, userId, req)
}

func (s *WithdrawService) GetAllWithdrawals(ctx context.Context, userId int) ([]models.Withdraw, error) {
	return s.storage.GetAllWithdrawals(ctx, userId)
}
