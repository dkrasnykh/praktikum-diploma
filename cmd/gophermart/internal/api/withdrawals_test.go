package api

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-http-utils/headers"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/service"
	mock_service "github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/service/mocks"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/errs"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
)

func TestGetBalance(t *testing.T) {
	type mockBehavior func(r *mock_service.MockWithdraw)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehavior: func(r *mock_service.MockWithdraw) {
				r.EXPECT().GetUserBalance(context.Background(), 1).Return(&models.UserBalance{Current: 325.0, Withdrawn: 0}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"current":325,"withdrawn":0}`,
		},
		{
			name: "internal error",
			mockBehavior: func(r *mock_service.MockWithdraw) {
				r.EXPECT().GetUserBalance(context.Background(), 1).Return(nil, errors.New("internal error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"internal error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockWithdraw(c)
			test.mockBehavior(repo)

			services := &service.Service{Withdraw: repo}
			handler := Handler{services}

			r := gin.New()
			r.GET("/api/user/balance", func(c *gin.Context) {
				c.Set(_userKey, 1)
			}, handler.getBalance)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/user/balance", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
			assert.True(t, strings.Contains(w.Header().Get(headers.ContentType), "application/json"))
		})
	}
}

func TestWithdrawReward(t *testing.T) {
	type mockBehavior func(r *mock_service.MockWithdraw, request models.WithdrawRequest)

	tests := []struct {
		name               string
		inputBody          string
		request            models.WithdrawRequest
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:      "ok",
			inputBody: `{"order": "2377225624", "sum": 250}`,
			request:   models.WithdrawRequest{Order: "2377225624", Sum: 250},
			mockBehavior: func(r *mock_service.MockWithdraw, request models.WithdrawRequest) {
				r.EXPECT().WithdrawReward(context.Background(), 1, request).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "invalid input body",
			inputBody:          "",
			mockBehavior:       func(r *mock_service.MockWithdraw, request models.WithdrawRequest) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:      "invalid order number",
			inputBody: `{"order": "11", "sum": 250}`,
			request:   models.WithdrawRequest{Order: "11", Sum: 250},
			mockBehavior: func(r *mock_service.MockWithdraw, request models.WithdrawRequest) {
				r.EXPECT().WithdrawReward(context.Background(), 1, request).Return(errs.ErrInvalidOrderNumber)
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:      "not enough reward",
			inputBody: `{"order": "2377225624", "sum": 250}`,
			request:   models.WithdrawRequest{Order: "2377225624", Sum: 250},
			mockBehavior: func(r *mock_service.MockWithdraw, request models.WithdrawRequest) {
				r.EXPECT().WithdrawReward(context.Background(), 1, request).Return(errs.ErrNoReward)
			},
			expectedStatusCode: http.StatusPaymentRequired,
		},
		{
			name:      "internal error",
			inputBody: `{"order": "2377225624", "sum": 250}`,
			request:   models.WithdrawRequest{Order: "2377225624", Sum: 250},
			mockBehavior: func(r *mock_service.MockWithdraw, request models.WithdrawRequest) {
				r.EXPECT().WithdrawReward(context.Background(), 1, request).Return(errors.New("internal error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockWithdraw(c)
			test.mockBehavior(repo, test.request)

			services := &service.Service{Withdraw: repo}
			handler := Handler{services}

			r := gin.New()
			r.POST("/api/user/balance/withdraw", func(c *gin.Context) {
				c.Set(_userKey, 1)
			}, handler.withdrawReward)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/balance/withdraw",
				bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
		})
	}
}

func TestGetAllWithdrawals(t *testing.T) {
	type mockBehavior func(r *mock_service.MockWithdraw)

	withdraw := models.Withdraw{Order: "2377225624", Sum: 250.5, ProcessedAt: "2024-02-25T12:54:52Z"}
	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehavior: func(r *mock_service.MockWithdraw) {
				r.EXPECT().GetAllWithdrawals(context.Background(), 1).Return([]models.Withdraw{withdraw}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[{"order":"2377225624","sum":250.5,"processed_at":"2024-02-25T12:54:52Z"}]`,
		},
		{
			name: "result is empty",
			mockBehavior: func(r *mock_service.MockWithdraw) {
				r.EXPECT().GetAllWithdrawals(context.Background(), 1).Return([]models.Withdraw{}, nil)
			},
			expectedStatusCode:   http.StatusNoContent,
			expectedResponseBody: "",
		},
		{
			name: "internal error",
			mockBehavior: func(r *mock_service.MockWithdraw) {
				r.EXPECT().GetAllWithdrawals(context.Background(), 1).Return(nil, errors.New("internal error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"internal error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockWithdraw(c)
			test.mockBehavior(repo)

			services := &service.Service{Withdraw: repo}
			handler := Handler{services}

			r := gin.New()
			r.GET("/api/user/withdrawals", func(c *gin.Context) {
				c.Set(_userKey, 1)
			}, handler.getAllWithdrawals)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/user/withdrawals", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
			assert.True(t, strings.Contains(w.Header().Get(headers.ContentType), "application/json"))
		})
	}
}
