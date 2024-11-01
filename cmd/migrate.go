package cmd

import (
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/marvelalexius/jones/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migration cmd",
	Long:  `This subcommand to execute migration process`,
	Run:   migrate,
}

func init() {
	migrateCmd.PersistentFlags().String("direction", "", "migration direction")
	migrateCmd.PersistentFlags().String("new", "", "new migration file")
	rootCmd.AddCommand(migrateCmd)
}

func migrate(cmd *cobra.Command, args []string) {
	direction := cmd.Flag("direction").Value.String()
	newFile := cmd.Flag("new").Value.String()

	dsn, err := url.Parse(config.PostgresDSN())
	if err != nil {
		log.WithField("ParseDSN", config.PostgresDSN()).Fatal("Failed to parse dsn: ", err)
	}

	db := dbmate.New(dsn)
	db.MigrationsDir = []string{"./migration"}
	db.SchemaFile = "./migration/_schema.sql"

	if newFile != "" {
		err := db.NewMigration(newFile)
		if err != nil {
			log.WithField("NewMigration", config.PostgresDSN()).Fatal("Failed to connect db: ", err)
		}

		log.Infof("Success creating %s file!", newFile)
		return
	}

	if direction != "" {
		switch direction {
		case "down":
			err := db.Rollback()
			if err != nil {
				log.Fatal("Rollback failed: ", err)
			}
		default:
			err := db.Migrate()
			if err != nil {
				log.Fatal("Migrate failed: ", err)
			}
		}

		total := getTotalApplied(db)

		log.Infof("Applied %d migrations!\n", total)
	}
}

func getTotalApplied(db *dbmate.DB) int {
	results, err := db.FindMigrations()
	if err != nil {
		return 0
	}

	var totalApplied int

	for _, res := range results {
		if res.Applied {
			totalApplied++
			continue
		}
	}

	return totalApplied
}
