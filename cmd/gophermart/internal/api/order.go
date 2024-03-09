package api

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-http-utils/headers"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/errs"
)

func (h *Handler) Add(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	orderNumber, err := io.ReadAll(c.Request.Body)
	if err != nil || len(orderNumber) == 0 {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	if err = h.service.Add(ctx, userID, string(orderNumber)); err != nil {
		switch {
		case errors.Is(err, errs.ErrOrderExist):
			newErrorResponse(c, http.StatusOK, err.Error())
		case errors.Is(err, errs.ErrUnreachableOrder):
			newErrorResponse(c, http.StatusConflict, err.Error())
		case errors.Is(err, errs.ErrInvalidOrderNumber):
			newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.Status(http.StatusAccepted)
}

func (h *Handler) getAll(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp, err := h.service.GetAll(ctx, userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if len(resp) == 0 {
		newErrorResponse(c, http.StatusNoContent, "response is empty")
		return
	}
	c.Header(headers.ContentType, "application/json")
	c.JSON(http.StatusOK, resp)
}
