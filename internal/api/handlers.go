package api

import (
	"net/http"
	"portal/internal/models"
	"portal/internal/service"
	"github.com/gin-gonic/gin"
)

type PortalHandlers struct {
	portalService service.PortalService	
}

func NewPortalHandlers (portalService service.PortalService) *PortalHandlers {
	return &PortalHandlers{
		portalService: portalService,
	}
}
func RegisterHandlers (router *gin.Engine, portalService service.PortalService) {

	handlers := NewPortalHandlers(portalService)

	portalGroup := router.Group("/api/v1/portal")
	{	
		portalGroup.POST("/student", handlers.CreateStudentActivity)
		portalGroup.GET("/:student_id", handlers.GetAllStudentActivity)
	}
}

func (h *PortalHandlers) CreateStudentActivity(c *gin.Context) {

	var req models.RequestStudentActivity

	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadGateway, err, models.ErrInvalidRequest)
		return 
	}

	err := h.portalService.CreateStudentActivity(c.Request.Context(), &req)
	if err != nil {
		SendError(c, http.StatusInternalServerError, err, models.ErrInvalidOperation)
		return 
	}

	SendSuccess(c, http.StatusOK, "Student activity created successfully", nil)
}

func (h *PortalHandlers) GetAllStudentActivity(c *gin.Context) {
	
	studentId := c.Param("student_id")
	date := c.Query("date")

	activities, err := h.portalService.GetAllStudentActivity(c.Request.Context(), studentId, date)
	if err != nil {
		SendError(c, http.StatusInternalServerError, err, models.ErrInvalidOperation)
		return
	}

	SendSuccess(c, http.StatusOK, "Student activity data retrieved successfully", activities)
}


