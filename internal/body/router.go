package body

import (
	"portal/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, BodyHandler *BodyHandler) {
	group := r.Group("/api/v1/body", middleware.Secured())
	{
		group.GET("", BodyHandler.GetCheckIns)
		// group.GET("/:id", BodyHandler.GetCheckIn)
		group.POST("", BodyHandler.CreateCheckIn)
		// group.PUT("/:id", BodyHandler.UpdateCheckIn)
		// group.DELETE("/:id", BodyHandler.DeleteCheckIn)
	}
}
