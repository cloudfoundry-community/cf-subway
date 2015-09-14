package broker

import (
	"errors"
	"math/rand"
	"net/http"

	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
)

func (subway *Broker) randomBroker() *BackendBroker {
	if len(subway.BackendBrokers) == 0 {
		return nil
	}
	rand.Seed(4200)
	return subway.BackendBrokers[rand.Intn(len(subway.BackendBrokers))]
}

func (subway *Broker) routeProvision(instanceID string, planID string) (err error) {
	if len(subway.BackendBrokers) == 0 {
		return errors.New("No backend broker available for plan")
	}

	list := rand.Perm(len(subway.BackendBrokers))
	for i := range list {
		backendBroker := subway.BackendBrokers[i]
		err := subway.routeProvisionToBackendBroker(backendBroker, instanceID, planID)
		if err == nil {
			return nil
		}
	}
	return brokerapi.ErrInstanceLimitMet
}

func (subway *Broker) routeProvisionToBackendBroker(backendBroker *BackendBroker, instanceID string, planID string) (err error) {
	subway.Logger.Info("provision", lager.Data{
		"instance-id": instanceID,
		"plan-id":     planID,
		"backend-uri": backendBroker.URI,
	})

	// Dummy URI to generate test results
	if backendBroker.URI == "TEST-SUCCESS" {
		return nil
	}
	if backendBroker.URI == "TEST-NO-CAPACITY" {
		return brokerapi.ErrInstanceLimitMet
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", backendBroker.URI, nil)
	if err != nil {
		subway.Logger.Error("provision", err)
		return err
	}
	req.SetBasicAuth(backendBroker.Username, backendBroker.Password)
	_, err = client.Do(req)

	return err
}
