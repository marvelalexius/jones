package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marvelalexius/jones/model"
	"github.com/marvelalexius/jones/utils"
	"github.com/marvelalexius/jones/utils/logger"
)

func (h *HTTPService) FindAllUsers(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		logger.Errorln(c, "failed to get user id from context")
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorRes{
			Message: "something went wrong when finding all users",
		})

		return
	}

	users, count, err := h.UserService.FindAll(c, userID.(string))
	if err != nil {
		logger.Errorln(c, "failed to find all users", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorRes{
			Message: "something went wrong when finding all users",
			Errors:  err,
		})

		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.SuccessRes{
		Message: "success",
		Data:    users,
		Meta: map[string]interface{}{
			"total": count,
		},
	})
}

func (h *HTTPService) React(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		logger.Errorln(c, "failed to get user id from context")
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorRes{
			Message: "something went wrong when swiping",
		})

		return
	}

	var req model.ReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorln(c, "failed to bind json", err)
		ve := utils.ValidationResponse(err)
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ErrorRes{
			Message: "something went wrong when validating the requests",
			Errors:  ve,
		})

		return
	}

	req.UserID = userID.(string)
	reaction, err := h.ReactionService.Swipe(c, req)
	if err != nil {
		logger.Errorln(c, "failed to swipe", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorRes{
			Message: "something went wrong when swiping",
			Errors:  err.Error(),
		})

		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.SuccessRes{
		Message: "success",
		Data:    reaction,
	})
}
