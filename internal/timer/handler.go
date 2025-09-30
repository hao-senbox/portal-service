package timer

import (
	"context"
	"fmt"
	"portal/helper"
	"portal/pkg/constants"

	"github.com/gin-gonic/gin"
)

type TimerHandler struct {
	TimeService TimerService
}

func NewTimerHandler(TimeService TimerService) *TimerHandler {
	return &TimerHandler{
		TimeService: TimeService,
	}
}

func (h *TimerHandler) CreateTimer(c *gin.Context) {

	var req CreateTimerRequest
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

	timerID, err := h.TimeService.CreateTimer(ctx, &req, userID.(string))
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Create timer successfully", timerID)
	
}

func (h *TimerHandler) GetTimers(c *gin.Context) {

	student_id := c.Query("student")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	timers, err := h.TimeService.GetTimers(ctx, student_id)

	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Get timers successfully", timers)

}

func (h *TimerHandler) CreateIsTime(c *gin.Context) {

	var req CreateIsTimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.TimeService.CreateIsTime(ctx, &req, userID.(string))
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Create is time successfully", nil)

}

func (h *TimerHandler) GetIsTimes(c *gin.Context) {

	studentID := c.Query("student")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	isTimes, err := h.TimeService.GetIsTimes(ctx, studentID)

	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Get is times successfully", isTimes)

}