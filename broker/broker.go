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
	Catalog        []brokerapi.Service
	BackendBrokers []*BackendBroker

	Logger lager.Logger
}

// BackendBroker describes the location/creds for a backend broker providing actual services
type BackendBroker struct {
	URI      string
	Username string
	Password string
}

// NewBroker is a constructor for a Broker webapp struct
func NewBroker() (subway *Broker) {
	return &Broker{}
}

// LoadCatalog loads the homogenous catalog from a file
func (subway *Broker) LoadCatalog(catalogPath string) error {
	bytes, err := ioutil.ReadFile(catalogPath)
	if err != nil {
		return err
	}

	subway.Catalog = []brokerapi.Service{}
	return yaml.Unmarshal(bytes, &subway.Catalog)
}

func (subway *Broker) plans() []*brokerapi.ServicePlan {
	plans := []*brokerapi.ServicePlan{}
	for _, service := range subway.Catalog {
		for _, plan := range service.Plans {
			plans = append(plans, &plan)
		}
	}
	return plans
}

// Run starts the Martini webapp handler
func (subway *Broker) Run() {
	subway.Logger = lager.NewLogger("cf-subway")
	subway.Logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	subway.Logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	credentials := brokerapi.BrokerCredentials{
		Username: "username",
		Password: "password",
	}

	brokerAPI := brokerapi.New(subway, subway.Logger, credentials)
	http.Handle("/", brokerAPI)
	port := 3000
	subway.Logger.Fatal("http-listen", http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil))
}
