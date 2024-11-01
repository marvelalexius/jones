package http

import (
	"github.com/gin-gonic/gin"
	"github.com/marvelalexius/jones/config"
	"github.com/marvelalexius/jones/http/middleware"
	"github.com/marvelalexius/jones/service"
)

type HTTPService struct {
	Conf                *config.Config
	UserService         service.IUserService
	ReactionService     service.IReactionService
	SubscriptionService service.ISubscriptionService
}

func NewHTTPService(appconf *config.Config, userService *service.IUserService, reactionService *service.IReactionService, subscriptionService *service.ISubscriptionService) *HTTPService {
	return &HTTPService{Conf: appconf, UserService: *userService, ReactionService: *reactionService, SubscriptionService: *subscriptionService}
}

func (h *HTTPService) Routes(route *gin.Engine) {
	api := route.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.POST("/auth/register", h.Register)
			v1.POST("/auth/login", h.Login)
			v1.POST("/auth/refresh", h.RefreshAuthToken)

			authed := v1.Group("").Use(middleware.JWTAuthMiddleware(h.Conf))
			authed.GET("/users", h.FindAllUsers)
			authed.POST("/reactions", h.React)
			authed.GET("/reactions/likes", h.SeeLikes)
			authed.POST("/subscription", h.Subscribe)

			if h.Conf.FeatureFlag.EnableStripe {
				v1.POST("/payment/callback", h.HandleCallback)
				authed.GET("/subscription/portal", h.ManageSubscription)
			}
		}
	}
}
