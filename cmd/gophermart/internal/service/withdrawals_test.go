package service

import (
	"context"
	mock_storage "github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/storage/mocks"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllWithdrawals(t *testing.T) {
	type mockBehavior func(r *mock_storage.MockWithdraw, userID int)

	c := gomock.NewController(t)
	defer c.Finish()

	withdraw := models.Withdraw{Order: "2377225624", Sum: 250, ProcessedAt: "2024-02-25T12:54:52Z"}

	tests := []struct {
		name          string
		ctx           context.Context
		inputUserID   int
		mockBehavior  mockBehavior
		expectedList  []models.Withdraw
		expectedError error
	}{
		{
			name:        "ok",
			ctx:         context.Background(),
			inputUserID: 1,
			mockBehavior: func(r *mock_storage.MockWithdraw, userID int) {
				r.EXPECT().GetAllWithdrawals(context.Background(), userID).Return([]models.Withdraw{withdraw}, nil)
			},
			expectedList:  []models.Withdraw{withdraw},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_storage.NewMockWithdraw(c)
			test.mockBehavior(repo, test.inputUserID)

			services := &Service{Withdraw: NewWithdrawService(repo, nil)}

			list, err := services.GetAllWithdrawals(test.ctx, test.inputUserID)

			assert.Equal(t, list, test.expectedList)
			assert.Equal(t, err, test.expectedError)
		})
	}
}
