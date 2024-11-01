package cmd

import (
	"encoding/json"
	"io"
	"log"
	"os"

	_ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"
	"github.com/marvelalexius/jones/config"
	"github.com/marvelalexius/jones/model"
	"github.com/marvelalexius/jones/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var seederCmd = &cobra.Command{
	Use:   "seed",
	Short: "seeder cmd",
	Long:  `This subcommand to execute seeding process`,
	Run:   seeder,
}

func init() {
	rootCmd.AddCommand(seederCmd)
}

func seeder(cmd *cobra.Command, args []string) {
	appconf := config.InitConfig()

	db, err := appconf.NewDatabase()
	if err != nil {
		logrus.Fatalln("failed to connect database", err)
	}
	defer appconf.CloseDatabase(db)

	subscriptionRepo := repository.NewSubscriptionRepository(db)

	err = seedSubscriptionPlan(subscriptionRepo)
	continueOrFatal(err)

	log.Print("seed success")
}

func continueOrFatal(err error) {
	if err != nil {
		logrus.Fatal(err.Error())
	}
}

func seedSubscriptionPlan(subscriptionRepo repository.ISubscriptionRepository) error {
	var subsPlan []model.SubscriptionPlan

	subs, err := os.Open("seeder/subscription_plan.json")
	if err != nil {
		return err
	}

	byteSubs, _ := io.ReadAll(subs)

	err = json.Unmarshal(byteSubs, &subsPlan)
	if err != nil {
		return err
	}

	err = subscriptionRepo.BulkCreatePlan(subsPlan)
	if err != nil {
		return err
	}

	return nil
}
