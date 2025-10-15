package teacherassign

import (
	"portal/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *TeacherAssignmentHandler) {
	group := r.Group("/api/v1/teacher-assignment", middleware.Secured())
	{
		group.GET("", handler.GetAllTeacherAssignment)
		group.GET("/:id", handler.GetTeacherAssignment)
		group.POST("", handler.CreateTeacherAssignment)
		group.PUT("/:id", handler.UpdateTeacherAssignment)
		group.DELETE("/:id", handler.DeleteTeacherAssignment)
	}
}
