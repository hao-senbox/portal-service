package program_planner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProgramPlanerService interface {
	CreateProgramPlaner(ctx context.Context, req *CreateProgramPlanerRequest, userID string) (string, error)
	GetAllProgramPlaner(ctx context.Context) ([]*ProgramPlaner, error)
	GetProgramPlaner(ctx context.Context, id string) (*ProgramPlaner, error)
	UpdateProgramPlaner(ctx context.Context, req *UpdateProgramPlanerRequest, id string) error
	DeleteProgramPlaner(ctx context.Context, id string) error

	CreateWeekProgramPlaner(ctx context.Context, req *CreateWeekProgramPlanerRequest, id string) error
}

type programPlannerService struct {
	ProgramPlanerRepository ProgramPlanerRepository
}

func NewProgramPlanerService(ProgramPlanerRepository ProgramPlanerRepository) ProgramPlanerService {
	return &programPlannerService{
		ProgramPlanerRepository: ProgramPlanerRepository,
	}
}

func normalizeSingleDay(code string) string {
	switch strings.ToLower(strings.TrimSpace(code)) {
	case "m", "mo", "mon", "monday":
		return "mo"
	case "t", "tu", "tue", "tuesday":
		return "tu"
	case "w", "we", "wed", "wednesday":
		return "we"
	case "th", "thu", "thursday":
		return "th"
	case "f", "fr", "fri", "friday":
		return "fr"
	case "sa", "sat", "saturday":
		return "sa"
	case "su", "sun", "sunday":
		return "su"
	default:
		return ""
	}
}

func normalizeDaysSlice(days []string) []string {
	out := make([]string, 0, len(days))
	seen := map[string]bool{}
	for _, d := range days {
		n := normalizeSingleDay(d)
		if n != "" && !seen[n] {
			out = append(out, n)
			seen[n] = true
		}
	}
	return out
}

// parseTimeRange returns timePart and comma-separated normalized days string (mo,tu,...).
// Examples:
//
//	"8:00 M-F" -> "8:00", "mo,tu,we,th,fr"
//	"11:00 sa-su" -> "11:00", "sa,su"
//	"08:00" -> "08:00", ""
func parseTimeRange(timeRange string) (string, string) {
	tr := strings.TrimSpace(timeRange)
	if tr == "" {
		return "", ""
	}
	parts := strings.Fields(tr)
	timePart := parts[0]
	var dayPart string
	if len(parts) > 1 {
		dayPart = strings.ToUpper(strings.Join(parts[1:], ""))
	}

	if dayPart == "" {
		return timePart, ""
	}

	// handle M-F and SA-SU
	switch dayPart {
	case "M-F":
		return timePart, strings.Join([]string{"mo", "tu", "we", "th", "fr"}, ",")
	case "SA-SU", "SAT-SUN":
		return timePart, strings.Join([]string{"sa", "su"}, ",")
	}

	// accept separators like "," or ";" or "|"
	replacer := strings.NewReplacer(";", ",", "|", ",")
	dayPart = replacer.Replace(dayPart)
	raw := strings.Split(dayPart, ",")
	norm := normalizeDaysSlice(raw)
	if len(norm) == 0 {
		return timePart, ""
	}
	return timePart, strings.Join(norm, ",")
}

// daysSetFromComma returns set map from "mo,tu,..." string
func daysSetFromComma(s string) map[string]bool {
	m := map[string]bool{}
	if s == "" {
		return m
	}
	parts := strings.Split(s, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			m[p] = true
		}
	}
	return m
}

