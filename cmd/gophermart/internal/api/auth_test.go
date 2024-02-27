package api

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/service"
	mock_service "github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/service/mocks"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/errs"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
)

func TestSignUp(t *testing.T) {
	type mockBehavior func(r *mock_service.MockAuthorization, user models.User)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            models.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			inputBody: `{"login": "login", "password": "password"}`,
			inputUser: models.User{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user models.User) {
				r.EXPECT().CreateUser(context.Background(), user).Return(1, nil)
				r.EXPECT().GenerateToken(context.Background(), "login", "password").Return("token", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"token":"token"}`,
		},
		{
			name:                 "invalid input",
			inputBody:            `{"login": "login"}`,
			inputUser:            models.User{},
			mockBehavior:         func(r *mock_service.MockAuthorization, user models.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "status conflict",
			inputBody: `{"login": "login", "password": "password"}`,
			inputUser: models.User{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user models.User) {
				r.EXPECT().CreateUser(context.Background(), user).Return(0, errs.ErrLoginAlreadyExist)
			},
			expectedStatusCode:   http.StatusConflict,
			expectedResponseBody: `{"message":"login already exist"}`,
		},
		{
			name:      "internal server error",
			inputBody: `{"login": "login", "password": "password"}`,
			inputUser: models.User{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user models.User) {
				r.EXPECT().CreateUser(context.Background(), user).Return(0, errors.New("test error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"test error"}`,
		},
		{
			name:      "status unauthorized",
			inputBody: `{"login": "login", "password": "password"}`,
			inputUser: models.User{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user models.User) {
				r.EXPECT().CreateUser(context.Background(), user).Return(1, nil)
				r.EXPECT().GenerateToken(context.Background(), "login", "password").Return("", errs.ErrInvalidLoginOrPassword)
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid login or password"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockAuthorization(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{Authorization: repo}
			handler := Handler{services}

			r := gin.New()
			r.POST("/api/user/register", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/register",
				bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestSignIn(t *testing.T) {
	type mockBehavior func(r *mock_service.MockAuthorization, user models.User)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            models.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			inputBody: `{"login": "login", "password": "password"}`,
			inputUser: models.User{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user models.User) {
				r.EXPECT().GenerateToken(context.Background(), "login", "password").Return("token", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"token":"token"}`,
		},
		{
			name:                 "invalid input",
			inputBody:            `{"login": "login"}`,
			inputUser:            models.User{},
			mockBehavior:         func(r *mock_service.MockAuthorization, user models.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"Key: 'User.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:      "internal server error",
			inputBody: `{"login": "login", "password": "password"}`,
			inputUser: models.User{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user models.User) {
				r.EXPECT().GenerateToken(context.Background(), "login", "password").Return("", errors.New("test error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"test error"}`,
		},
		{
			name:      "status unauthorized",
			inputBody: `{"login": "login", "password": "password"}`,
			inputUser: models.User{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user models.User) {
				r.EXPECT().GenerateToken(context.Background(), "login", "password").Return("", errs.ErrInvalidLoginOrPassword)
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid login or password"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockAuthorization(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{Authorization: repo}
			handler := Handler{services}

			r := gin.New()
			r.POST("/api/user/login", handler.signIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/login",
				bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
