package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/marvelalexius/jones/config"
	"github.com/marvelalexius/jones/utils"
	"github.com/marvelalexius/jones/utils/logger"
	"github.com/marvelalexius/jones/utils/str"
)

func extractToken(c *gin.Context) (string, error) {
	bearerToken := c.Request.Header.Get("Authorization")
	err := errors.New("no Authorization token detected")

	// Apple already reserved header for Authorization
	// https://developer.apple.com/documentation/foundation/nsurlrequest
	if bearerToken == "" {
		bearerToken = c.Request.Header.Get("X-Authorization")
	}

	if len(strings.Split(bearerToken, " ")) == 2 {
		bearerToken = strings.Split(bearerToken, " ")[1]
	}

	if bearerToken == "" {
		return "", err
	}

	return bearerToken, nil
}

func JWTAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		extractedToken, err := extractToken(ctx)
		if err != nil {
			logger.Errorln(ctx, "failed to extract token", err)
			utils.ErrorResponse(ctx, http.StatusUnauthorized, utils.ErrorRes{
				Message: "Invalid token",
				Errors:  err.Error(),
			})
			ctx.Abort()
			return
		}

		parsedToken, err := str.ParseJWT(extractedToken, cfg.App.Secret)
		if err != nil {
			logger.Errorln(ctx, "failed to parse token", err)
			utils.ErrorResponse(ctx, http.StatusUnauthorized, utils.ErrorRes{
				Message: "Invalid token",
				Errors:  err.Error(),
			})
			ctx.Abort()
			return
		}

		ctx.Set("userID", parsedToken.UserID)
		ctx.Next()
	}
}
