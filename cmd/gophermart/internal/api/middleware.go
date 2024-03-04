package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	_authorizationHeader = "Authorization"
	_userKey             = "userID"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(_authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}
	header = strings.TrimPrefix(header, "Bearer ")
	if strings.TrimSpace(header) == "" {
		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}
	userID, err := h.service.Authorization.ParseToken(header)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(_userKey, userID)
}

func getUserID(c *gin.Context) (int, error) {
	id, ok := c.Get(_userKey)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
