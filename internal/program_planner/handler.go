package program_planner

import (
	"context"
	"fmt"
	"net/http"
	"portal/helper"
	"portal/pkg/constants"

	"github.com/gin-gonic/gin"
)

type ProgramPlanerHandler struct {
	ProgramPlanerService ProgramPlanerService
}

func NewProgramPlanerHandler(ProgramPlanerService ProgramPlanerService) *ProgramPlanerHandler {
	return &ProgramPlanerHandler{
		ProgramPlanerService: ProgramPlanerService,
	}
}

func (handler *ProgramPlanerHandler) CreateProgramPlaner(c *gin.Context) {

	var req CreateProgramPlanerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
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

	programPlanerID, err := handler.ProgramPlanerService.CreateProgramPlaner(ctx, &req, userID.(string))
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Create program planer successfully", programPlanerID)

}

func (handler *ProgramPlanerHandler) GetAllProgramPlaner(c *gin.Context) {

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	programPlaners, err := handler.ProgramPlanerService.GetAllProgramPlaner(ctx)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Get all program planer successfully", programPlaners)

}

func (handler *ProgramPlanerHandler) GetProgramPlaner(c *gin.Context) {

	id := c.Param("id")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	programPlaner, err := handler.ProgramPlanerService.GetProgramPlaner(ctx, id)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Get program planer successfully", programPlaner)

}

func (handler *ProgramPlanerHandler) UpdateProgramPlaner(c *gin.Context) {

	id := c.Param("id")

	var req UpdateProgramPlanerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := handler.ProgramPlanerService.UpdateProgramPlaner(ctx, &req, id)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Update program planer successfully", nil)

}

func (handler *ProgramPlanerHandler) DeleteProgramPlaner(c *gin.Context) {

	id := c.Param("id")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := handler.ProgramPlanerService.DeleteProgramPlaner(ctx, id)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Delete program planer successfully", nil)

}

func (handler *ProgramPlanerHandler) UpdateProgramPlanerWeek(c *gin.Context) {

	var req UpdateWeekProgramPlanerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	id := c.Param("id")

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := handler.ProgramPlanerService.UpdateProgramPlanerWeek(ctx, &req, id)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Create week program planer successfully", nil)

}