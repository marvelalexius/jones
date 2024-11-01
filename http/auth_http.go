package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marvelalexius/jones/model"
	"github.com/marvelalexius/jones/utils"
	"github.com/marvelalexius/jones/utils/logger"
)

func (h *HTTPService) Register(c *gin.Context) {
	var req model.RegisterUser
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorln(c, "failed to bind json", err)
		ve := utils.ValidationResponse(err)

		utils.ErrorResponse(c, http.StatusBadRequest, utils.ErrorRes{
			Message: "something went wrong when validating the requests",
			Errors:  ve,
		})

		return
	}

	user, err := h.UserService.Register(c, &req)
	if err != nil {
		logger.Errorln(c, "failed to register", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorRes{
			Message: "something went wrong when registering",
			Errors:  err.Error(),
		})

		return
	}

	token, refreshToken, err := h.UserService.GenerateAuthTokens(user)
	if err != nil {
		logger.Errorln(c, "failed to generate auth tokens", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorRes{
			Message: "something went wrong when generating auth tokens",
			Errors:  err.Error(),
		})

		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.SuccessRes{
		Message: "user registered successfully",
		Data: map[string]interface{}{
			"user": model.AuthUser{
				User: model.User{
					ID:         user.ID,
					Name:       user.Name,
					Email:      user.Email,
					Password:   user.Password,
					Bio:        user.Bio,
					Gender:     user.Gender,
					Preference: user.Preference,
					Age:        user.Age,
					Images:     user.Images,
					CreatedAt:  user.CreatedAt,
					UpdatedAt:  user.UpdatedAt,
				},
				AuthToken:    token,
				RefreshToken: refreshToken,
			},
		},
	})
}

func (h *HTTPService) Login(c *gin.Context) {
	var req model.LoginUser
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorln(c, "failed to bind json", err)
		ve := utils.ValidationResponse(err)

		utils.ErrorResponse(c, http.StatusBadRequest, utils.ErrorRes{
			Message: "something went wrong when validating the requests",
			Errors:  ve,
		})

		return
	}

	user, err := h.UserService.Login(c, req)
	if err != nil {
		logger.Errorln(c, "failed to login", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorRes{
			Message: "something went wrong when logging in",
			Errors:  err.Error(),
		})

		return
	}

	token, refreshToken, err := h.UserService.GenerateAuthTokens(user)
	if err != nil {
		logger.Errorln(c, "failed to generate auth tokens", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorRes{
			Message: "something went wrong when generating auth tokens",
			Errors:  err.Error(),
		})

		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.SuccessRes{
		Message: "user logged in successfully",
		Data: map[string]interface{}{
			"user": model.AuthUser{
				User: model.User{
					ID:         user.ID,
					Name:       user.Name,
					Email:      user.Email,
					Password:   user.Password,
					Bio:        user.Bio,
					Gender:     user.Gender,
					Preference: user.Preference,
					Age:        user.Age,
					Images:     user.Images,
					CreatedAt:  user.CreatedAt,
					UpdatedAt:  user.UpdatedAt,
				},
				AuthToken:    token,
				RefreshToken: refreshToken,
			},
		},
	})
}

func (h *HTTPService) RefreshAuthToken(c *gin.Context) {
	var req model.RefreshToken
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorln(c, "failed to bind json", err)
		ve := utils.ValidationResponse(err)

		utils.ErrorResponse(c, http.StatusBadRequest, utils.ErrorRes{
			Message: "something went wrong when validating the requests",
			Errors:  ve,
		})

		return
	}

	token, refreshToken, err := h.UserService.RefreshAuthToken(c, req.RefreshToken)
	if err != nil {
		logger.Errorln(c, "failed to refresh auth token", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorRes{
			Message: "something went wrong when refreshing auth token",
			Errors:  err,
		})

		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.SuccessRes{
		Message: "auth token refreshed successfully",
		Data: map[string]interface{}{
			"auth_token":    token,
			"refresh_token": refreshToken,
		},
	})
}
