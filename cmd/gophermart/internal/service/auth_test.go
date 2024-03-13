package service

import (
	"context"
	mock_storage "github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/storage/mocks"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUser(t *testing.T) {
	type mockBehavior func(r *mock_storage.MockAuthorization, user models.User)

	c := gomock.NewController(t)
	defer c.Finish()

	tests := []struct {
		name          string
		ctx           context.Context
		inputUser     models.User
		mockBehavior  mockBehavior
		expectedID    int
		expectedError error
	}{
		{
			name: "ok",
			ctx:  context.Background(),
			inputUser: models.User{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(r *mock_storage.MockAuthorization, user models.User) {
				user.Password = "63687364616a6376687364626a68636a6462685baa61e4c9b93f3f0682250b6cf8331b7ee68fd8"
				r.EXPECT().CreateUser(context.Background(), user).Return(1, nil)
			},
			expectedID:    1,
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repoAuth := mock_storage.NewMockAuthorization(c)
			test.mockBehavior(repoAuth, test.inputUser)

			services := &Service{Authorization: NewAuthService(repoAuth)}

			id, err := services.CreateUser(test.ctx, test.inputUser)

			assert.Equal(t, id, test.expectedID)
			assert.Equal(t, err, test.expectedError)
		})
	}
}
