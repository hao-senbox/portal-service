package program_planner

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProgramPlanerService interface {
	CreateProgramPlaner(ctx context.Context, req *CreateProgramPlanerRequest, userID string) (string, error)
	GetAllProgramPlaner(ctx context.Context) ([]*ProgramPlaner, error)
	GetProgramPlaner(ctx context.Context, id string) (*ProgramPlaner, error)
	UpdateProgramPlaner(ctx context.Context, req *UpdateProgramPlanerRequest, id string) error
	DeleteProgramPlaner(ctx context.Context, id string) error

	UpdateProgramPlanerWeek(ctx context.Context, req *UpdateWeekProgramPlanerRequest, id string) error
}

type programPlannerService struct {
	ProgramPlanerRepository ProgramPlanerRepository
}

func NewProgramPlanerService(ProgramPlanerRepository ProgramPlanerRepository) ProgramPlanerService {
	return &programPlannerService{
		ProgramPlanerRepository: ProgramPlanerRepository,
	}
}

func (service *programPlannerService) CreateProgramPlaner(ctx context.Context, req *CreateProgramPlanerRequest, userID string) (string, error) {

	if req.StudentID == "" {
		return "", fmt.Errorf("student_id is required")
	}
	if req.OrganizationID == "" {
		return "", fmt.Errorf("organization_id is required")
	}
	if req.Month == 0 || req.Year == 0 {
		return "", fmt.Errorf("month and year is required")
	}
	if len(req.SelectedSlots) == 0 {
		return "", fmt.Errorf("selected_slots is required")
	}

	var totalFee float64
	for _, slot := range req.SelectedSlots {
		totalFee += slot.Fee
	}

	rawWeeks := getWeeksInMonth(req.Year, req.Month)

	weeks := make([]WeekPlan, 0)
	for _, w := range rawWeeks {
		_, isoWeek := w.WeekStart.ISOWeek()
		wp := WeekPlan{
			WeekNumber: isoWeek,
			WeekStart:  w.WeekStart,
			WeekEnd:    w.WeekEnd,
			WeekFee:    0,
			Slots:      []DailySlot{},
		}

		dayCodes := []string{"mo", "tu", "we", "th", "fr", "sa", "su"}
		baseMonday := w.WeekStart
		for baseMonday.Weekday() != time.Monday {
			baseMonday = baseMonday.AddDate(0, 0, -1)
		}
		dayDateMap := map[string]time.Time{}
		for i, code := range dayCodes {
			dayDate := baseMonday.AddDate(0, 0, i)
			dayDateMap[code] = dayDate
		}

		for _, sel := range req.SelectedSlots {

			if !sel.Selected {
				continue
			}

			var timePart string
			switch sel.TimeRange {
			case "8:00 M-F":
				timePart = "8:00"
			case "11:00 M-F":
				timePart = "11:00"
			case "17:00 M-F":
				timePart = "17:00"
			case "20:00 M-F":
				timePart = "20:00"
			default:
				continue
			}

			for _, code := range sel.Days {
				dayDate := dayDateMap[code]
				if dayDate.Before(w.WeekStart) || dayDate.After(w.WeekEnd) {
					continue
				}
				wp.Slots = append(wp.Slots, DailySlot{
					DayOfWeek:  code,
					Time:       timePart,
					Selected:   true,
					Fee:        0,
					IsOriginal: true,
				})
			}

			// for _, extra := range sel.Days {
			// 	if extra != "sa" && extra != "su" {
			// 		continue
			// 	}
			// 	dayDate := dayDateMap[extra]
			// 	if dayDate.Before(w.WeekStart) || dayDate.After(w.WeekEnd) {
			// 		continue
			// 	}
			// 	wp.Slots = append(wp.Slots, DailySlot{
			// 		DayOfWeek:  extra,
			// 		Time:       timePart,
			// 		Selected:   true,
			// 		Fee:        0,
			// 		IsOriginal: true,
			// 	})
			// }
		}

		weeks = append(weeks, wp)
	}

	programPlaner := &ProgramPlaner{
		ID:             primitive.NewObjectID(),
		StudentID:      req.StudentID,
		OrganizationID: req.OrganizationID,
		Month:          req.Month,
		Year:           req.Year,
		TotalFee:       totalFee,
		SelectedSlots:  req.SelectedSlots,
		Weeks:          weeks,
		CreatedBy:      userID,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return service.ProgramPlanerRepository.CreateProgramPlaner(ctx, programPlaner)
}

func (service *programPlannerService) GetAllProgramPlaner(ctx context.Context) ([]*ProgramPlaner, error) {
	return service.ProgramPlanerRepository.GetAllProgramPlaner(ctx)
}

func (service *programPlannerService) GetProgramPlaner(ctx context.Context, id string) (*ProgramPlaner, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return service.ProgramPlanerRepository.GetProgramPlaner(ctx, objectID)
}

func (service *programPlannerService) UpdateProgramPlaner(ctx context.Context, req *UpdateProgramPlanerRequest, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	program, err := service.ProgramPlanerRepository.GetProgramPlaner(ctx, objectID)
	if err != nil {
		return err
	}

	if req.Month != 0 {
		program.Month = req.Month
	}

	if req.Year != 0 {
		program.Year = req.Year
	}

	if len(req.SelectedSlots) != 0 {
		program.SelectedSlots = req.SelectedSlots
	}

	var totalFeeUpdate float64
	for _, slot := range program.SelectedSlots {
		totalFeeUpdate += slot.Fee
	}

	program.TotalFee = totalFeeUpdate

	return service.ProgramPlanerRepository.UpdateProgramPlaner(ctx, program, objectID)

}

func (service *programPlannerService) DeleteProgramPlaner(ctx context.Context, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return service.ProgramPlanerRepository.DeleteProgramPlaner(ctx, objectID)
}

func (service *programPlannerService) UpdateProgramPlanerWeek(ctx context.Context, req *UpdateWeekProgramPlanerRequest, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	program, err := service.ProgramPlanerRepository.GetProgramPlaner(ctx, objectID)
	if err != nil {
		return err
	}

	if program == nil {
		return errors.New("program planer not found")
	}

	weekIndex := -1

	for i, week := range program.Weeks {
		if week.WeekNumber == req.WeekNumber {
			weekIndex = i
			break
		}
	}

	if weekIndex == -1 {
		return fmt.Errorf("week %d not found", req.WeekNumber)
	}

	slotIndex := -1

	for i, slot := range program.Weeks[weekIndex].Slots {
		if slot.DayOfWeek == req.DayOfWeek && slot.Time == req.Time {
			slotIndex = i
			break
		}
	}

	slotFee := req.SlotFee

	if slotIndex != -1 {
		slot := &program.Weeks[weekIndex].Slots[slotIndex]
		if slot.IsOriginal {
			if slot.Selected {
				slot.Selected = false
				slot.Fee = -slotFee
			} else {
				slot.Selected = true
				slot.Fee = 0
			}
		} else {
			if slot.Selected {
				slot.Selected = false
				slot.Fee = 0
			} else {
				slot.Selected = true
				slot.Fee = slotFee
			}
		}
	} else {
		newSlot := DailySlot{
			DayOfWeek:  req.DayOfWeek,
			Time:       req.Time,
			Selected:   true,
			Fee:        slotFee,
			IsOriginal: false,
		}
		program.SelectedSlots = append(program.SelectedSlots, SelectedSlot{
			TimeRange: req.Time,
			Days:      []string{req.DayOfWeek},
			Selected:  true,
			Fee:       slotFee,
		})
		program.Weeks[weekIndex].Slots = append(program.Weeks[weekIndex].Slots, newSlot)
	}

	program.Weeks[weekIndex].WeekFee = service.calculateWeekFee(program.Weeks[weekIndex].Slots)

	program.TotalFee = service.calculateTotalFee(program.Weeks, program.SelectedSlots)

	return service.ProgramPlanerRepository.UpdateProgramPlanerWeek(ctx, program, objectID)
}

func (service *programPlannerService) calculateWeekFee(slots []DailySlot) float64 {

	total := 0.0

	for _, s := range slots {
		total += s.Fee
	}

	return total
}

func (service *programPlannerService) calculateTotalFee(weeks []WeekPlan, selectedSlots []SelectedSlot) float64 {

	baseFee := 0.0
	for _, slot := range selectedSlots {
		if slot.Selected {
			baseFee += slot.Fee
		}
	}

	for _, week := range weeks {
		baseFee += week.WeekFee
	}

	return baseFee
}

func getWeeksInMonth(year int, month int) []struct {
	WeekStart time.Time
	WeekEnd   time.Time
} {
	loc := time.UTC
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
	nextMonthFirst := firstOfMonth.AddDate(0, 1, 0)
	lastOfMonth := nextMonthFirst.AddDate(0, 0, -1)

	weekStart := firstOfMonth
	for weekStart.Weekday() != time.Monday {
		weekStart = weekStart.AddDate(0, 0, -1)
	}

	var weeks []struct {
		WeekStart time.Time
		WeekEnd   time.Time
	}

	for ws := weekStart; ws.Before(lastOfMonth) || ws.Equal(lastOfMonth); ws = ws.AddDate(0, 0, 7) {
		we := ws.AddDate(0, 0, 6)
		actualStart := ws
		if actualStart.Before(firstOfMonth) {
			actualStart = firstOfMonth
		}
		actualEnd := we
		if actualEnd.After(lastOfMonth) {
			actualEnd = lastOfMonth
		}
		weeks = append(weeks, struct {
			WeekStart time.Time
			WeekEnd   time.Time
		}{WeekStart: actualStart, WeekEnd: actualEnd})
	}
	return weeks
}
