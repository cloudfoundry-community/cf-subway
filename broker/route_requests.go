package broker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
)

func (subway *Broker) routeProvision(instanceID string, details brokerapi.ProvisionDetails) (err error) {
	if len(subway.BackendBrokers) == 0 {
		return errors.New("No backend broker available for plan")
	}

	list := rand.Perm(len(subway.BackendBrokers))
	for i := range list {
		backendBroker := subway.BackendBrokers[i]
		err := subway.routeProvisionToBackendBroker(backendBroker, instanceID, details)
		if err == nil {
			return nil
		}
	}
	return brokerapi.ErrInstanceLimitMet
}

func (subway *Broker) routeProvisionToBackendBroker(backendBroker *BackendBroker, instanceID string, details brokerapi.ProvisionDetails) (err error) {
	subway.Logger.Info("provision", lager.Data{
		"instance-id": instanceID,
		"plan-id":     details.PlanID,
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
	url := fmt.Sprintf("%s/v2/service_instances/%s", backendBroker.URI, instanceID)
	buffer := &bytes.Buffer{}
	json.NewEncoder(buffer).Encode(details)
	req, err := http.NewRequest("PUT", url, buffer)
	if err != nil {
		subway.Logger.Error("backend-provision", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(backendBroker.Username, backendBroker.Password)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	// TODO: ProvisioningResponse
	return err
}
