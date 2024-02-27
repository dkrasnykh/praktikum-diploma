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

func TestAdd(t *testing.T) {
	type mockBehavior func(r *mock_service.MockOrder, orderNumber string)

	tests := []struct {
		name               string
		inputBody          string
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:      "ok",
			inputBody: "2377225624",
			mockBehavior: func(r *mock_service.MockOrder, orderNumber string) {
				r.EXPECT().Add(context.Background(), 1, "2377225624").Return(nil)
			},
			expectedStatusCode: http.StatusAccepted,
		},
		{
			name:               "invalid input body",
			inputBody:          "",
			mockBehavior:       func(r *mock_service.MockOrder, orderNumber string) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:      "order already loaded by user",
			inputBody: "2377225624",
			mockBehavior: func(r *mock_service.MockOrder, orderNumber string) {
				r.EXPECT().Add(context.Background(), 1, "2377225624").Return(errs.ErrOrderExist)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:      "order already loaded by another user",
			inputBody: "2377225624",
			mockBehavior: func(r *mock_service.MockOrder, orderNumber string) {
				r.EXPECT().Add(context.Background(), 1, "2377225624").Return(errs.ErrUnreachableOrder)
			},
			expectedStatusCode: http.StatusConflict,
		},
		{
			name:      "invalid input number",
			inputBody: "5",
			mockBehavior: func(r *mock_service.MockOrder, orderNumber string) {
				r.EXPECT().Add(context.Background(), 1, "5").Return(errs.ErrInvalidOrderNumber)
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:      "internal error",
			inputBody: "2377225624",
			mockBehavior: func(r *mock_service.MockOrder, orderNumber string) {
				r.EXPECT().Add(context.Background(), 1, "2377225624").Return(errors.New("internal error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockOrder(c)
			test.mockBehavior(repo, test.inputBody)

			services := &service.Service{Order: repo}
			handler := Handler{services}

			r := gin.New()
			r.POST("/api/user/orders", func(c *gin.Context) {
				c.Set(userCtx, 1)
			}, handler.Add)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/orders",
				bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
		})
	}
}

func TestGetAll(t *testing.T) {
	type mockBehavior func(r *mock_service.MockOrder)

	accrual := float32(325.0)
	order := models.Order{Number: "2377225624", Status: models.Processed, Accrual: &accrual, UploadedAt: "2024-02-25T12:54:52Z"}

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehavior: func(r *mock_service.MockOrder) {
				r.EXPECT().GetAll(context.Background(), 1).Return([]models.Order{order}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[{"number":"2377225624","status":"PROCESSED","accrual":325,"uploaded_at":"2024-02-25T12:54:52Z"}]`,
		},
		{
			name: "empty response",
			mockBehavior: func(r *mock_service.MockOrder) {
				r.EXPECT().GetAll(context.Background(), 1).Return([]models.Order{}, nil)
			},
			expectedStatusCode:   http.StatusNoContent,
			expectedResponseBody: "",
		},
		{
			name: "internal error",
			mockBehavior: func(r *mock_service.MockOrder) {
				r.EXPECT().GetAll(context.Background(), 1).Return([]models.Order{}, errors.New("internal error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"internal error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockOrder(c)
			test.mockBehavior(repo)

			services := &service.Service{Order: repo}
			handler := Handler{services}

			r := gin.New()
			r.GET("/api/user/orders", func(c *gin.Context) {
				c.Set(userCtx, 1)
			}, handler.getAll)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/user/orders", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
			assert.True(t, strings.Contains(w.Header().Get(headers.ContentType), "application/json"))
		})
	}
}
