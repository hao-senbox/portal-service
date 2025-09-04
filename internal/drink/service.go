package drink

import (
	"context"
	"fmt"
	"portal/internal/user"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DrinkService interface {
	CreateDrink(ctx context.Context, req *CreateDrinkRequest, userID string) (string, error)
	GetDrinks(ctx context.Context, studentID string, date string) ([]*DrinkResponse, error)
	GetDrink(ctx context.Context, id string) (*DrinkResponse, error)
	GetStatistics(ctx context.Context, studentID string, date string) (*DrinkDailyTotals, error)
}

type drinkService struct {
	DrinkRepository DrinkRepository
	UserService     user.UserService
}

func NewDrinkService(DrinkRepository DrinkRepository, UserService user.UserService) DrinkService {
	return &drinkService{
		DrinkRepository: DrinkRepository,
		UserService:     UserService,
	}
}

func (s *drinkService) CreateDrink(ctx context.Context, req *CreateDrinkRequest, userID string) (string, error) {

	if req.Date == "" {
		return "", fmt.Errorf("date is required")
	}

	parseDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return "", fmt.Errorf("invalid date format")
	}

	if req.StudentID == "" {
		return "", fmt.Errorf("student_id is required")
	}

	if len(req.Liquids) == 0 {
		return "", fmt.Errorf("liquids is required")
	}

	var liquids []Liquid

	for _, liquid := range req.Liquids {

		if liquid.Type == "" {
			return "", fmt.Errorf("type is required")
		}

		if liquid.Amount == 0 {
			return "", fmt.Errorf("amount is required")
		}

		liquid := Liquid{
			Type:   liquid.Type,
			Amount: liquid.Amount,
		}

		liquids = append(liquids, liquid)

	}

	drink := Drink{
		Date:      parseDate,
		StudentID: req.StudentID,
		Liquids:   liquids,
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.DrinkRepository.CreateDrink(ctx, &drink)

}

func (s *drinkService) GetDrinks(ctx context.Context, studentID string, date string) ([]*DrinkResponse, error) {

	var dateRepo *time.Time

	if studentID == "" {
		studentID = ""
	}

	if date != "" {
		parseDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format")
		}
		dateRepo = &parseDate
	} else {
		dateRepo = nil
	}

	var result []*DrinkResponse

	res, err := s.DrinkRepository.GetDrinks(ctx, studentID, dateRepo)

	if err != nil {
		return nil, err
	}

	for _, drink := range res {

		student, err := s.UserService.GetStudentInfor(ctx, drink.StudentID)
		if err != nil {
			return nil, err
		}

		teacher, err := s.UserService.GetUserInfor(ctx, drink.CreatedBy)
		if err != nil {
			return nil, err
		}

		drinkResponse := DrinkResponse{
			ID:        drink.ID,
			Student:   student,
			Date:      drink.Date.Format("2006-01-02"),
			Liquids:   drink.Liquids,
			Teacher:   teacher,
			CreatedAt: drink.CreatedAt,
			UpdatedAt: drink.UpdatedAt,
		}

		result = append(result, &drinkResponse)
	}

	return result, nil

}

func (s *drinkService) GetDrink(ctx context.Context, id string) (*DrinkResponse, error) {

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	drink, err := s.DrinkRepository.GetDrink(ctx, objectID)

	if err != nil {
		return nil, err
	}

	student, err := s.UserService.GetStudentInfor(ctx, drink.StudentID)
	if err != nil {
		return nil, err
	}

	teacher, err := s.UserService.GetUserInfor(ctx, drink.CreatedBy)
	if err != nil {
		return nil, err
	}

	drinkResponse := DrinkResponse{
		ID:        drink.ID,
		Student:   student,
		Date:      drink.Date.Format("2006-01-02"),
		Liquids:   drink.Liquids,
		Teacher:   teacher,
		CreatedAt: drink.CreatedAt,
		UpdatedAt: drink.UpdatedAt,
	}

	return &drinkResponse, nil

}

func (s *drinkService) GetStatistics(ctx context.Context, studentID string, date string) (*DrinkDailyTotals, error) {

	var dateRepo *time.Time

	if date != "" {
		parseDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format")
		}
		dateRepo = &parseDate
	} else {
		dateRepo = nil
	}

	res, err := s.DrinkRepository.GetDrinks(ctx, studentID, dateRepo)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	if res == nil {
		return nil, nil
	}

	totals := make(map[string]float64)
	for _, drink := range res {
		for _, liquid := range drink.Liquids {
			totals[liquid.Type] += liquid.Amount
		}
	}

	var statistics []Satistic

	for liquidType, total := range totals {
		statistics = append(statistics, Satistic{
			Type:  liquidType,
			Total: total,
		})
	}

	sort.Slice(statistics, func(i, j int) bool {
		return statistics[i].Type < statistics[j].Type
	})

	teacher, err := s.UserService.GetUserInfor(ctx, res[0].CreatedBy)
	if err != nil {
		return nil, err
	}

	student, err := s.UserService.GetStudentInfor(ctx, res[0].StudentID)
	if err != nil {
		return nil, err
	}

	return &DrinkDailyTotals{
		Teacher:    teacher,
		Student:    student,
		Date:       date,
		Statistics: statistics,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}
