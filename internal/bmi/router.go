package bmi

import (
	"portal/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, BMIHandler *BMIHandler) {
	group := r.Group("/api/v1/bmi", middleware.Secured())
	{
		group.GET("/", BMIHandler.GetBMIs)
		group.GET("/:id", BMIHandler.GetBMI)
		group.POST("", BMIHandler.CreateBMI)
		// group.PUT("/:id", BMIHandler.UpdateBMI)
		// group.DELETE("/:id", BMIHandler.DeleteBMI)
	}
}