package main

import (
	"os"

	"github.com/cloudfoundry-community/cf-subway/broker"
	"github.com/codegangsta/cli"
)

func runBroker(c *cli.Context) {
	broker := broker.NewBroker()
	broker.Run()
}

func main() {
	app := cli.NewApp()
	app.Name = "cf-subway"
	app.Version = "0.1.0"
	app.Usage = "Underground tunnel to multiplex many service brokers"
	app.Commands = []cli.Command{
		{
			Name:  "broker",
			Usage: "run the broker",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Value: "broker_config.yml",
					Usage: "configuration",
				},
			},
			Action: runBroker,
		},
	}
	app.Run(os.Args)

}
