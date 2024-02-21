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

func (s *WithdrawService) GetUserBalance(ctx context.Context, userID int) (*models.UserBalance, error) {
	return s.storage.GetUserBalance(ctx, userID)
}

func (s *WithdrawService) WithdrawReward(ctx context.Context, userID int, req models.WithdrawRequest) error {
	if err := validateOrderNumber(req.Order); err != nil {
		return err
	}
	balance, err := s.GetUserBalance(ctx, userID)
	if err != nil {
		return err
	}
	if int64(balance.Current*100) < int64(req.Sum*100) {
		return errs.ErrNoReward
	}
	return s.storage.WithdrawReward(ctx, userID, req)
}

func (s *WithdrawService) GetAllWithdrawals(ctx context.Context, userID int) ([]models.Withdraw, error) {
	return s.storage.GetAllWithdrawals(ctx, userID)
}
