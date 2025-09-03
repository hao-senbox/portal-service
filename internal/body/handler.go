package body

import (
	"context"
	"fmt"
	"portal/helper"
	"portal/pkg/constants"

	"github.com/gin-gonic/gin"
)

type BodyHandler struct {
	BodyService BodyService
}

func NewBodyHandler(bodyService BodyService) *BodyHandler {
	return &BodyHandler{
		BodyService: bodyService,
	}
}

func (h *BodyHandler) CreateCheckIn(c *gin.Context) {

	var req CreateCheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	checkInID, err := h.BodyService.CreateCheckIn(ctx, &req, userID.(string))
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Create check in successfully", checkInID)
	
}