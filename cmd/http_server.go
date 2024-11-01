package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/marvelalexius/jones/config"
	"github.com/marvelalexius/jones/http"
	"github.com/marvelalexius/jones/http/middleware"
	stripePkg "github.com/marvelalexius/jones/pkg/stripe"
	"github.com/marvelalexius/jones/repository"
	"github.com/marvelalexius/jones/service"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var runHTTPServer = &cobra.Command{
	Use:   "serve",
	Short: "Run HTTP API server",
	Long:  "Run HTTP API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		echan := make(chan error)
		go func() {
			echan <- initHTTP()
		}()

		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)

		select {
		case <-term:
			logrus.Infoln("signal termination detected")
			return nil
		case err := <-echan:
			return errors.Wrap(err, "service runtime error")
		}
	},
}

func initHTTP() error {
	// appName := os.Getenv("APP_NAME")
	appconf := config.InitConfig()

	db, err := appconf.NewDatabase()
	if err != nil {
		logrus.Fatalln("failed to connect database", err)
	}
	defer appconf.CloseDatabase(db)

	stripeClient := stripePkg.NewStripeClient(appconf.Stripe.Secret, appconf.Stripe.WebhookSecret)

	userRepo := repository.NewUserRepository(db)
	reactionRepo := repository.NewReactionRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	userService := service.NewUserService(appconf, userRepo, reactionRepo)
	reactionService := service.NewReactionService(userRepo, reactionRepo, subscriptionRepo, notificationRepo)
	subscriptionService := service.NewSubscriptionService(appconf, stripeClient, userRepo, subscriptionRepo)

	route := gin.New()
	route.Use(gin.Recovery())
	route.Use(gin.Logger())
	route.Use(gin.ErrorLogger())
	route.Use(middleware.CORS())

	httpService := http.NewHTTPService(appconf, &userService, &reactionService, &subscriptionService)
	httpService.Routes(route)

	return route.Run(":8080")
}
