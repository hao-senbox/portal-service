package drink

import (
	"context"
	"fmt"
	"portal/helper"
	"portal/pkg/constants"

	"github.com/gin-gonic/gin"
)

type DrinkHandler struct {
	DrinkService DrinkService
}

func NewDrinkHandler(DrinkService DrinkService) *DrinkHandler {
	return &DrinkHandler{
		DrinkService: DrinkService,
	}
}

func (h *DrinkHandler) CreateDrink(c *gin.Context) {

	var req CreateDrinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	drinkID, err := h.DrinkService.CreateDrink(ctx, &req, userID.(string))
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Create drink successfully", drinkID)

}

func (h *DrinkHandler) GetDrinks(c *gin.Context) {

	studentID := c.Query("student")
	date := c.Query("date")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	drinks, err := h.DrinkService.GetDrinks(ctx, studentID, date)

	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Get drinks successfully", drinks)

}

func (h *DrinkHandler) GetDrink(c *gin.Context) {

	id := c.Param("id")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	drink, err := h.DrinkService.GetDrink(ctx, id)

	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Get drink successfully", drink)

}

func (h *DrinkHandler) GetStatistics(c *gin.Context) {

	student_id := c.Query("student")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	statistics, err := h.DrinkService.GetStatistics(ctx, student_id)

	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Get statistics successfully", statistics)

}