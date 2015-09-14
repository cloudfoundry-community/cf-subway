package broker_test

import (
	"github.com/cloudfoundry-community/cf-subway/broker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf/brokerapi"
)

var _ = Describe("Service broker", func() {
	var broker broker.Broker

	BeforeEach(func() {
		broker.Catalog = []brokerapi.Service{
			{
				Plans: []brokerapi.ServicePlan{
					{ID: "plan-uuid"},
				},
			},
		}
	})

	Describe(".Provision", func() {
		Context("when the plan is recognized", func() {
			It("creates an instance", func() {
				err := broker.Provision("some-id", brokerapi.ProvisionDetails{PlanID: "plan-uuid"})
				Ω(err).ToNot(HaveOccurred())

				// Expect(len(someCreatorAndBinder.createdInstanceIds)).To(Equal(1))
				// Expect(someCreatorAndBinder.createdInstanceIds[0]).To(Equal(instanceID))
			})
		})

		Context("when the plan is not recognized", func() {
			It("creates an instance", func() {
				err := broker.Provision("some-id", brokerapi.ProvisionDetails{PlanID: "unknown-uuid"})
				Ω(err).To(HaveOccurred())

				// Expect(len(someCreatorAndBinder.createdInstanceIds)).To(Equal(1))
				// Expect(someCreatorAndBinder.createdInstanceIds[0]).To(Equal(instanceID))
			})
		})
	})
})
