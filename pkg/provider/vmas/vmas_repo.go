/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vmas

import (
	"context"
	"strings"
	"sync"
	"time"

	"k8s.io/utils/ptr"

	"sigs.k8s.io/cloud-provider-azure/pkg/consts"
)

type AvailabilitySetEntry struct {
	VMAS          *compute.AvailabilitySet
	ResourceGroup string
}

func (as *availabilitySet) newVMASCache() (azcache.Resource, error) {
	getter := func(ctx context.Context, _ string) (interface{}, error) {
		localCache := &sync.Map{}

		allResourceGroups, err := as.GetResourceGroups()
		if err != nil {
			return nil, err
		}

		for _, resourceGroup := range allResourceGroups.UnsortedList() {
			allAvailabilitySets, rerr := as.AvailabilitySetsClient.List(ctx, resourceGroup)
			if rerr != nil {
				klog.Errorf("AvailabilitySetsClient.List failed: %v", rerr)
				return nil, rerr.Error()
			}

			for i := range allAvailabilitySets {
				vmas := allAvailabilitySets[i]
				if strings.EqualFold(ptr.Deref(vmas.Name, ""), "") {
					klog.Warning("failed to get the name of the VMAS")
					continue
				}
				localCache.Store(ptr.Deref(vmas.Name, ""), &AvailabilitySetEntry{
					VMAS:          &vmas,
					ResourceGroup: resourceGroup,
				})
			}
		}

		return localCache, nil
	}

	if as.Config.AvailabilitySetsCacheTTLInSeconds == 0 {
		as.Config.AvailabilitySetsCacheTTLInSeconds = consts.VMASCacheTTLDefaultInSeconds
	}

	return azcache.NewTimedCache(time.Duration(as.Config.AvailabilitySetsCacheTTLInSeconds)*time.Second, getter, as.Cloud.Config.DisableAPICallCache)
}

// RefreshCaches invalidates and renew all related caches.
func (as *availabilitySet) RefreshCaches() error {
	var err error
	as.vmasCache, err = as.newVMASCache()
	if err != nil {
		klog.Errorf("as.RefreshCaches: failed to create or refresh VMAS cache: %s", err)
		return err
	}
	return nil
}
