/*
Copyright 2020 The Kubernetes Authors.

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

package config

import (
	"strings"

	"github.com/Azure/go-autorest/autorest/azure"
	"k8s.io/klog/v2"

	"sigs.k8s.io/cloud-provider-azure/pkg/azclient"
	"sigs.k8s.io/cloud-provider-azure/pkg/consts"
)

// AzureAuthConfig holds auth related part of cloud config
type AzureAuthConfig struct {
	azclient.ARMClientConfig `json:",inline" yaml:",inline"`
	azclient.AzureAuthConfig `json:",inline" yaml:",inline"`

	// The ID of the Azure Subscription that the cluster is deployed in
	SubscriptionID string `json:"subscriptionId,omitempty" yaml:"subscriptionId,omitempty"`
	// IdentitySystem indicates the identity provider. Relevant only to hybrid clouds (Azure Stack).
	// Allowed values are 'azure_ad' (default), 'adfs'.
	IdentitySystem string `json:"identitySystem,omitempty" yaml:"identitySystem,omitempty"`

	// The ID of the Azure Subscription that the network resources are deployed in
	NetworkResourceSubscriptionID string `json:"networkResourceSubscriptionID,omitempty" yaml:"networkResourceSubscriptionID,omitempty"`
}

// UsesNetworkResourceInDifferentTenant determines whether the AzureAuthConfig indicates to use network resources in
// different AAD Tenant than those for the cluster. Return true when NetworkResourceTenantID is specified  and not equal
// to one defined in global configs
func (config *AzureAuthConfig) UsesNetworkResourceInDifferentTenant() bool {
	return len(config.NetworkResourceTenantID) > 0 && !strings.EqualFold(config.NetworkResourceTenantID, config.TenantID)
}

// UsesNetworkResourceInDifferentSubscription determines whether the AzureAuthConfig indicates to use network resources
// in different Subscription than those for the cluster. Return true when NetworkResourceSubscriptionID is specified
// and not equal to one defined in global configs
func (config *AzureAuthConfig) UsesNetworkResourceInDifferentSubscription() bool {
	return len(config.NetworkResourceSubscriptionID) > 0 && !strings.EqualFold(config.NetworkResourceSubscriptionID, config.SubscriptionID)
}

// ParseAzureEnvironment returns the azure environment.
// If 'resourceManagerEndpoint' is set, the environment is computed by querying the cloud's resource manager endpoint.
// Otherwise, a pre-defined Environment is looked up by name.
func ParseAzureEnvironment(cloudName, resourceManagerEndpoint, identitySystem string) (*azure.Environment, error) {
	var env azure.Environment
	var err error
	if resourceManagerEndpoint != "" {
		klog.V(4).Infof("Loading environment from resource manager endpoint: %s", resourceManagerEndpoint)
		nameOverride := azure.OverrideProperty{Key: azure.EnvironmentName, Value: cloudName}
		env, err = azure.EnvironmentFromURL(resourceManagerEndpoint, nameOverride)
		if err == nil {
			azureStackOverrides(&env, resourceManagerEndpoint, identitySystem)
		}
	} else if cloudName == "" {
		klog.V(4).Info("Using public cloud environment")
		env = azure.PublicCloud
	} else {
		klog.V(4).Infof("Using %s environment", cloudName)
		env, err = azure.EnvironmentFromName(cloudName)
	}
	return &env, err
}

// azureStackOverrides ensures that the Environment matches what AKSe currently generates for Azure Stack
func azureStackOverrides(env *azure.Environment, resourceManagerEndpoint, identitySystem string) {
	env.ManagementPortalURL = strings.Replace(resourceManagerEndpoint, "https://management.", "https://portal.", -1)
	env.ServiceManagementEndpoint = env.TokenAudience
	env.ResourceManagerVMDNSSuffix = strings.Replace(resourceManagerEndpoint, "https://management.", "cloudapp.", -1)
	env.ResourceManagerVMDNSSuffix = strings.TrimSuffix(env.ResourceManagerVMDNSSuffix, "/")
	if strings.EqualFold(identitySystem, consts.ADFSIdentitySystem) {
		env.ActiveDirectoryEndpoint = strings.TrimSuffix(env.ActiveDirectoryEndpoint, "/")
		env.ActiveDirectoryEndpoint = strings.TrimSuffix(env.ActiveDirectoryEndpoint, "adfs")
	}
}
