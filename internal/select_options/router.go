package selectoptions

import (
	"portal/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *SelectOptionsHandler) {
	group := r.Group("/api/v1/select-options", middleware.Secured())
	{
		group.POST("", handler.CreateSelectOption)
	}
}
