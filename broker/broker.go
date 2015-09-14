package broker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
	"gopkg.in/yaml.v2"
)

// Broker is the core struct for the Broker webapp
type Broker struct {
	Catalog []brokerapi.Service
}

// NewBroker is a constructor for a Broker webapp struct
func NewBroker() (broker *Broker) {
	return &Broker{}
}

// LoadCatalog loads the homogenous catalog from a file
func (broker *Broker) LoadCatalog(catalogPath string) error {
	bytes, err := ioutil.ReadFile(catalogPath)
	if err != nil {
		return err
	}

	broker.Catalog = []brokerapi.Service{}
	return yaml.Unmarshal(bytes, &broker.Catalog)
}

func (broker *Broker) plans() []*brokerapi.ServicePlan {
	plans := []*brokerapi.ServicePlan{}
	for _, service := range broker.Catalog {
		for _, plan := range service.Plans {
			plans = append(plans, &plan)
		}
	}
	return plans
}

// Run starts the Martini webapp handler
func (broker *Broker) Run() {
	logger := lager.NewLogger("cf-subway")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	credentials := brokerapi.BrokerCredentials{
		Username: "username",
		Password: "password",
	}

	brokerAPI := brokerapi.New(broker, logger, credentials)
	http.Handle("/", brokerAPI)
	port := 3000
	logger.Fatal("http-listen", http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil))
}
