package broker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/frodenas/brokerapi"
	"github.com/mitchellh/mapstructure"
	"github.com/pivotal-golang/lager"
)

// Services is used by Cloud Foundry to learn the available catalog of services
func (subway *Broker) Services() brokerapi.CatalogResponse {
	err := subway.LoadCatalog()
	if err != nil {
		subway.Logger.Error("catalog", err)
	}

	return subway.BackendCatalog
}

// Provision requests the creation of a service instance from an available sub-broker
func (subway *Broker) Provision(instanceID string, details brokerapi.ProvisionDetails, acceptsIncomplete bool) (resp brokerapi.ProvisioningResponse, doesAcceptIncomplete bool, err error) {
	if details.PlanID == "" {
		return resp, false, errors.New("plan_id required")
	}

	planID := ""
	for _, plan := range subway.plans() {
		if plan.ID == details.PlanID {
			planID = details.PlanID
			break
		}
	}

	if planID == "" {
		return resp, false, errors.New("plan_id not recognized")
	}

	err = subway.routeProvision(instanceID, details)
	return resp, false, err
}

// Update service instance
func (subway *Broker) Update(instanceID string, details brokerapi.UpdateDetails, acceptsIncomplete bool) (doesAcceptIncomplete bool, err error) {
	return false, fmt.Errorf("Update not supported yet")
}

// Bind requests the creation of a service instance bindings from associated sub-broker
func (subway *Broker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (bindResp brokerapi.BindingResponse, err error) {
	subway.Logger.Info("bind", lager.Data{
		"instance-id": instanceID,
		"binding-id":  bindingID,
		"service-id":  details.ServiceID,
		"plan-id":     details.PlanID,
	})

	for _, backendBroker := range subway.BackendBrokers {
		// Dummy URI to generate test results
		if backendBroker.URI == "TEST-FOUND-INSTANCE" {
			bindResp.Credentials = map[string]interface{}{"host": "10.10.10.10"}
			return bindResp, nil
		} else if backendBroker.URI == "TEST-UNKNOWN-INSTANCE" {
			// Skip test backend broker
		} else {
			client := &http.Client{}
			url := fmt.Sprintf("%s/v2/service_instances/%s/service_bindings/%s", backendBroker.URI, instanceID, bindingID)
			buffer := &bytes.Buffer{}
			if err := json.NewEncoder(buffer).Encode(details); err != nil {
				subway.Logger.Error("backend-bind-encode-details", err)
				return bindResp, err
			}
			req, err := http.NewRequest("PUT", url, buffer)
			if err != nil {
				subway.Logger.Error("backend-bind-req", err)
				return bindResp, err
			}
			req.Header.Set("Content-Type", "application/json")
			req.SetBasicAuth(backendBroker.Username, backendBroker.Password)
			debug(httputil.DumpRequestOut(req, true))

			resp, err := client.Do(req)
			if err != nil {
				subway.Logger.Error("backend-bind-resp", err)
				return bindResp, err
			}
			defer resp.Body.Close()

			debug(httputil.DumpResponse(resp, true))

			if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
				jsonData, err := ioutil.ReadAll(resp.Body)

				rawBindingResponse := map[string]interface{}{}
				if err = json.Unmarshal(jsonData, &rawBindingResponse); err != nil {
					return bindResp, err
				}
				fmt.Printf("%#v\n", rawBindingResponse)
				if err = mapstructure.WeakDecode(rawBindingResponse, &bindResp); err != nil {
					return bindResp, err
				}
				// HACK for some reason WeakDecode doesn't parse "syslog_drain_url" into .SyslogDrainURL
				if rawBindingResponse["syslog_drain_url"] != nil {
					bindResp.SyslogDrainURL = rawBindingResponse["syslog_drain_url"].(string)
				}
				if err == nil {
					subway.Logger.Info("bind-success", lager.Data{
						"instance-id": instanceID,
						"binding-id":  bindingID,
						"plan-id":     details.PlanID,
						"backend-uri": backendBroker.URI,
					})
					fmt.Printf("%#v\n", bindResp)
					return bindResp, nil
				}
			}
		}
	}

	return bindResp, brokerapi.ErrInstanceDoesNotExist
}

// Unbind requests the destructions of a service instance binding from associated sub-broker
func (subway *Broker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	subway.Logger.Info("unbind", lager.Data{
		"instance-id": instanceID,
		"binding-id":  bindingID,
	})

	for _, backendBroker := range subway.BackendBrokers {
		// Dummy URI to generate test results
		if backendBroker.URI == "TEST-FOUND-INSTANCE" {
			return nil
		} else if backendBroker.URI == "TEST-UNKNOWN-INSTANCE" {
			// Skip test backend broker
		} else {
			client := &http.Client{}
			url := fmt.Sprintf("%s/v2/service_instances/%s/service_bindings/%s?plan_id=%s&service_id=%s",
				backendBroker.URI, instanceID, bindingID, details.PlanID, details.ServiceID)

			req, err := http.NewRequest("DELETE", url, nil)
			if err != nil {
				subway.Logger.Error("backend-unbind-req", err)
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.SetBasicAuth(backendBroker.Username, backendBroker.Password)

			resp, err := client.Do(req)
			if err != nil {
				subway.Logger.Error("backend-unbind-resp", err)
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				subway.Logger.Info("unbind-success", lager.Data{
					"instance-id": instanceID,
					"binding-id":  bindingID,
					"backend-uri": backendBroker.URI,
				})
				return nil
			}
		}
	}

	return brokerapi.ErrInstanceDoesNotExist
}

// Deprovision requests the destruction of a service instance from associated sub-broker
func (subway *Broker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, acceptsIncomplete bool) (doesAcceptIncomplete bool, err error) {
	subway.Logger.Info("deprovision", lager.Data{
		"instance-id": instanceID,
	})

	for _, backendBroker := range subway.BackendBrokers {
		// Dummy URI to generate test results
		if backendBroker.URI == "TEST-FOUND-INSTANCE" {
			return false, nil
		} else if backendBroker.URI == "TEST-UNKNOWN-INSTANCE" {
			// Skip test backend broker
		} else {
			client := &http.Client{}
			url := fmt.Sprintf("%s/v2/service_instances/%s?plan_id=%s&service_id=%s",
				backendBroker.URI, instanceID, details.PlanID, details.ServiceID)
			req, err := http.NewRequest("DELETE", url, nil)
			if err != nil {
				subway.Logger.Error("backend-deprovision-req", err)
				return acceptsIncomplete, err
			}
			req.Header.Set("Content-Type", "application/json")
			req.SetBasicAuth(backendBroker.Username, backendBroker.Password)

			resp, err := client.Do(req)
			if err != nil {
				subway.Logger.Error("backend-deprovision-resp", err)
				return false, err
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				subway.Logger.Info("deprovision-success", lager.Data{
					"instance-id": instanceID,
					"backend-uri": backendBroker.URI,
				})
				return false, nil
			}
		}
	}

	return false, brokerapi.ErrInstanceDoesNotExist
}

// LastOperation returns the status of the last instance operation
func (subway *Broker) LastOperation(instanceID string) (resp brokerapi.LastOperationResponse, err error) {
	return resp, fmt.Errorf("Async not supported yet")
}
