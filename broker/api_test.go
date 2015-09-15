package broker_test

import (
	"github.com/cloudfoundry-community/cf-subway/broker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf/brokerapi"
)

var _ = Describe("Service broker", func() {
	var subway *broker.Broker

	BeforeEach(func() {
		subway = broker.NewBroker()
		subway.Catalog = []brokerapi.Service{
			{
				Plans: []brokerapi.ServicePlan{{ID: "plan-uuid"}},
			},
		}
		subway.BackendBrokers = []*broker.BackendBroker{{URI: "TEST-SUCCESS"}}
	})

	Describe(".Provision", func() {
		Context("when the plan is recognized", func() {
			It("creates an instance if first backend ok", func() {
				subway.BackendBrokers = []*broker.BackendBroker{{URI: "TEST-SUCCESS"}}
				err := subway.Provision("some-id", brokerapi.ProvisionDetails{PlanID: "plan-uuid"})
				Ω(err).ToNot(HaveOccurred())
			})

			It("creates an instance if one backend ok", func() {
				subway.BackendBrokers = []*broker.BackendBroker{
					{URI: "TEST-NO-CAPACITY"},
					{URI: "TEST-NO-CAPACITY"},
					{URI: "TEST-NO-CAPACITY"},
					{URI: "TEST-SUCCESS"},
					{URI: "TEST-NO-CAPACITY"},
				}
				err := subway.Provision("some-id", brokerapi.ProvisionDetails{PlanID: "plan-uuid"})
				Ω(err).ToNot(HaveOccurred())
			})

			It("fails to create an instance if no backend ok", func() {
				subway.BackendBrokers = []*broker.BackendBroker{
					{URI: "TEST-NO-CAPACITY"},
					{URI: "TEST-NO-CAPACITY"},
				}
				err := subway.Provision("some-id", brokerapi.ProvisionDetails{PlanID: "plan-uuid"})
				Ω(err).To(HaveOccurred())
			})
		})

		Context("when the plan is recognized but no backend brokers", func() {
			It("creates an instance", func() {
				subway.BackendBrokers = nil
				err := subway.Provision("some-id", brokerapi.ProvisionDetails{PlanID: "plan-uuid"})
				Ω(err).To(HaveOccurred())
				Ω(err.Error()).To(Equal("No backend broker available for plan"))
			})
		})

		Context("when the plan is not recognized", func() {
			It("creates an instance", func() {
				err := subway.Provision("service-id", brokerapi.ProvisionDetails{PlanID: "unknown-uuid"})
				Ω(err).To(HaveOccurred())
			})
		})
	})

	Describe(".Bind", func() {
		It("one broker recognizes servicesinstance", func() {
			subway.BackendBrokers = []*broker.BackendBroker{
				{URI: "TEST-UNKNOWN-INSTANCE"},
				{URI: "TEST-FOUND-INSTANCE"},
				{URI: "TEST-UNKNOWN-INSTANCE"},
			}
			creds, err := subway.Bind("service-id", "bind-id", brokerapi.BindDetails{PlanID: "plan-uuid"})
			Ω(err).ToNot(HaveOccurred())
			Ω(creds).ToNot(BeNil())
			credentials := creds.(brokerapi.BindingResponse).Credentials.(map[string]interface{})
			Ω(credentials["host"]).To(Equal("10.10.10.10"))
		})

		It("no broker recognizes service instance", func() {
			subway.BackendBrokers = []*broker.BackendBroker{
				{URI: "TEST-UNKNOWN-INSTANCE"},
				{URI: "TEST-UNKNOWN-INSTANCE"},
			}
			_, err := subway.Bind("service-id", "bind-id", brokerapi.BindDetails{PlanID: "plan-uuid"})
			Ω(err).To(HaveOccurred())
		})

	})
})
