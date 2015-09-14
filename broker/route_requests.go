package broker

import (
	"errors"
	"math/rand"
	"net/http"
)

func (subway *Broker) randomBroker() *BackendBroker {
	if len(subway.BackendBrokers) == 0 {
		return nil
	}
	rand.Seed(4200)
	return subway.BackendBrokers[rand.Intn(len(subway.BackendBrokers))]
}

func (subway *Broker) routeProvision(instanceID string, planID string) (err error) {
	backendBroker := subway.randomBroker()
	if backendBroker == nil {
		return errors.New("No backend broker available for plan")
	}

	if backendBroker.URI == "TESTDUMMY" {
		return nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", backendBroker.URI, nil)
	if err != nil {
		subway.Logger.Error("provision", err)
		return err
	}
	req.SetBasicAuth(backendBroker.Username, backendBroker.Password)
	_, err = client.Do(req)

	// forward response from backend broker to CF
	// TODO: if error, then keep looking for another broker until tried all backend brokers
	return err
}
