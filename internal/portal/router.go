package portal

import (
	"portal/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *PortalHandlers) {
	portalGroup := r.Group("/api/v1/portal", middleware.Secured())
	{
		portalGroup.POST("/student", handler.CreateStudentActivity)
		portalGroup.GET("", handler.GetAllStudentActivity)
	}
}