// getWeeksInMonth returns a slice of (weekStart, weekEnd) pairs covering the month,
// where each week is Monday..Sunday but clipped to the month boundaries.
func getWeeksInMonth(year int, month int) []struct {
	WeekStart time.Time
	WeekEnd   time.Time
} {
	loc := time.UTC
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
	// last day of month
	nextMonthFirst := firstOfMonth.AddDate(0, 1, 0)
	lastOfMonth := nextMonthFirst.AddDate(0, 0, -1)

	// find first Monday on or before the first day (but we will clip to month start later)
	weekStart := firstOfMonth
	for weekStart.Weekday() != time.Monday {
		weekStart = weekStart.AddDate(0, 0, -1)
	}

	var weeks []struct {
		WeekStart time.Time
		WeekEnd   time.Time
	}

	for ws := weekStart; ws.Before(lastOfMonth) || ws.Equal(lastOfMonth); ws = ws.AddDate(0, 0, 7) {
		we := ws.AddDate(0, 0, 6) // Sunday
		// clip to month boundaries
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

// main CreateProgramPlaner function (replace your existing one)
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

	// compute totalFee
	var totalFee float64
	for _, slot := range req.SelectedSlots {
		totalFee += slot.Fee
	}

	// get all weeks in the month (clipped to month)
	rawWeeks := getWeeksInMonth(req.Year, req.Month)
	fmt.Printf("rawWeeks: %v\n", rawWeeks)
	// day order mo..su
	dayOrder := []string{"mo", "tu", "we", "th", "fr", "sa", "su"}

	// prepare WeekPlan slice
	var weeks []WeekPlan
	for i, w := range rawWeeks {
		// weekNumber use ISO week number from weekStart (using Monday of that week)
		isoYear, isoWeek := w.WeekStart.ISOWeek()
		_ = isoYear // we keep only week number as you requested (ISO week)
		wp := WeekPlan{
			WeekNumber: isoWeek,
			WeekStart:  w.WeekStart,
			WeekEnd:    w.WeekEnd,
			WeekFee:    0,
			Slots:      []DailySlot{},
		}
		// for each day in dayOrder we will append slots for each selectedSlot that is Selected==true
		// but only for dates that fall within this week's actual date range
		// map day code -> actual date in this week (to know if that weekday is within clipped week)
		// compute Monday of that calendar week (full week starting Monday)
		// We need baseMonday = w.WeekStart adjusted to Monday of that calendar week (not clipped)
		// To compute correct day-of-week mapping, find Monday of week containing w.WeekStart:
		// We can recompute baseMonday as the Monday <= w.WeekStart (even if clipped)
		baseMonday := w.WeekStart
		// if w.WeekStart is clipped and not a Monday, find its Monday
		for baseMonday.Weekday() != time.Monday {
			baseMonday = baseMonday.AddDate(0, 0, -1)
		}

		// Build map day->date
		dayDateMap := map[string]time.Time{}
		for idx, dowCode := range dayOrder {
			d := baseMonday.AddDate(0, 0, idx) // Monday + idx days
			dayDateMap[dowCode] = d
		}

		// Now for each day in dayOrder, if the date is within w.WeekStart..w.WeekEnd (clipped),
		// create slots for all selected selectedSlots (Selected==true)
		for _, dowCode := range dayOrder {
			dayDate := dayDateMap[dowCode]
			if dayDate.Before(w.WeekStart) || dayDate.After(w.WeekEnd) {
				// outside clipped week -> skip creating daily slots for that date
				continue
			}
			// for each selected slot create one DailySlot if applicable: Selected true -> we create the slot,
			// and DailySlot.Selected will be true only if this dowCode is in the slot's days set
			for _, sel := range req.SelectedSlots {
				if !sel.Selected {
					continue
				}
				// parse timeRange to get timePart and days from timeRange
				timePart, daysFromTR := parseTimeRange(sel.TimeRange)
				// if daysFromTR empty use sel.Days
				if daysFromTR == "" && len(sel.Days) > 0 {
					norm := normalizeDaysSlice(sel.Days)
					if len(norm) > 0 {
						daysFromTR = strings.Join(norm, ",")
					}
				}
				// default if still empty
				if daysFromTR == "" {
					daysFromTR = "mo,tu,we,th,fr"
				}
				daysSet := daysSetFromComma(daysFromTR)
				selected := false
				if daysSet[dowCode] {
					selected = true
				}
				// if timePart empty, skip creating (time required)
				if timePart == "" {
					continue
				}
				ds := DailySlot{
					DayOfWeek: dowCode,
					Time:      timePart,
					Selected:  selected,
					Fee:       0, // per requirement
				}
				wp.Slots = append(wp.Slots, ds)
			}
		}

		weeks = append(weeks, wp)
		_ = i
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

func (service *programPlannerService) CreateWeekProgramPlaner(ctx context.Context, req *CreateWeekProgramPlanerRequest, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if len(req.Weeks) == 0 {
		return fmt.Errorf("week_plan is required")
	}

	for _, week := range req.Weeks {
		if len(week.Slots) == 0 {
			return fmt.Errorf("slots is required")
		}
	}

	return service.ProgramPlanerRepository.CreateWeekProgramPlaner(ctx, req, objectID)
}

