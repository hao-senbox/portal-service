package timer

import (
	"portal/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, TimerHandler *TimerHandler) {
	group := r.Group("/api/v1/timer", middleware.Secured())
	{
		group.GET("", TimerHandler.GetTimers)
		group.POST("", TimerHandler.CreateTimer)


		group.POST("/is-time", TimerHandler.CreateIsTime)
		group.GET("/is-time", TimerHandler.GetIsTimes)
	}
}