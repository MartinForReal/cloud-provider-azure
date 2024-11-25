package privatelinkserviceclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	armnetwork "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/cache"
	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/utils"
)

type repo struct {
	Interface
	cache cache.Resource[armnetwork.PrivateLinkService]
}

const (
	DefaultCacheTTL = 120 * time.Second
)

func NewRepo(
	client Interface,
	cacheTTL time.Duration,
) (Interface, error) {
	getter := func(ctx context.Context, key string) (*armnetwork.PrivateLinkService, error) {
		resourceGroup, resourceName := parseRepositoryCacheKey(key)
		if resourceName == "" {
			return nil, fmt.Errorf("missing resource name")
		}

		resource, err := client.Get(ctx, resourceGroup, resourceName, nil)
		found, err := utils.CheckResourceExistsFromAzcoreError(err)
		if err != nil {
			return nil, err
		}
		if !found {
			return nil, nil
		}
		return resource, nil
	}

	if cacheTTL == 0 {
		cacheTTL = DefaultCacheTTL
	}
	c, err := cache.NewTimedCache(cacheTTL, getter, false)
	if err != nil {
		return nil, fmt.Errorf("new resource cache: %w", err)
	}

	return &repo{
		Interface: client,
		cache:     c,
	}, nil
}

func (r *repo) Get(ctx context.Context, resourceGroup, name string, expand *string) (*armnetwork.PrivateLinkService, error) {
	key := getRepositoryCacheKey(resourceGroup, name)
	cachedResource, err := r.cache.Get(ctx, key, cache.CacheReadTypeDefault)
	if err != nil {
		return nil, fmt.Errorf("get resource: %w", err)
	}

	return cachedResource, nil
}

func (r *repo) CreateOrUpdate(ctx context.Context, resourceGroup, name string, resource armnetwork.PrivateLinkService) (*armnetwork.PrivateLinkService, error) {
	resp, err := r.Interface.CreateOrUpdate(ctx, resourceGroup, name, resource)
	if err != nil {
		return nil, fmt.Errorf("create or update resource: %w", err)
	}
	// clear cache
	_ = r.cache.Delete(getRepositoryCacheKey(resourceGroup, name))

	return resp, nil
}

const cacheKeySeparator = "/"

func getRepositoryCacheKey(resourceGroup, name string) string {
	return strings.ToLower(strings.Join([]string{resourceGroup, name}, cacheKeySeparator))
}
func parseRepositoryCacheKey(key string) (resourceGroup, name string) {
	result := strings.Split(key, cacheKeySeparator)
	if len(result) != 2 {
		return "", ""
	}
	return result[0], result[1]
}

func (r *repo) Delete(ctx context.Context, resourceGroup, resourceName string) error {
	if resourceName == "" {
		return fmt.Errorf("resource name is empty")
	}
	cacheKey := getRepositoryCacheKey(resourceGroup, resourceName)

	if err := r.Interface.Delete(ctx, resourceGroup, resourceName); err != nil {
		return fmt.Errorf("delete resource error: %w", err)
	}
	// clear cache
	_ = r.cache.Delete(cacheKey)

	return nil
}
