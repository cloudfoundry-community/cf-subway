package broker

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"code.cloudfoundry.org/lager"
)

// BackendBroker describes the location/creds for a backend broker providing actual services
type BackendBroker struct {
	URI      string
	Username string
	Password string
}

// LoadBackendBrokersFromEnv allows registration of backend brokers via environment variables
// BACKEND_BROKER_1=https://username1:password1@hostname1
// BACKEND_BROKER_2=https://username2:password2@hostname2
func (subway *Broker) LoadBackendBrokersFromEnv() {
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if strings.HasPrefix(pair[0], "BACKEND_BROKER") {
			backendURI := pair[1]
			uri, err := url.Parse(backendURI)
			if err != nil {
				subway.Logger.Error("backend-brokers", fmt.Errorf("Could not parse $%s %s", pair[0], backendURI))
				continue
			}
			if uri.User == nil {
				subway.Logger.Error("backend-brokers", fmt.Errorf("Missing username:password in URL: %s", backendURI))
				continue
			}
			password, _ := uri.User.Password()
			username := uri.User.Username()
			if username == "" || password == "" {
				subway.Logger.Error("backend-brokers", fmt.Errorf("Missing username:password in URL: %s", backendURI))
				continue
			}
			url := fmt.Sprintf("%s://%s", uri.Scheme, uri.Host)
			backendBroker := BackendBroker{url, username, password}
			subway.BackendBrokers = append(subway.BackendBrokers, &backendBroker)
			subway.Logger.Info("backend-brokers", lager.Data{"backend-broker": backendURI})
		}
	}
}
