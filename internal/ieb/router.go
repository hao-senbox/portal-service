package ieb

import (
	"portal/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRouters(r *gin.Engine, IEBHandler *IEBHandler) {
	group := r.Group("/api/v1/ieb", middleware.Secured())
	{
		// group.GET("", IEBHandler.GetIEBs)
		group.GET("", IEBHandler.GetIEB)
		group.POST("", IEBHandler.CreateIEB)
		// group.PUT("/:id", IEBHandler.UpdateIEB)
		// group.DELETE("/:id", IEBHandler.DeleteIEB)
	}

}
