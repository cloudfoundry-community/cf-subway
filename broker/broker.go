package broker

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

// Broker is the core struct for the Broker webapp
type Broker struct {
}

// NewBroker is a constructor for a Broker webapp struct
func NewBroker() *Broker {
	return &Broker{}
}

// Run starts the Martini webapp handler
func (broker *Broker) Run() {
	m := martini.Classic()
	m.Use(render.Renderer())
	// m.Use(auth.Basic(webserverConfig.Auth.Username, webserverConfig.Auth.Password))
	// m.Get("/", brokerShowHelp)
	m.Run()
}
