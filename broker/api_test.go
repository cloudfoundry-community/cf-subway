package broker_test

import (
	"github.com/cloudfoundry-community/cf-subway/broker"
	"github.com/frodenas/brokerapi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service broker", func() {
	var subway *broker.Broker

	BeforeEach(func() {
		subway = broker.NewBroker()
		subway.BackendCatalog = brokerapi.CatalogResponse{
			Services: []brokerapi.Service{
				{
					ID: "service1-id",
					Metadata: &brokerapi.ServiceMetadata{
						LongDescription: "long description",
						ImageURL:        "image.png",
					},
					Plans: []brokerapi.ServicePlan{{ID: "plan-uuid"}},
				},
			},
		}
		subway.BackendBrokers = []*broker.BackendBroker{{URI: "TEST-SUCCESS"}}
	})

	Describe(".Services", func() {
		It("returns service catalog", func() {
			catalog := subway.Services()
			Ω(len(catalog.Services)).To(Equal(1))
			service := catalog.Services[0]
			Ω(service.ID).To(Equal("service1-id"))
			Ω(service.Metadata.LongDescription).To(Equal("long description"))
			Ω(service.Metadata.ImageURL).To(Equal("image.png"))
		})
	})

	Describe(".Provision", func() {
		Context("when the plan is recognized", func() {
			It("creates an instance if first backend ok", func() {
				subway.BackendBrokers = []*broker.BackendBroker{{URI: "TEST-SUCCESS"}}
				_, _, err := subway.Provision("some-id", brokerapi.ProvisionDetails{PlanID: "plan-uuid"}, false)
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
				_, _, err := subway.Provision("some-id", brokerapi.ProvisionDetails{PlanID: "plan-uuid"}, false)
				Ω(err).ToNot(HaveOccurred())
			})

			It("fails to create an instance if no backend ok", func() {
				subway.BackendBrokers = []*broker.BackendBroker{
					{URI: "TEST-NO-CAPACITY"},
					{URI: "TEST-NO-CAPACITY"},
				}
				_, _, err := subway.Provision("some-id", brokerapi.ProvisionDetails{PlanID: "plan-uuid"}, false)
				Ω(err).To(HaveOccurred())
			})
		})

		Context("when the plan is recognized but no backend brokers", func() {
			It("creates an instance", func() {
				subway.BackendBrokers = nil
				_, _, err := subway.Provision("some-id", brokerapi.ProvisionDetails{PlanID: "plan-uuid"}, false)
				Ω(err).To(HaveOccurred())
				Ω(err.Error()).To(Equal("No backend broker available for plan"))
			})
		})

		Context("when the plan is not recognized", func() {
			It("creates an instance", func() {
				_, _, err := subway.Provision("service-id", brokerapi.ProvisionDetails{PlanID: "unknown-uuid"}, false)
				Ω(err).To(HaveOccurred())
			})
		})
	})

	Describe(".Bind", func() {
		It("one broker recognizes service instance", func() {
			subway.BackendBrokers = []*broker.BackendBroker{
				{URI: "TEST-UNKNOWN-INSTANCE"},
				{URI: "TEST-FOUND-INSTANCE"},
				{URI: "TEST-UNKNOWN-INSTANCE"},
			}
			bindResponse, err := subway.Bind("service-id", "bind-id", brokerapi.BindDetails{PlanID: "plan-uuid"})
			Ω(err).ToNot(HaveOccurred())
			var creds map[string]interface{}
			creds = bindResponse.Credentials.(map[string]interface{})
			Ω(creds["host"]).To(Equal("10.10.10.10"))
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

	Describe(".Unbind", func() {
		var details brokerapi.UnbindDetails

		BeforeEach(func() {
			details = brokerapi.UnbindDetails{
				PlanID:    "plan-id",
				ServiceID: "service-id",
			}
		})

		It("one broker recognizes service instance", func() {
			subway.BackendBrokers = []*broker.BackendBroker{
				{URI: "TEST-UNKNOWN-INSTANCE"},
				{URI: "TEST-FOUND-INSTANCE"},
				{URI: "TEST-UNKNOWN-INSTANCE"},
			}
			err := subway.Unbind("service-id", "bind-id", details)
			Ω(err).ToNot(HaveOccurred())
		})

		It("no broker recognizes service instance", func() {
			subway.BackendBrokers = []*broker.BackendBroker{
				{URI: "TEST-UNKNOWN-INSTANCE"},
				{URI: "TEST-UNKNOWN-INSTANCE"},
			}
			err := subway.Unbind("service-id", "bind-id", details)
			Ω(err).To(HaveOccurred())
		})

	})

	Describe(".Deprovision", func() {
		var details brokerapi.DeprovisionDetails

		BeforeEach(func() {
			details = brokerapi.DeprovisionDetails{
				PlanID:    "plan-id",
				ServiceID: "service-id",
			}
		})

		It("one broker recognizes service instance", func() {
			subway.BackendBrokers = []*broker.BackendBroker{
				{URI: "TEST-UNKNOWN-INSTANCE"},
				{URI: "TEST-FOUND-INSTANCE"},
				{URI: "TEST-UNKNOWN-INSTANCE"},
			}
			_, err := subway.Deprovision("service-id", details, false)
			Ω(err).ToNot(HaveOccurred())
		})

		It("no broker recognizes service instance", func() {
			subway.BackendBrokers = []*broker.BackendBroker{
				{URI: "TEST-UNKNOWN-INSTANCE"},
				{URI: "TEST-UNKNOWN-INSTANCE"},
			}
			_, err := subway.Deprovision("service-id", details, false)
			Ω(err).To(HaveOccurred())
		})

	})
})
