package privatelinkserviceclient

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/cache"
	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/privatelinkserviceclient/mock_privatelinkserviceclient"
)

var _ = Describe("Repository", func() {
	var impl Interface
	var repository Interface
	var cntl *gomock.Controller
	BeforeEach(func() {
		var err error
		cntl = gomock.NewController(GinkgoT())
		impl = mock_privatelinkserviceclient.NewMockInterface(cntl)
		repository, err = NewRepo(impl, DefaultCacheTTL)
		Expect(err).To(BeNil())
	})
	AfterEach(func() {
		cntl.Finish()
	})
	When("Get is called", func() {
		It("should return the resource and put object into cache", func(ctx context.Context) {
			impl.(*mock_privatelinkserviceclient.MockInterface).EXPECT().Get(ctx, strings.ToLower(resourceGroupName), strings.ToLower(resourceName), nil).Return(newResource, nil)
			resource, err := repository.Get(ctx, resourceGroupName, resourceName, nil)
			Expect(err).To(BeNil())
			Expect(resource).NotTo(BeNil())
			Expect(len(repository.(*repo).cache.(*cache.TimedCache[armnetwork.PrivateLinkService]).GetStore().List())).To(Equal(1))

			//get from cache
			resource, err = repository.Get(ctx, resourceGroupName, resourceName, nil)
			Expect(err).To(BeNil())
			Expect(resource).NotTo(BeNil())
			Expect(len(repository.(*repo).cache.(*cache.TimedCache[armnetwork.PrivateLinkService]).GetStore().List())).To(Equal(1))
		})
	})
	When("CreateOrUpdate is called", func() {
		It("should delete cached items", func(ctx context.Context) {
			repository.(*repo).cache.(*cache.TimedCache[armnetwork.PrivateLinkService]).Set(getRepositoryCacheKey(resourceGroupName, resourceName), &armnetwork.PrivateLinkService{
				Name: to.Ptr(resourceName),
			})
			impl.(*mock_privatelinkserviceclient.MockInterface).EXPECT().CreateOrUpdate(gomock.Any(), resourceGroupName, resourceName, gomock.Any()).Return(newResource, nil)
			resource, err := repository.CreateOrUpdate(ctx, resourceGroupName, resourceName, armnetwork.PrivateLinkService{})
			Expect(err).To(BeNil())
			Expect(resource).NotTo(BeNil())
			Expect(repository.(*repo).cache.(*cache.TimedCache[armnetwork.PrivateLinkService]).GetStore().List()).To(BeEmpty())

		})
	})
	When("Delete is called", func() {
		It("should delete cached items", func(ctx context.Context) {
			repository.(*repo).cache.(*cache.TimedCache[armnetwork.PrivateLinkService]).Set(getRepositoryCacheKey(resourceGroupName, resourceName), &armnetwork.PrivateLinkService{
				Name: to.Ptr(resourceName),
			})
			impl.(*mock_privatelinkserviceclient.MockInterface).EXPECT().Delete(gomock.Any(), resourceGroupName, resourceName).Return(nil)
			err := repository.Delete(ctx, resourceGroupName, resourceName)
			Expect(err).To(BeNil())
			Expect(repository.(*repo).cache.(*cache.TimedCache[armnetwork.PrivateLinkService]).GetStore().List()).To(BeEmpty())
		})
	})
})
