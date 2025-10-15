package studyprogram

import (
	"portal/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, studyProgramHandler *StudyProgramHandler) {
	group := r.Group("/api/v1/study-program", middleware.Secured())
	{
		group.GET("", studyProgramHandler.GetStudyPrograms)
		group.GET("/:id", studyProgramHandler.GetStudyProgram)
		group.POST("", studyProgramHandler.CreateStudyProgram)
		group.PUT("/:id", studyProgramHandler.UpdateStudyProgram)
		group.DELETE("/:id", studyProgramHandler.DeleteStudyProgram)
	}
}
