package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"portal/internal/models"
	"portal/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

	fmt.Printf("data : %v\n", req.Data)

	studentActivity := &models.StudentActivity{
		ID:           primitive.NewObjectID(),
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
				SessionID:    activity.ID.Hex(),
				TypeActivity: activity.TypeActivity,
				Data:         activity.Data,
				SubmittedAt:  activity.SubmittedAt,
				AssignedBy:   activity.AssignedBy,
				CreatedAt:    activity.CreatedAt,
				UpdatedAt:    activity.UpdatedAt,
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
	case "sleep_rest":
		return s.generateSleepRestStatistics(details)
	case "toileting":
		return s.generateToiletingStatistics(details)
	case "fluids":
		return s.generateFluidsStatistics(details)
	case "work":
		return s.generateWorkStatistics(details)
	case "exercise":
		return s.generateExerciseStatistics(details)
	case "social_play":
		return s.generateSocialPlayStatistics(details)
	case "food":
		return s.generateFoodStatistics(details)
	default:
		return nil
	}
}

func (s *portalService) generateFoodStatistics(details []models.ActivityDetail) map[string]interface{} {

	dishConsumption := make(map[string][]float64)

	for _, detail := range details {
		for _, d := range detail.Data {
			if strings.HasPrefix(d.Key, "what_is_the_name_of_the") && strings.HasSuffix(d.Key, "_dish") {
				dishName := d.Value

				dishNumber := strings.TrimPrefix(d.Key, "what_is_the_name_of_the_")
				dishNumber = strings.TrimSuffix(dishNumber, "_dish")
				consumtionKey := fmt.Sprintf("how_much_the_student_ate_the_%s_dish", dishNumber)

				for _, consumptionData := range detail.Data {
					if consumptionData.Key == consumtionKey {
						if consumption, err := strconv.ParseFloat(consumptionData.Value, 64); err == nil {
							dishConsumption[dishName] = append(dishConsumption[dishName], consumption)
						}
						break
					}
				}
			}
		}
	}
	
	var summary []map[string]interface{}

	for dishName, consumption := range dishConsumption {
		if len(consumption) > 0 {
			sum := 0.0
			for _, c := range consumption {
				sum += c
			}
			average := sum / float64(len(consumption))	
			summary = append(summary, map[string]interface{}{
				"dish_name": dishName,
				"total":     math.Round(average*100) / 100,
			})
		}
	}

	return map[string]interface{}{
		"dishes": summary,
	}

}

func (s *portalService) generateExerciseStatistics(details []models.ActivityDetail) map[string]interface{} {

	var total int
	for _, detail := range details {
		for _, d := range detail.Data {
			if d.Key == "duration_of_session" {
				if val, err := strconv.Atoi(d.Value); err == nil {
					total += val
				}
			}
		}
	}

	return map[string]interface{}{
		"total": s.parseSecondToHoursAndMinutes(total),
	}

}

func (s *portalService) generateSocialPlayStatistics(details []models.ActivityDetail) map[string]interface{} {
	return nil
}

func (s *portalService) generateWorkStatistics(details []models.ActivityDetail) map[string]interface{} {
	return nil
}

func (s *portalService) generateFluidsStatistics(details []models.ActivityDetail) map[string]interface{} {

	var totalWater int
	var totalJuice int
	var totalSmoothies int
	var totalMilk int
	var totalOther int

	for _, detail := range details {
		for _, d := range detail.Data {
			switch d.Key {
			case "water":
				if value, ok := s.parseFluidsValue(d.Value); ok {
					fmt.Printf("water: %v\n", value)
					totalWater += value
				}
			case "juice":
				if value, ok := s.parseFluidsValue(d.Value); ok {
					totalJuice += value
				}
			case "smoothie":
				if value, ok := s.parseFluidsValue(d.Value); ok {
					totalSmoothies += value
				}
			case "milk":
				if value, ok := s.parseFluidsValue(d.Value); ok {
					totalMilk += value
				}
			case "other_fluid":
				if value, ok := s.parseFluidsValue(d.Value); ok {
					totalOther += value
				} 
			}
		}
	}

	return map[string]interface{}{
		"water":     fmt.Sprintf("%dml", totalWater),
		"juice":     fmt.Sprintf("%dml", totalJuice),
		"smoothies": fmt.Sprintf("%dml", totalSmoothies),
		"milk":      fmt.Sprintf("%dml", totalMilk),
		"other":     fmt.Sprintf("%dml", totalOther),
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
		"number_1": number1,
		"number_2": number2,
		"number_3": number3,
		"max":      max,
	}
}

func (s *portalService) generateSleepRestStatistics(details []models.ActivityDetail) map[string]interface{} {

	var totalSleep int
	var totalRest int

	for _, detail := range details {
		for _, d := range detail.Data {
			switch d.Key {
			case "duration_of_sleep":
				if val, err := strconv.Atoi(d.Value); err == nil {
					totalSleep += val
				} else {
					fmt.Printf("invalid durian_of_sleep: %v\n", d.Value)
				}
			case "duration_of_rest":
				if val, err := strconv.Atoi(d.Value); err == nil {
					totalRest += val
				} else {
					fmt.Printf("invalid durian_of_rest: %v\n", d.Value)
				}
			}
		}
	}

	return map[string]interface{}{
		"sleep": s.parseSecondToHoursAndMinutes(totalSleep),
		"rest":  s.parseSecondToHoursAndMinutes(totalRest),
		"total": s.parseSecondToHoursAndMinutes(totalSleep + totalRest),
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

func (s *portalService) parseFluidsValue(input string) (int, bool) {
	type FluidDetails struct {
		Capacity     int `json:"capacity"`      // Dung tích nước
		ActualPoured int `json:"actual_poured"` // Dung tích nước rót thực tế  
		Consumed     int `json:"consumed"`      // Học sinh uống bao nhiêu
		Remaining    int `json:"remaining"`     // Còn lại bao nhiêu
	}

	if strings.Contains(input, "{") {
		
		var fluidDetails FluidDetails
		err := json.Unmarshal([]byte(input), &fluidDetails)
		if err != nil {
			return 0, false
		}

		return fluidDetails.Consumed, true

	}

	return 0, false
}
