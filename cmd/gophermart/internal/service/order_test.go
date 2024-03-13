package service

import (
	"context"
	mock_storage "github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/storage/mocks"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAll(t *testing.T) {
	type mockBehavior func(r *mock_storage.MockOrder, userID int)

	c := gomock.NewController(t)
	defer c.Finish()
	value := float32(325.0)
	order := models.Order{Number: "2377225624", Status: models.Processed, Accrual: &value, UploadedAt: "2024-02-25T12:54:52Z"}

	tests := []struct {
		name          string
		ctx           context.Context
		inputUserID   int
		mockBehavior  mockBehavior
		expectedList  []models.Order
		expectedError error
	}{
		{
			name:        "ok",
			ctx:         context.Background(),
			inputUserID: 1,
			mockBehavior: func(r *mock_storage.MockOrder, userID int) {
				r.EXPECT().GetAll(context.Background(), userID).Return([]models.Order{order}, nil)
			},
			expectedList:  []models.Order{order},
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repoOrder := mock_storage.NewMockOrder(c)
			test.mockBehavior(repoOrder, test.inputUserID)

			services := &Service{Order: NewOrderService(repoOrder, nil, nil)}

			list, err := services.GetAll(test.ctx, test.inputUserID)

			assert.Equal(t, list, test.expectedList)
			assert.Equal(t, err, test.expectedError)
		})
	}
}
