package program_planner

import (
	"portal/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *ProgramPlanerHandler) {
	group := r.Group("/api/v1/program-planner", middleware.Secured())
	{
		group.GET("", handler.GetAllProgramPlaner)
		group.GET("/:id", handler.GetProgramPlaner)
		group.POST("", handler.CreateProgramPlaner)
		group.PUT("/:id", handler.UpdateProgramPlaner)
		group.DELETE("/:id", handler.DeleteProgramPlaner)

		// group.POST("/week/:id", handler.UpdateProgramPlanerWeek)
	}
}
