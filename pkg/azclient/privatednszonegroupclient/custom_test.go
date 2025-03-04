// /*
// Copyright The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

// Code generated by client-gen. DO NOT EDIT.
package privatednszonegroupclient

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	armnetwork "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	armprivatedns "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns"
	. "github.com/onsi/gomega"
)

var (
	networkClientFactory     *armnetwork.ClientFactory
	privateDNSFactory        *armprivatedns.ClientFactory
	loadBalancersClient      *armnetwork.LoadBalancersClient
	pipClient                *armnetwork.PublicIPAddressesClient
	vnetClient               *armnetwork.VirtualNetworksClient
	privatelinkserviceClient *armnetwork.PrivateLinkServicesClient
	privateendpointclient    *armnetwork.PrivateEndpointsClient
	privatezoneclient        *armprivatedns.PrivateZonesClient
	vnetLinkClient           *armprivatedns.VirtualNetworkLinksClient
)
var (
	pipName            string = "pip1"
	pipResource        *armnetwork.PublicIPAddress
	lbName             string = "lb1"
	lbResource         *armnetwork.LoadBalancer
	vnetName           string = "vnet1"
	vnetResource       *armnetwork.VirtualNetwork
	plsName            string = "pls1"
	plsResource        *armnetwork.PrivateLinkService
	privateEndpoint    *armnetwork.PrivateEndpoint
	privateZone        *armprivatedns.PrivateZone
	virtualnetworklink *armprivatedns.VirtualNetworkLink
)

