package service

import (
	"context"
	"fmt"
	"portal/internal/models"
	"portal/internal/repository"
	"strconv"
	"strings"
	"time"
)

type PortalService interface {
	CreateStudentActivity(ctx context.Context, req *models.RequestStudentActivity) error
	GetAllStudentActivity(ctx context.Context, studentID string, date string) ([]*models.StudentDailyActivities, error)
}
type portalService struct {
	repoPortal repository.PortalRepository
}

func NewPortalService(repo repository.PortalRepository) PortalService {
	return &portalService{
		repoPortal: repo,
	}
}

func (s *portalService) CreateStudentActivity(ctx context.Context, req *models.RequestStudentActivity) error {

	if req.StudentID == "" {
		return fmt.Errorf("student ID cannot be empty")
	}

	if req.TypeActivity == "" {
		return fmt.Errorf("type activity cannot be empty")
	}

	if req.Date == "" {
		return fmt.Errorf("date cannot be empty")
	}

	dateParse, err := time.Parse("2006-01-02T15:04:05Z07:00", req.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	if req.Data == nil {
		return fmt.Errorf("data cannot be empty")
	}

	if req.AssignedBy == "" {
		return fmt.Errorf("assigned by cannot be empty")
	}

	studentActivity := &models.StudentActivity{
		StudentID:    req.StudentID,
		TypeActivity: req.TypeActivity,
		Date:         dateParse,
		Data:         req.Data,
		SubmittedAt:  time.Now(),
		AssignedBy:   req.AssignedBy,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.repoPortal.CreateStudentActivity(ctx, studentActivity)
	if err != nil {
		return fmt.Errorf("failed to create student activity: %w", err)
	}

	return nil
}

func (s *portalService) GetAllStudentActivity(ctx context.Context, studentID string, date string) ([]*models.StudentDailyActivities, error) {

	var parsedDate *time.Time

	if date != "" {
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %w", err)
		}
		parsedDate = &t
	}

	activities, err := s.repoPortal.GetAllStudentActivity(ctx, studentID, parsedDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get student activity: %w", err)
	}

	transformData := s.transformStudentActivities(activities)

	return transformData, nil
}

func (s *portalService) transformStudentActivities(rawActivities []*models.StudentActivity) []*models.StudentDailyActivities {

	groupedData := make(map[string]map[string][]*models.StudentActivity)

	for _, activity := range rawActivities {
		studentID := activity.StudentID
		date := activity.Date.Format("2006-01-02")

		if groupedData[studentID] == nil {
			groupedData[studentID] = make(map[string][]*models.StudentActivity)
		}

		groupedData[studentID][date] = append(groupedData[studentID][date], activity)
	}

	var result []*models.StudentDailyActivities

	for studentID, dateMap := range groupedData {
		for date, activities := range dateMap {
			dailyActivity := &models.StudentDailyActivities{
				StudentID:  studentID,
				Date:       date,
				Activities: s.groupActivitiesByType(activities),
			}
			result = append(result, dailyActivity)
		}
	}

	return result
}

func (s *portalService) groupActivitiesByType(activities []*models.StudentActivity) []models.ActivitySummary {

	typeGroups := make(map[string][]*models.StudentActivity)

	for _, activity := range activities {
		typeGroups[activity.TypeActivity] = append(typeGroups[activity.TypeActivity], activity)
	}

	var result []models.ActivitySummary

	for typeActivity, typeActivities := range typeGroups {
		var details []models.ActivityDetail
		for _, activity := range typeActivities {
			details = append(details, models.ActivityDetail{
				Data:        activity.Data,
				SubmittedAt: activity.SubmittedAt,
				AssignedBy:  activity.AssignedBy,
				CreatedAt:   activity.CreatedAt,
				UpdatedAt:   activity.UpdatedAt,
			})
		}

		summary := models.ActivitySummary{
			TypeActivity: typeActivity,
			Summary: models.ActivitySummaryData{
				TotalSessions: len(details),
				Statistics:    s.generateStatistics(typeActivity, details),
			},
			Details: details,
		}

		result = append(result, summary)
	}

	return result

}

func (s *portalService) generateStatistics(typeActivity string, details []models.ActivityDetail) map[string]interface{} {
	switch typeActivity {
	case "sleep-rest":
		return s.generateSleepRestStatistics(details)
	case "toileting":
		return s.generateToiletingStatistics(details)
	default:
		return nil
	}
}

func (s *portalService) generateToiletingStatistics(details []models.ActivityDetail) map[string]interface{} {

	var number1 int
	var number2 int
	var number3 int
	var max int

	for _, detail := range details {
		for _, d := range detail.Data {
			switch d.Key {
			case "number_1":
				value := strings.ToLower(d.Value)
				if strings.Contains(value, "nothing") {
					number1 -= 1
				} else if strings.Contains(value, "small") {
					number1 += 1
				} else if strings.Contains(value, "big") {
					number1 += 1
				}
			case "number_2":
				value := strings.ToLower(d.Value)
				if strings.Contains(value, "nothing") {
					if number2 == 0 {
						number2 = 0
					} else {
						number2 -= 1
					}
				} else if strings.Contains(value, "small") {
					number2 += 1
				} else if strings.Contains(value, "big") {
					number2 += 1
				}
			case "number_3":
				value := strings.ToLower(d.Value)
				if strings.Contains(value, "nothing") {
					if number3 == 0 {
						number3 = 0
					} else {
						number3 -= 1
					}
				} else if strings.Contains(value, "small") {
					number3 += 1
				} else if strings.Contains(value, "big") {
					number3 += 1
				}
			}
		}
	}

	if number1 < 0 {
		number1 = 0
	}

	if number2 < 0 {
		number2 = 0
	}

	if number3 < 0 {
		number3 = 0
	}

	max = number1 + number2 + number3
	return map[string]interface{}{
		"total_number_1": number1,
		"total_number_2": number2,
		"total_number_3": number3,
		"max":            max,
	}
}

func (s *portalService) generateSleepRestStatistics(details []models.ActivityDetail) map[string]interface{} {

	var totalSleep int
	var totalRest int

	for _, detail := range details {
		for _, d := range detail.Data {
			switch d.Key {
			case "durian_of_sleep":
				if val, err := strconv.Atoi(d.Value); err == nil {
					totalSleep += val
				} else {
					fmt.Printf("invalid durian_of_sleep: %v\n", d.Value)
				}
			case "durian_of_rest":
				if val, err := strconv.Atoi(d.Value); err == nil {
					totalRest += val
				} else {
					fmt.Printf("invalid durian_of_rest: %v\n", d.Value)
				}
			}
		}
	}

	return map[string]interface{}{
		"total_sleep": s.parseSecondToHoursAndMinutes(totalSleep),
		"total_rest":  s.parseSecondToHoursAndMinutes(totalRest),
	}

}

func (s *portalService) parseSecondToHoursAndMinutes(seconds int) string {

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60

	var result string
	if hours > 0 {
		result += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 || hours == 0 {
		result += fmt.Sprintf("%dm", minutes)
	}

	return result

}
