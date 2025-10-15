package studyprogram

import (
	"context"
	"fmt"
	"net/http"
	"portal/helper"
	"portal/pkg/constants"

	"github.com/gin-gonic/gin"
)

type StudyProgramHandler struct {
	StudyProgramService StudyProgramService
}

func NewStudyProgramHandler(service StudyProgramService) *StudyProgramHandler {
	return &StudyProgramHandler{
		StudyProgramService: service,
	}
}

func (h *StudyProgramHandler) CreateStudyProgram(c *gin.Context) {

	var req CreateStudyProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
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

	studyProgramID, err := h.StudyProgramService.CreateStudyProgram(ctx, &req, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Create study program successfully", studyProgramID)

}

func (h *StudyProgramHandler) GetStudyPrograms(c *gin.Context) {

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.StudyProgramService.GetStudyPrograms(ctx)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get study programs successfully", data)
}

func (h *StudyProgramHandler) GetStudyProgram(c *gin.Context) {

	id := c.Param("id")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.StudyProgramService.GetStudyProgram(ctx, id)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get study program successfully", data)

}

func (h *StudyProgramHandler) UpdateStudyProgram(c *gin.Context) {

	id := c.Param("id")

	var req UpdateStudyProgramRequest

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

	err := h.StudyProgramService.UpdateStudyProgram(ctx, id, &req)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Update study program successfully", nil)

}

func (h *StudyProgramHandler) DeleteStudyProgram(c *gin.Context) {

	id := c.Param("id")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.StudyProgramService.DeleteStudyProgram(ctx, id)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Delete study program successfully", nil)
	
}