package drink

import (
	"portal/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, DrinkHandler *DrinkHandler) {
	group := r.Group("/api/v1/drink", middleware.Secured())
	{
		group.GET("", DrinkHandler.GetDrinks)
		group.GET("/:id", DrinkHandler.GetDrink)
		group.POST("", DrinkHandler.CreateDrink)
		// group.PUT("/:id", DrinkHandler.UpdateDrink)
		// group.DELETE("/:id", DrinkHandler.DeleteDrink)
		group.GET("/statistics", DrinkHandler.GetStatistics)
	}
}