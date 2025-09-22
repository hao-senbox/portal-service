package portal

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *PortalHandlers) {
	portalGroup := r.Group("/api/v1/portal")
	{
		portalGroup.POST("/student", handler.CreateStudentActivity)
		portalGroup.GET("", handler.GetAllStudentActivity)
	}
}
