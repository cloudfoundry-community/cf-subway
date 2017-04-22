package broker_test

import (
	"context"

	"github.com/cloudfoundry-community/cf-subway/broker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi"
)

var _ = Describe("Service broker", func() {
	var subway *broker.Broker

	BeforeEach(func() {
		subway = broker.NewBroker()
		subway.BackendCatalog = []brokerapi.Service{
			brokerapi.Service{
				ID: "service1-id",
				Metadata: &brokerapi.ServiceMetadata{
					LongDescription: "long description",
					ImageUrl:        "image.png",
				},
				Plans: []brokerapi.ServicePlan{{ID: "plan-uuid"}},
			},
		}
		subway.BackendBrokers = []*broker.BackendBroker{{URI: "TEST-SUCCESS"}}
	})

	Describe(".Services", func() {
		It("returns service catalog", func() {
			services := subway.Services(context.Background())
			Ω(len(services)).To(Equal(1))
			service := services[0]
			Ω(service.ID).To(Equal("service1-id"))
			Ω(service.Metadata.LongDescription).To(Equal("long description"))
			Ω(service.Metadata.ImageUrl).To(Equal("image.png"))
		})
	})

	Describe(".Provision", func() {
		Context("when the plan is recognized", func() {
			It("creates an instance if first backend ok", func() {
				subway.BackendBrokers = []*broker.BackendBroker{{URI: "TEST-SUCCESS"}}
				_, err := subway.Provision(context.Background(), "some-id", brokerapi.ProvisionDetails{PlanID: "plan-uuid"}, false)
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
				_, err := subway.Provision(context.Background(), "some-id", brokerapi.ProvisionDetails{PlanID: "plan-uuid"}, false)
				Ω(err).ToNot(HaveOccurred())
			})

			It("fails to create an instance if no backend ok", func() {
				subway.BackendBrokers = []*broker.BackendBroker{
					{URI: "TEST-NO-CAPACITY"},
					{URI: "TEST-NO-CAPACITY"},
				}
				_, err := subway.Provision(context.Background(), "some-id", brokerapi.ProvisionDetails{PlanID: "plan-uuid"}, false)
				Ω(err).To(HaveOccurred())
			})
		})

		Context("when the plan is recognized but no backend brokers", func() {
			It("creates an instance", func() {
				subway.BackendBrokers = nil
				_, err := subway.Provision(context.Background(), "some-id", brokerapi.ProvisionDetails{PlanID: "plan-uuid"}, false)
				Ω(err).To(HaveOccurred())
				Ω(err.Error()).To(Equal("No backend broker available for plan"))
			})
		})

		Context("when the plan is not recognized", func() {
			It("creates an instance", func() {
				_, err := subway.Provision(context.Background(), "service-id", brokerapi.ProvisionDetails{PlanID: "unknown-uuid"}, false)
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
			bindResponse, err := subway.Bind(context.Background(), "service-id", "bind-id", brokerapi.BindDetails{PlanID: "plan-uuid"})
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
			_, err := subway.Bind(context.Background(), "service-id", "bind-id", brokerapi.BindDetails{PlanID: "plan-uuid"})
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
			err := subway.Unbind(context.Background(), "service-id", "bind-id", details)
			Ω(err).ToNot(HaveOccurred())
		})

		It("no broker recognizes service instance", func() {
			subway.BackendBrokers = []*broker.BackendBroker{
				{URI: "TEST-UNKNOWN-INSTANCE"},
				{URI: "TEST-UNKNOWN-INSTANCE"},
			}
			err := subway.Unbind(context.Background(), "service-id", "bind-id", details)
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
			_, err := subway.Deprovision(context.Background(), "service-id", details, false)
			Ω(err).ToNot(HaveOccurred())
		})

		It("no broker recognizes service instance", func() {
			subway.BackendBrokers = []*broker.BackendBroker{
				{URI: "TEST-UNKNOWN-INSTANCE"},
				{URI: "TEST-UNKNOWN-INSTANCE"},
			}
			_, err := subway.Deprovision(context.Background(), "service-id", details, false)
			Ω(err).To(HaveOccurred())
		})

	})
})
