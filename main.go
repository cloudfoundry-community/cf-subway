package main

import (
	"os"

	"github.com/dajulia3/cli"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func runBroker(c *cli.Context) {
	// pipelinesConfigPath := c.String("config")
	// var err error
	// webserverConfig, err = config.LoadConfigFromYAMLFile(pipelinesConfigPath)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	m := martini.Classic()
	m.Use(render.Renderer())
	// m.Use(auth.Basic(webserverConfig.Auth.Username, webserverConfig.Auth.Password))
	// m.Get("/", brokerShowHelp)
	m.Run()
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
					Name:  "config",
					Value: "config.yml",
					Usage: "configuration",
				},
			},
			Action: runBroker,
		},
	}
	app.Run(os.Args)

}
