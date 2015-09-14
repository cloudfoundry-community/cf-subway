package broker

import "github.com/pivotal-cf/brokerapi"

// Services is used by Cloud Foundry to learn the available catalog of services
func (broker *Broker) Services() []brokerapi.Service {
	return []brokerapi.Service{}
}

// Provision requests the creation of a service instance from an available sub-broker
func (broker *Broker) Provision(instanceID string, details brokerapi.ProvisionDetails) error {
	// Provision a new instance here
	return nil
}

// Deprovision requests the destruction of a service instance from associated sub-broker
func (broker *Broker) Deprovision(instanceID string) error {
	// Deprovision instances here
	return nil
}

// Provision requests the creation of a service instance bindings from associated sub-broker
func (broker *Broker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (interface{}, error) {
	// Bind to instances here
	// Return credentials which will be marshalled to JSON
	credentialsMap := map[string]interface{}{}
	return credentialsMap, nil
}

// Provision requests the destructions of a service instance binding from associated sub-broker
func (broker *Broker) Unbind(instanceID, bindingID string) error {
	// Unbind from instances here
	return nil
}
