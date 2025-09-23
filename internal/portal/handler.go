package portal

import (
	"net/http"
	"portal/helper"

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

	activities, err := h.portalService.GetAllStudentActivity(c.Request.Context(), studentId, date)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Student activity data retrieved successfully", activities)
}
