package broker

import (
	"net/http"

	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
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
	logger := lager.NewLogger("my-service-broker")
	credentials := brokerapi.BrokerCredentials{
		Username: "username",
		Password: "password",
	}

	brokerAPI := brokerapi.New(broker, logger, credentials)
	http.Handle("/", brokerAPI)
	http.ListenAndServe(":3000", nil)
}
