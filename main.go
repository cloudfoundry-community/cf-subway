package main

import (
	"log"
	"os"

	"github.com/cloudfoundry-community/cf-subway/broker"
	"github.com/codegangsta/cli"
)

func runBroker(c *cli.Context) {
	catalogPath := c.String("catalog")

	broker := broker.NewBroker()
	err := broker.LoadCatalog(catalogPath)
	if err != nil {
		log.Fatalln(err)
	}
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
					Name:  "catalog",
					Value: "service-catalog.yml",
					Usage: "service catalog supported by all backend sub-brokers",
				},
			},
			Action: runBroker,
		},
	}
	app.Run(os.Args)

}
