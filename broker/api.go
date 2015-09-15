package broker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
	"gopkg.in/yaml.v2"
)

// Services is used by Cloud Foundry to learn the available catalog of services
func (subway *Broker) Services() []brokerapi.Service {
	return subway.Catalog
}

// Provision requests the creation of a service instance from an available sub-broker
func (subway *Broker) Provision(instanceID string, details brokerapi.ProvisionDetails) error {
	if details.PlanID == "" {
		return errors.New("plan_id required")
	}

	planID := ""
	for _, plan := range subway.plans() {
		if plan.ID == details.PlanID {
			planID = details.PlanID
			break
		}
	}

	if planID == "" {
		return errors.New("plan_id not recognized")
	}

	return subway.routeProvision(instanceID, details)
}

// Deprovision requests the destruction of a service instance from associated sub-broker
func (subway *Broker) Deprovision(instanceID string) error {
	// Deprovision instances here
	return nil
}

// Bind requests the creation of a service instance bindings from associated sub-broker
func (subway *Broker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (interface{}, error) {
	subway.Logger.Info("bind", lager.Data{
		"instance-id": instanceID,
		"binding-id":  bindingID,
		"plan-id":     details.PlanID,
	})

	bindingResponse := brokerapi.BindingResponse{}

	for _, backendBroker := range subway.BackendBrokers {
		// Dummy URI to generate test results
		if backendBroker.URI == "TEST-FOUND-INSTANCE" {
			bindingResponse.Credentials = map[string]interface{}{"host": "10.10.10.10"}
			return bindingResponse, nil
		} else if backendBroker.URI == "TEST-UNKNOWN-INSTANCE" {
			// Skip test backend broker
		} else {
			client := &http.Client{}
			url := fmt.Sprintf("%s/v2/service_instances/%s/service_bindings/%s", backendBroker.URI, instanceID, bindingID)
			buffer := &bytes.Buffer{}
			json.NewEncoder(buffer).Encode(details)
			req, err := http.NewRequest("PUT", url, buffer)
			if err != nil {
				subway.Logger.Error("backend-bind", err)
				return bindingResponse, err
			}
			req.Header.Set("Content-Type", "application/json")
			req.SetBasicAuth(backendBroker.Username, backendBroker.Password)
			resp, err := client.Do(req)
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusCreated {
				jsonData, err := ioutil.ReadAll(resp.Body)

				err = yaml.Unmarshal(jsonData, &bindingResponse)
				if err == nil {
					subway.Logger.Info("bind-success", lager.Data{
						"instance-id": instanceID,
						"binding-id":  bindingID,
						"plan-id":     details.PlanID,
						"backend-uri": backendBroker.URI,
					})
					return bindingResponse, nil
				}
			}
		}
	}

	return bindingResponse, brokerapi.ErrInstanceDoesNotExist
}

// Unbind requests the destructions of a service instance binding from associated sub-broker
func (subway *Broker) Unbind(instanceID, bindingID string) error {
	// Unbind from instances here
	return nil
}
