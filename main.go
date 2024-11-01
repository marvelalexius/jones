package main

import (
	"os"

	"github.com/marvelalexius/jones/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.Execute(); err != nil {
		logrus.Fatalln("error on running the application", err.Error())
		os.Exit(1)
	}
}
