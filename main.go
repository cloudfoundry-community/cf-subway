package main

import (
	"log"
	"os"

	"github.com/cloudfoundry-community/cf-subway/broker"
	"github.com/codegangsta/cli"
)

func runBroker(c *cli.Context) {
	subway := broker.NewBroker()
	subway.LoadBackendBrokersFromEnv()

	err := subway.LoadCatalog()
	if err != nil {
		log.Fatalln(err)
	}

	subway.Run()
}

func main() {
	app := cli.NewApp()
	app.Name = "cf-subway"
	app.Version = "0.1.0"
	app.Usage = "Underground tunnel to multiplex many service brokers"
	app.Commands = []cli.Command{
		{
			Name:   "broker",
			Usage:  "run the broker",
			Flags:  []cli.Flag{},
			Action: runBroker,
		},
	}
	app.Run(os.Args)

}
