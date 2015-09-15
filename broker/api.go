package broker

import (
	"errors"

	"github.com/pivotal-cf/brokerapi"
)

// Services is used by Cloud Foundry to learn the available catalog of services
func (broker *Broker) Services() []brokerapi.Service {
	return broker.Catalog
}

// Provision requests the creation of a service instance from an available sub-broker
func (broker *Broker) Provision(instanceID string, details brokerapi.ProvisionDetails) error {
	if details.PlanID == "" {
		return errors.New("plan_id required")
	}

	planID := ""
	for _, plan := range broker.plans() {
		if plan.ID == details.PlanID {
			planID = details.PlanID
			break
		}
	}

	if planID == "" {
		return errors.New("plan_id not recognized")
	}

	return broker.routeProvision(instanceID, details)
}

// Deprovision requests the destruction of a service instance from associated sub-broker
func (broker *Broker) Deprovision(instanceID string) error {
	// Deprovision instances here
	return nil
}

// Bind requests the creation of a service instance bindings from associated sub-broker
func (broker *Broker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (interface{}, error) {
	// Bind to instances here
	// Return credentials which will be marshalled to JSON
	credentialsMap := map[string]interface{}{}
	return credentialsMap, nil
}

// Unbind requests the destructions of a service instance binding from associated sub-broker
func (broker *Broker) Unbind(instanceID, bindingID string) error {
	// Unbind from instances here
	return nil
}
