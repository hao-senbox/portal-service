package portal

import (
	"context"
	"fmt"
	"net/http"
	"portal/helper"
	"portal/pkg/constants"

	"github.com/gin-gonic/gin"
)

type PortalHandlers struct {
	portalService PortalService
}

func NewPortalHandlers(portalService PortalService) *PortalHandlers {
	return &PortalHandlers{
		portalService: portalService,
	}
}
func (h *PortalHandlers) CreateStudentActivity(c *gin.Context) {

	var req RequestStudentActivity

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadGateway, err, helper.ErrInvalidRequest)
		return
	}

	err := h.portalService.CreateStudentActivity(c.Request.Context(), &req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Student activity created successfully", nil)
}

func (h *PortalHandlers) GetAllStudentActivity(c *gin.Context) {

	studentId := c.Query("student_id")
	date := c.Query("date")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	activities, err := h.portalService.GetAllStudentActivity(ctx, studentId, date)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Student activity data retrieved successfully", activities)
}
