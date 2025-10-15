package teacherassign

import (
	"context"
	"fmt"
	"net/http"
	"portal/helper"
	"portal/pkg/constants"

	"github.com/gin-gonic/gin"
)

type TeacherAssignmentHandler struct {
	service TeacherAssignmentService
}

func NewTeacherAssignmentHandler(service TeacherAssignmentService) *TeacherAssignmentHandler {
	return &TeacherAssignmentHandler{
		service: service,
	}
}

func (h *TeacherAssignmentHandler) CreateTeacherAssignment(c *gin.Context) {

	var req CreateTeacherAssignmentRequest

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

	teacherAssignmentID, err := h.service.CreateTeacherAssignment(ctx, &req, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Create teacher assignment successfully", teacherAssignmentID)

}

func (h *TeacherAssignmentHandler) GetAllTeacherAssignment(c *gin.Context) {

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	teacherAssignments, err := h.service.GetAllTeacherAssignment(ctx)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get all teacher assignment successfully", teacherAssignments)

}

func (h *TeacherAssignmentHandler) GetTeacherAssignment(c *gin.Context) {

	id := c.Param("id")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	teacherAssignment, err := h.service.GetTeacherAssignment(ctx, id)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get teacher assignment successfully", teacherAssignment)

}

func (h *TeacherAssignmentHandler) UpdateTeacherAssignment(c *gin.Context) {

	id := c.Param("id")

	var req UpdateTeacherAssignmentRequest

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

	err := h.service.UpdateTeacherAssignment(ctx, id, &req)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Update teacher assignment successfully", nil)

}

func (h *TeacherAssignmentHandler) DeleteTeacherAssignment(c *gin.Context) {

	id := c.Param("id")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.service.DeleteTeacherAssignment(ctx, id)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Delete teacher assignment successfully", nil)
	
}