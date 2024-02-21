package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/errs"
	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/pkg/models"
)

func (h *Handler) signUp(c *gin.Context) {
	ctx := c.Request.Context()
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	id, err := h.service.Authorization.CreateUser(ctx, user)
	if err != nil {
		switch err {
		case errs.ErrLoginAlreadyExist:
			newErrorResponse(c, http.StatusConflict, err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) signIn(c *gin.Context) {
	ctx := c.Request.Context()
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.service.Authorization.GenerateToken(ctx, user.Login, user.Password)
	if err != nil {
		switch err {
		case errs.ErrInvalidLoginOrPassword:
			newErrorResponse(c, http.StatusUnauthorized, err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.Header(authorizationHeader, token)
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
