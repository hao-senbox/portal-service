package portal

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	attendancePkg "portal/internal/attendance"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PortalService interface {
	CreateStudentActivity(ctx context.Context, req *RequestStudentActivity) error
	GetAllStudentActivity(ctx context.Context, studentID string, date string) ([]*StudentDailyActivities, error)
}
type portalService struct {
	repoPortal        PortalRepository
	attendanceService attendancePkg.AttendanceService
}

func NewPortalService(
	repo PortalRepository,
	attendanceService attendancePkg.AttendanceService,
) PortalService {
	return &portalService{
		repoPortal:        repo,
		attendanceService: attendanceService,
	}
}

func (s *portalService) CreateStudentActivity(ctx context.Context, req *RequestStudentActivity) error {

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

	studentActivity := &StudentActivity{
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

func (s *portalService) GetAllStudentActivity(ctx context.Context, studentID string, date string) ([]*StudentDailyActivities, error) {

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

	attendanceInfo, err := s.attendanceService.GetAttendanceInfor(ctx, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendance information: %w", err)
	}

	transformData := s.transformStudentActivities(activities, attendanceInfo, date)

	return transformData, nil
}

func (s *portalService) transformStudentActivities(rawActivities []*StudentActivity, attendanceInfo []*attendancePkg.AttendanceUserInfo, filterDate string) []*StudentDailyActivities {

	// Group activities by student and date
	groupedActivities := make(map[string]map[string][]*StudentActivity)
	studentID := ""

	for _, activity := range rawActivities {
		studentID = activity.StudentID
		date := activity.Date.Format("2006-01-02")

		if groupedActivities[studentID] == nil {
			groupedActivities[studentID] = make(map[string][]*StudentActivity)
		}

		groupedActivities[studentID][date] = append(groupedActivities[studentID][date], activity)
	}

	// Group attendance by student and date
	groupedAttendance := make(map[string]map[string][]attendancePkg.AttendanceUserInfo)
	for _, attendance := range attendanceInfo {
		if studentID == "" {
			studentID = attendance.StudentID
		}

		// Parse attendance date (format: "2006-01-02T15:04:05Z")
		attendanceDate, err := time.Parse("2006-01-02T15:04:05Z", attendance.Date)
		if err != nil {
			continue
		}
		date := attendanceDate.Format("2006-01-02")

		// Filter by date if provided
		if filterDate != "" && date != filterDate {
			continue
		}

		if groupedAttendance[studentID] == nil {
			groupedAttendance[studentID] = make(map[string][]attendancePkg.AttendanceUserInfo)
		}

		groupedAttendance[studentID][date] = append(groupedAttendance[studentID][date], *attendance)
	}

	var result []*StudentDailyActivities

	// Collect all unique dates from both activities and attendance
	allDates := make(map[string]map[string]bool)
	for studentID, dateMap := range groupedActivities {
		if allDates[studentID] == nil {
			allDates[studentID] = make(map[string]bool)
		}
		for date := range dateMap {
			allDates[studentID][date] = true
		}
	}
	for studentID, dateMap := range groupedAttendance {
		if allDates[studentID] == nil {
			allDates[studentID] = make(map[string]bool)
		}
		for date := range dateMap {
			allDates[studentID][date] = true
		}
	}

	// Create daily activities for each student and date combination
	for studentID, dateMap := range allDates {
		for date := range dateMap {
			activities := groupedActivities[studentID][date]

			dailyActivity := &StudentDailyActivities{
				StudentID:  studentID,
				Date:       date,
				Activities: s.groupActivitiesByType(activities, groupedAttendance[studentID][date]),
			}
			result = append(result, dailyActivity)
		}
	}

	return result
}

func (s *portalService) groupActivitiesByType(activities []*StudentActivity, attendances []attendancePkg.AttendanceUserInfo) []ActivitySummary {

	typeGroups := make(map[string][]*StudentActivity)

	for _, activity := range activities {
		typeGroups[activity.TypeActivity] = append(typeGroups[activity.TypeActivity], activity)
	}

	var result []ActivitySummary

	for typeActivity, typeActivities := range typeGroups {
		var details []ActivityDetail
		for _, activity := range typeActivities {
			details = append(details, ActivityDetail{
				SessionID:    activity.ID.Hex(),
				TypeActivity: activity.TypeActivity,
				Data:         activity.Data,
				SubmittedAt:  activity.SubmittedAt,
				AssignedBy:   activity.AssignedBy,
				CreatedAt:    activity.CreatedAt,
				UpdatedAt:    activity.UpdatedAt,
			})
		}

		var statistics interface{}
		if typeActivity == "fluids" {
			statistics = s.generateFluidsStatisticsOrdered(details)
		} else {
			statistics = s.generateStatistics(typeActivity, details)
		}

		summary := ActivitySummary{
			TypeActivity: typeActivity,
			Summary: ActivitySummaryData{
				TotalSessions: len(details),
				Statistics:    statistics,
			},
			Details: details,
		}

		result = append(result, summary)
	}

	// Add attendance as a separate activity type
	if len(attendances) > 0 {
		attendanceDetails := s.createAttendanceDetails(attendances)
		if len(attendanceDetails) > 0 {
			attendanceStatistics := s.generateAttendanceStatistics(attendanceDetails)

			attendanceSummary := ActivitySummary{
				TypeActivity: "attendance",
				Summary: ActivitySummaryData{
					TotalSessions: len(attendanceDetails),
					Statistics:    attendanceStatistics,
				},
				Details: attendanceDetails,
			}

			result = append(result, attendanceSummary)
		}
	}

	return result

}

func (s *portalService) generateStatistics(typeActivity string, details []ActivityDetail) map[string]interface{} {
	switch typeActivity {
	case "sleep_rest":
		return s.generateSleepRestStatistics(details)
	case "toileting":
		return s.generateToiletingStatistics(details)
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

func (s *portalService) generateFoodStatistics(details []ActivityDetail) map[string]interface{} {

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

func (s *portalService) generateExerciseStatistics(details []ActivityDetail) map[string]interface{} {

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

func (s *portalService) generateSocialPlayStatistics(details []ActivityDetail) map[string]interface{} {
	return nil
}

func (s *portalService) generateWorkStatistics(details []ActivityDetail) map[string]interface{} {
	return nil
}

func (s *portalService) generateToiletingStatistics(details []ActivityDetail) map[string]interface{} {

	var number1 int
	var number2 int
	var number3 int
	var max int

	for _, detail := range details {
		for _, d := range detail.Data {
			switch d.Key {
			case "number_1":
				value := strings.ToLower(d.Value)
				if strings.Contains(value, "nothing") || strings.Contains(strings.ToLower(value), "independent") {
					if number1 == 0 {
						number1 = 0
					} else {
						number1 -= 1
					}
				} else if strings.Contains(value, "small") {
					number1 += 1
				} else if strings.Contains(value, "big") {
					number1 += 1
				}
			case "number_2":
				value := strings.ToLower(d.Value)
				if strings.Contains(value, "nothing") || strings.Contains(strings.ToLower(value), "independent") {
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
				if strings.Contains(value, "nothing") || strings.Contains(strings.ToLower(value), "independent") {
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

func (s *portalService) generateSleepRestStatistics(details []ActivityDetail) map[string]interface{} {

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

type FluidsOrdered struct {
	Water string `json:"water"`
	Milk  string `json:"milk"`
}

func (s *portalService) generateFluidsStatisticsOrdered(details []ActivityDetail) FluidsOrdered {
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

	return FluidsOrdered{
		Water: fmt.Sprintf("%dml", totalWater),
		Milk:  fmt.Sprintf("%dml", totalMilk),
	}
}

func (s *portalService) createAttendanceDetails(attendances []attendancePkg.AttendanceUserInfo) []ActivityDetail {
	var details []ActivityDetail

	for _, info := range attendances {
		// Parse attendance date
		attendanceDate, err := time.Parse("2006-01-02T15:04:05Z", info.Date)
		if err != nil {
			continue
		}

		// Create attendance data
		attendanceData := []StudentActivityData{
			{
				Key:   "attendance_id",
				Label: "Attendance ID",
				Value: info.AttendanceID,
			},
			{
				Key:   "date",
				Label: "Date",
				Value: info.Date,
			},
			{
				Key:   "temperature",
				Label: "Temperature",
				Value: fmt.Sprintf("%.2f", info.Temperature),
			},
		}

		detail := ActivityDetail{
			SessionID:    info.AttendanceID,
			TypeActivity: "attendance",
			Data:         attendanceData,
			SubmittedAt:  attendanceDate, // Use the attendance date as submitted time
			AssignedBy:   "",             // Not available from attendance data
			CreatedAt:    attendanceDate,
			UpdatedAt:    attendanceDate,
		}

		details = append(details, detail)
	}

	return details
}

func (s *portalService) generateAttendanceStatistics(details []ActivityDetail) map[string]interface{} {
	if len(details) == 0 {
		return map[string]interface{}{
			"total_records":   0,
			"temperature_avg": "0.00°C",
		}
	}

	totalTemp := 0.0
	validReadings := 0

	for _, detail := range details {
		for _, data := range detail.Data {
			if data.Key == "temperature" {
				if temp, err := strconv.ParseFloat(data.Value, 64); err == nil {
					totalTemp += temp
					validReadings++
				}
			}
		}
	}

	var avgTemp float64
	if validReadings > 0 {
		avgTemp = totalTemp / float64(validReadings)
	}

	return map[string]interface{}{
		"total_records":        len(details),
		"temperature_avg":      fmt.Sprintf("%.2f°C", avgTemp),
		"temperature_readings": validReadings,
	}
}
