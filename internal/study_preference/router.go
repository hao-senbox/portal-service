package studypreference

import (
	"portal/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *StudyPreferenceHandler) {
	group := r.Group("/api/v1/study-preference", middleware.Secured())
	{
		group.POST("", handler.CreateStudyPreference)
		group.GET("/:id", handler.GetStudyPreferenceByID)
		group.PUT("/:id", handler.UpdateStudyPreference)
	}
}
