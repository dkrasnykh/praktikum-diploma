package service

import (
	"context"

	"github.com/ShiraazMoollatjie/goluhn"

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
	current, withdrawn, err := s.storage.GetUserBalance(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &models.UserBalance{Current: float32(current) / 100, Withdrawn: float32(withdrawn) / 100}, nil
}

func (s *WithdrawService) WithdrawReward(ctx context.Context, userID int, req models.WithdrawRequest) error {
	if err := goluhn.Validate(req.Order); err != nil {
		return errs.ErrInvalidOrderNumber
	}
	balance, err := s.GetUserBalance(ctx, userID)
	if err != nil {
		return err
	}
	if balance.Current < req.Sum {
		return errs.ErrNoReward
	}
	return s.storage.WithdrawReward(ctx, userID, req)
}

func (s *WithdrawService) GetAllWithdrawals(ctx context.Context, userID int) ([]models.Withdraw, error) {
	return s.storage.GetAllWithdrawals(ctx, userID)
}
