package eks_test

import (
	"github.com/aws/aws-sdk-go/aws"
	awseks "github.com/aws/aws-sdk-go/service/eks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("EKS API wrapper", func() {
	Describe("can update cluster configuration for logging", func() {
		var (
			ctl *ClusterProvider

			cfg *api.ClusterConfig
			err error

			clusterLogging []*awseks.LogSetup
		)

		BeforeEach(func() {
			p := mockprovider.NewMockProvider()
			ctl = &ClusterProvider{
				Provider: p,
			}

			cfg = api.NewClusterConfig()

			updateClusterConfigOutput := &awseks.UpdateClusterConfigOutput{
				Update: &awseks.Update{
					Id:   aws.String("u123"),
					Type: aws.String(awseks.UpdateTypeLoggingUpdate),
				},
			}

			p.MockEKS().On("UpdateClusterConfig", mock.MatchedBy(func(input *awseks.UpdateClusterConfigInput) bool {
				Expect(input.Logging).ToNot(BeNil())

				Expect(input.Logging.ClusterLogging[0].Enabled).ToNot(BeNil())
				Expect(input.Logging.ClusterLogging[1].Enabled).ToNot(BeNil())

				Expect(*input.Logging.ClusterLogging[0].Enabled).To(BeTrue())
				Expect(*input.Logging.ClusterLogging[1].Enabled).To(BeFalse())

				clusterLogging = input.Logging.ClusterLogging

				return true
			})).Return(updateClusterConfigOutput, nil)

			describeUpdateInput := &awseks.DescribeUpdateInput{}

			describeUpdateOutput := &awseks.DescribeUpdateOutput{
				Update: &awseks.Update{
					Id:     aws.String("u123"),
					Type:   aws.String(awseks.UpdateTypeLoggingUpdate),
					Status: aws.String(awseks.UpdateStatusSuccessful),
				},
			}

			p.MockEKS().On("DescribeUpdateRequest", mock.MatchedBy(func(input *awseks.DescribeUpdateInput) bool {
				*describeUpdateInput = *input
				return input.UpdateId != nil && *input.UpdateId == *describeUpdateOutput.Update.Id
			})).Return(p.Client.MockRequestForGivenOutput(describeUpdateInput, describeUpdateOutput), describeUpdateOutput)
		})

		It("should have no logging by default", func() {
			err = api.SetClusterConfigDefaults(cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.EnableLogging).To(BeEmpty())
		})


		It("should handle unknown facilities", func() {
			cfg.EnableLogging = []string{"anything", "anyOtherThing"}

			err = api.SetClusterConfigDefaults(cfg)
			Expect(err).To(HaveOccurred())
		})

		It("should expand [] to none", func() {
			cfg.EnableLogging = []string{}

			err = api.SetClusterConfigDefaults(cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.EnableLogging).To(BeEmpty())

			err = ctl.UpdateClusterConfigForLogging(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(clusterLogging[0].Types).To(BeEmpty())

			Expect(clusterLogging[1].Types).ToNot(BeEmpty())
			Expect(clusterLogging[1].Types).To(Equal(aws.StringSlice(api.SupportedLoggingFacilities())))
		})

		It("should expand '*' to _all_", func() {
			cfg.EnableLogging = []string{"*"}

			err = api.SetClusterConfigDefaults(cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.EnableLogging).To(Equal(api.SupportedLoggingFacilities()))


			err = ctl.UpdateClusterConfigForLogging(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(clusterLogging[0].Types).ToNot(BeEmpty())
			Expect(clusterLogging[0].Types).To(Equal(aws.StringSlice(cfg.EnableLogging)))

			Expect(clusterLogging[1].Types).To(BeEmpty())
		})

		It("should expand 'all' to _all_", func() {
			cfg.EnableLogging = []string{"all"}

			err = api.SetClusterConfigDefaults(cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.EnableLogging).To(Equal(api.SupportedLoggingFacilities()))

			Expect(clusterLogging[0].Types).ToNot(BeEmpty())
			Expect(clusterLogging[0].Types).To(Equal(aws.StringSlice(cfg.EnableLogging)))

			Expect(clusterLogging[1].Types).To(BeEmpty())
		})

		It("should enable some logging facilities and disable others", func() {
			cfg.EnableLogging = []string{"authenticator", "controllerManager"}

			err = api.SetClusterConfigDefaults(cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.EnableLogging).To(Equal([]string{"authenticator", "controllerManager"}))

			err = ctl.UpdateClusterConfigForLogging(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(clusterLogging[0].Types).ToNot(BeEmpty())
			Expect(clusterLogging[0].Types).To(Equal(aws.StringSlice(cfg.EnableLogging)))

			Expect(clusterLogging[1].Types).ToNot(BeEmpty())
			Expect(clusterLogging[1].Types).To(Equal(aws.StringSlice([]string{"api", "audit", "scheduler"})))
		})

		It("should enable some logging facilities and disable others", func() {
			cfg.EnableLogging = []string{"audit", "scheduler"}

			err = api.SetClusterConfigDefaults(cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.EnableLogging).To(Equal([]string{"audit", "scheduler"}))

			err = ctl.UpdateClusterConfigForLogging(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(clusterLogging[0].Types).ToNot(BeEmpty())
			Expect(clusterLogging[0].Types).To(Equal(aws.StringSlice(cfg.EnableLogging)))

			Expect(clusterLogging[1].Types).ToNot(BeEmpty())
			Expect(clusterLogging[1].Types).To(Equal(aws.StringSlice([]string{"api", "authenticator", "controllerManager"})))
		})
	})
})
