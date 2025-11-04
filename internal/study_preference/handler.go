package studypreference

import (
	"context"
	"fmt"
	"net/http"
	"portal/helper"
	"portal/pkg/constants"

	"github.com/gin-gonic/gin"
)

type StudyPreferenceHandler struct {
	service StudyPreferenceService
}

func NewStudyPreferenceHandler(service StudyPreferenceService) *StudyPreferenceHandler {
	return &StudyPreferenceHandler{service: service}
}

func (h *StudyPreferenceHandler) CreateStudyPreference(c *gin.Context) {
	var req CreateStudyPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	studyPreferenceID, err := h.service.CreateStudyPreference(ctx, &req, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Create study preference successfully", studyPreferenceID)
}

func (h *StudyPreferenceHandler) GetStudyPreferencesByStudentID(c *gin.Context) {
	studentID := c.Query("student_id")
	if studentID == "" {
		helper.SendError(c, http.StatusBadRequest, nil, "Student ID is required")
		return
	}

	orgID := c.Query("organization_id")
	if orgID == "" {
		helper.SendError(c, http.StatusBadRequest, nil, "Organization ID is required")
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, http.StatusBadRequest, nil, "Token not found")
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	studyPreferences, err := h.service.GetStudyPreferencesByStudentID(ctx, studentID, orgID)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get study preferences by student ID successfully", studyPreferences)
}

// func (h *StudyPreferenceHandler) GetStudyPreferenceByID(c *gin.Context) {
	
// 	id := c.Param("id")

// 	if id == "" {
// 		helper.SendError(c, http.StatusBadRequest, nil, "Study preference ID is required")
// 		return
// 	}

// 	token, exists := c.Get(constants.Token)
// 	if !exists {
// 		helper.SendError(c, http.StatusBadRequest, nil, "Token not found")
// 		return
// 	}

// 	ctx := context.WithValue(c, constants.TokenKey, token)

// 	studyPreference, err := h.service.GetStudyPreferenceByID(ctx, id)
// 	if err != nil {
// 		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
// 		return
// 	}

// 	helper.SendSuccess(c, http.StatusOK, "Get study preference by ID successfully", studyPreference)


// }

func (h *StudyPreferenceHandler) UpdateStudyPreference(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, nil, "Study preference ID is required")
		return
	}

	var req UpdateStudyPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	if len(req.ParentSelections) == 0 {
		helper.SendError(c, http.StatusBadRequest, nil, "Parent selections cannot be empty")
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, http.StatusBadRequest, nil, "User ID not found")
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, http.StatusBadRequest, nil, "Token not found")
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.service.UpdateStudyPreference(ctx, id, &req, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "Failed to update study preference")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Update study preference successfully", gin.H{
		"id": id,
	})
}

func (h *StudyPreferenceHandler) GetStudyPreferenceStatistical(c *gin.Context) {

	orgID := c.Query("organization_id")
	if orgID == "" {
		helper.SendError(c, http.StatusBadRequest, nil, "Organization ID is required")
		return
	}

	studentID := c.Query("student_id")
	if studentID == "" {
		helper.SendError(c, http.StatusBadRequest, nil, "Student ID is required")
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, http.StatusBadRequest, nil, "Token not found")
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	statistical, err := h.service.GetStudyPreferenceStatistical(ctx, orgID, studentID)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "Failed to get study preference statistical")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get study preference statistical successfully", statistical)
}