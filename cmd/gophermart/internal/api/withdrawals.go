package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/errs"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
)

func (h *Handler) getBalance(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var resp *models.UserBalance
	if resp, err = h.service.Withdraw.GetUserBalance(ctx, userID); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) withdrawReward(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var request models.WithdrawRequest
	if err = c.BindJSON(&request); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	if err = h.service.Withdraw.WithdrawReward(ctx, userID, request); err != nil {
		switch err {
		case errs.ErrInvalidOrderNumber:
			newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		case errs.ErrNoReward:
			newErrorResponse(c, http.StatusPaymentRequired, err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
}

func (h *Handler) getAllWithdrawals(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp, err := h.service.Withdraw.GetAllWithdrawals(ctx, userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if len(resp) == 0 {
		newErrorResponse(c, http.StatusNoContent, "result list is empty")
		return
	}
	c.JSON(http.StatusOK, resp)
}
