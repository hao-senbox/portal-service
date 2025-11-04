package selectoptions

import (
	"context"
	"errors"
	"net/http"

	"portal/helper"
	"portal/pkg/constants"

	"github.com/gin-gonic/gin"
)

type SelectOptionsHandler struct {
	service SelectOptionsService
}

func NewSelectOptionsHandler(service SelectOptionsService) *SelectOptionsHandler {
	return &SelectOptionsHandler{
		service: service,
	}
}

func (h *SelectOptionsHandler) CreateSelectOption(c *gin.Context) {

	var req CreateSelectOptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, errors.New("user ID not found"), "Unauthorized")
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, http.StatusBadRequest, nil, "Token not found")
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	id, err := h.service.CreateSelectOption(ctx, req, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, "Failed to create select option")
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Select option created successfully", id)

}
