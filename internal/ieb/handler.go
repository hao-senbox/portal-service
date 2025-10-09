package ieb

import (
	"context"
	"fmt"
	"portal/helper"
	"portal/pkg/constants"

	"github.com/gin-gonic/gin"
)

type IEBHandler struct {
	IEBService IEBService
}

func NewIEBHandler(IEBService IEBService) *IEBHandler {
	return &IEBHandler{
		IEBService: IEBService,
	}
}

func (handler *IEBHandler) CreateIEB(c *gin.Context) {

	var req CreateIEBRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	iebID, err := handler.IEBService.CreateIEB(ctx, &req, userID.(string))
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Create ieb successfully", iebID)

}

func (handler *IEBHandler) GetIEB(c *gin.Context) {

	userID := c.Query("user_id")
	termID := c.Query("term_id")
	languageKey := c.Query("language_key")
	regionKey := c.Query("region_key")
	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	ieb, err := handler.IEBService.GetIEB(ctx, userID, termID, languageKey, regionKey)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Get ieb successfully", ieb)

}
