package cmd

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Jones API",
	Short: "Run JonesHTTP API server",
	Long:  "Run Jones HTTP API server",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(runHTTPServer)

	//load environment variable
	if err := godotenv.Load(); err != nil {
		if os.Getenv("APP_ENV") == "development" {
			logrus.Fatalln("unable to load environment variable", err.Error())
		} else {
			logrus.Warningln("Can't find env.file. To use system's env vars for now")
		}
	}
}

func Execute() error {
	return rootCmd.Execute()
}