func init() {
	additionalTestCases = func() {
	}

	beforeAllFunc = func(ctx context.Context) {
		networkClientOption := clientOption
		networkClientOption.Telemetry.ApplicationID = "ccm-network-client"
		networkClientFactory, err = armnetwork.NewClientFactory(subscriptionID, recorder.TokenCredential(), &arm.ClientOptions{
			ClientOptions: networkClientOption,
		})
		Expect(err).NotTo(HaveOccurred())
		dnsClientOption := clientOption
		dnsClientOption.Telemetry.ApplicationID = "ccm-network-client"
		privateDNSFactory, err = armprivatedns.NewClientFactory(subscriptionID, recorder.TokenCredential(), &arm.ClientOptions{
			ClientOptions: dnsClientOption,
		})
		Expect(err).NotTo(HaveOccurred())
		pipClient = networkClientFactory.NewPublicIPAddressesClient()
		poller, err := pipClient.BeginCreateOrUpdate(ctx, resourceGroupName, pipName, armnetwork.PublicIPAddress{
			Location: to.Ptr(location),
			Properties: &armnetwork.PublicIPAddressPropertiesFormat{
				PublicIPAddressVersion:   to.Ptr(armnetwork.IPVersionIPv4),
				PublicIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodStatic),
			},
			SKU: &armnetwork.PublicIPAddressSKU{
				Name: to.Ptr(armnetwork.PublicIPAddressSKUNameStandard),
				Tier: to.Ptr(armnetwork.PublicIPAddressSKUTierRegional),
			},
		}, nil)
		Expect(err).NotTo(HaveOccurred())
		resp, err := poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: 1 * time.Second})
		Expect(err).NotTo(HaveOccurred())
		pipResource = &resp.PublicIPAddress
		loadBalancersClient = networkClientFactory.NewLoadBalancersClient()
		lbpoller, err := loadBalancersClient.BeginCreateOrUpdate(ctx, resourceGroupName, lbName, armnetwork.LoadBalancer{
			Location: to.Ptr(location),
			Properties: &armnetwork.LoadBalancerPropertiesFormat{
				FrontendIPConfigurations: []*armnetwork.FrontendIPConfiguration{
					{
						Name: to.Ptr("frontendConfig1"),
						Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
							PublicIPAddress: pipResource,
						},
					},
				},
			},
			SKU: &armnetwork.LoadBalancerSKU{
				Name: to.Ptr(armnetwork.LoadBalancerSKUNameStandard),
				Tier: to.Ptr(armnetwork.LoadBalancerSKUTierRegional),
			},
		}, nil)
		Expect(err).NotTo(HaveOccurred())
		lbresp, err := lbpoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: 1 * time.Second})
		Expect(err).NotTo(HaveOccurred())
		lbResource = &lbresp.LoadBalancer

		vnetClient = networkClientFactory.NewVirtualNetworksClient()
		vnetpoller, err := vnetClient.BeginCreateOrUpdate(ctx, resourceGroupName, vnetName, armnetwork.VirtualNetwork{
			Location: to.Ptr(location),
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						to.Ptr("10.1.0.0/16"),
					},
				},
				Subnets: []*armnetwork.Subnet{
					{
						Name: to.Ptr("subnet1"),
						Properties: &armnetwork.SubnetPropertiesFormat{
							AddressPrefix:                     to.Ptr("10.1.0.0/24"),
							PrivateEndpointNetworkPolicies:    to.Ptr(armnetwork.VirtualNetworkPrivateEndpointNetworkPoliciesDisabled),
							PrivateLinkServiceNetworkPolicies: to.Ptr(armnetwork.VirtualNetworkPrivateLinkServiceNetworkPoliciesDisabled),
						},
					},
				},
			},
		}, nil)
		Expect(err).NotTo(HaveOccurred())
		vnetresp, err := vnetpoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: 1 * time.Second})
		Expect(err).NotTo(HaveOccurred())
		vnetResource = &vnetresp.VirtualNetwork

		privatelinkserviceClient = networkClientFactory.NewPrivateLinkServicesClient()
		plsPoller, err := privatelinkserviceClient.BeginCreateOrUpdate(ctx, resourceGroupName, plsName, armnetwork.PrivateLinkService{
			Location: to.Ptr(location),
			Properties: &armnetwork.PrivateLinkServiceProperties{
				IPConfigurations: []*armnetwork.PrivateLinkServiceIPConfiguration{
					{
						Name: to.Ptr("ipConfig1"),
						Properties: &armnetwork.PrivateLinkServiceIPConfigurationProperties{
							Subnet:                  vnetResource.Properties.Subnets[0],
							Primary:                 to.Ptr(true),
							PrivateIPAddressVersion: to.Ptr(armnetwork.IPVersionIPv4),
						},
					},
				},
				LoadBalancerFrontendIPConfigurations: lbResource.Properties.FrontendIPConfigurations,
			},
		}, nil)
		Expect(err).NotTo(HaveOccurred())
		plsresp, err := plsPoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: 1 * time.Second})
		Expect(err).NotTo(HaveOccurred())
		plsResource = &plsresp.PrivateLinkService

		privateEndpoint = &armnetwork.PrivateEndpoint{
			Location: to.Ptr(location),
			Properties: &armnetwork.PrivateEndpointProperties{
				Subnet: vnetResource.Properties.Subnets[0],
				PrivateLinkServiceConnections: []*armnetwork.PrivateLinkServiceConnection{
					{
						Name: to.Ptr("plsConnection1"),
						Properties: &armnetwork.PrivateLinkServiceConnectionProperties{
							PrivateLinkServiceID: plsResource.ID,
						},
					},
				},
			},
		}
		privateendpointclient = networkClientFactory.NewPrivateEndpointsClient()
		peendointPoller, err := privateendpointclient.BeginCreateOrUpdate(ctx, resourceGroupName, resourceName, *privateEndpoint, nil)
		Expect(err).NotTo(HaveOccurred())
		peresp, err := peendointPoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: 1 * time.Second})
		Expect(err).NotTo(HaveOccurred())
		privateEndpoint = &peresp.PrivateEndpoint
		privateendpointName = *privateEndpoint.Name
		privateZone = &armprivatedns.PrivateZone{
			Name:       to.Ptr("privatezone1.testzone.local"),
			Location:   to.Ptr("global"),
			Properties: &armprivatedns.PrivateZoneProperties{},
		}
		privatezoneclient = privateDNSFactory.NewPrivateZonesClient()
		privateZonePoller, err := privatezoneclient.BeginCreateOrUpdate(ctx, resourceGroupName, *privateZone.Name, *privateZone, nil)
		Expect(err).NotTo(HaveOccurred())
		pzresp, err := privateZonePoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: 1 * time.Second})
		Expect(err).NotTo(HaveOccurred())
		privateZone = &pzresp.PrivateZone

		virtualnetworklink = &armprivatedns.VirtualNetworkLink{
			Name:     to.Ptr(resourceName),
			Location: to.Ptr("global"),
			Properties: &armprivatedns.VirtualNetworkLinkProperties{
				RegistrationEnabled: to.Ptr(true),
				VirtualNetwork: &armprivatedns.SubResource{
					ID: vnetResource.ID,
				},
			},
		}
		vnetLinkClient = privateDNSFactory.NewVirtualNetworkLinksClient()
		vnetLinkPoller, err := vnetLinkClient.BeginCreateOrUpdate(ctx, resourceGroupName, *privateZone.Name, resourceName, *virtualnetworklink, nil)
		Expect(err).NotTo(HaveOccurred())
		vnlresp, err := vnetLinkPoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: 1 * time.Second})
		Expect(err).NotTo(HaveOccurred())
		virtualnetworklink = &vnlresp.VirtualNetworkLink

		newResource = &armnetwork.PrivateDNSZoneGroup{
			Name: to.Ptr(resourceName),
			Properties: &armnetwork.PrivateDNSZoneGroupPropertiesFormat{
				PrivateDNSZoneConfigs: []*armnetwork.PrivateDNSZoneConfig{
					{
						Name: to.Ptr("zoneConfig1"),
						Properties: &armnetwork.PrivateDNSZonePropertiesFormat{
							PrivateDNSZoneID: privateZone.ID,
						},
					},
				},
			},
		}
	}
	afterAllFunc = func(ctx context.Context) {
		vnetLinkPoller, err := vnetLinkClient.BeginDelete(ctx, resourceGroupName, *privateZone.Name, *virtualnetworklink.Name, nil)
		Expect(err).NotTo(HaveOccurred())
		_, err = vnetLinkPoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: 1 * time.Second})
		Expect(err).NotTo(HaveOccurred())

		privateZonePoller, err := privatezoneclient.BeginDelete(ctx, resourceGroupName, *privateZone.Name, nil)
		Expect(err).NotTo(HaveOccurred())
		_, err = privateZonePoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: 1 * time.Second})
		Expect(err).NotTo(HaveOccurred())

		pepoller, err := privateendpointclient.BeginDelete(ctx, resourceGroupName, *privateEndpoint.Name, nil)
		Expect(err).NotTo(HaveOccurred())
		_, err = pepoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: 1 * time.Second})
		Expect(err).NotTo(HaveOccurred())

		plsPoller, err := privatelinkserviceClient.BeginDelete(ctx, resourceGroupName, plsName, nil)
		Expect(err).NotTo(HaveOccurred())
		_, err = plsPoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: 1 * time.Second})
		Expect(err).NotTo(HaveOccurred())
		lbPoller, err := loadBalancersClient.BeginDelete(ctx, resourceGroupName, lbName, nil)
		Expect(err).NotTo(HaveOccurred())
		_, err = lbPoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: 1 * time.Second})
		Expect(err).NotTo(HaveOccurred())

		pipPoller, err := pipClient.BeginDelete(ctx, resourceGroupName, pipName, nil)
		Expect(err).NotTo(HaveOccurred())
		_, err = pipPoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: 1 * time.Second})
		Expect(err).NotTo(HaveOccurred())

		vnetPoller, err := vnetClient.BeginDelete(ctx, resourceGroupName, vnetName, nil)
		Expect(err).NotTo(HaveOccurred())
		_, err = vnetPoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: 1 * time.Second})
		Expect(err).NotTo(HaveOccurred())
	}
}
