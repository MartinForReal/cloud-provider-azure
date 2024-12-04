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
package virtualnetworklinkclient

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	armprivatedns "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns"

	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/metrics"
	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/utils"
)

type Client struct {
	*armprivatedns.VirtualNetworkLinksClient
	subscriptionID string
	tracer         tracing.Tracer
}

func New(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (Interface, error) {
	if options == nil {
		options = utils.GetDefaultOption()
	}
	tr := options.TracingProvider.NewTracer(utils.ModuleName, utils.ModuleVersion)

	client, err := armprivatedns.NewVirtualNetworkLinksClient(subscriptionID, credential, options)
	if err != nil {
		return nil, err
	}
	return &Client{
		VirtualNetworkLinksClient: client,
		subscriptionID:            subscriptionID,
		tracer:                    tr,
	}, nil
}

const GetOperationName = "VirtualNetworkLinksClient.Get"

// Get gets the VirtualNetworkLink
func (client *Client) Get(ctx context.Context, resourceGroupName string, parentResourceName string, resourceName string) (result *armprivatedns.VirtualNetworkLink, err error) {

	metricsCtx := metrics.BeginARMRequest(client.subscriptionID, resourceGroupName, "VirtualNetworkLink", "get")
	defer func() { metricsCtx.Observe(ctx, err) }()
	ctx, endSpan := runtime.StartSpan(ctx, GetOperationName, client.tracer, nil)
	defer endSpan(err)
	resp, err := client.VirtualNetworkLinksClient.Get(ctx, resourceGroupName, parentResourceName, resourceName, nil)
	if err != nil {
		return nil, err
	}
	//handle statuscode
	return &resp.VirtualNetworkLink, nil
}

const CreateOrUpdateOperationName = "VirtualNetworkLinksClient.Create"

// CreateOrUpdate creates or updates a VirtualNetworkLink.
func (client *Client) CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, parentResourceName string, resource armprivatedns.VirtualNetworkLink) (result *armprivatedns.VirtualNetworkLink, err error) {
	metricsCtx := metrics.BeginARMRequest(client.subscriptionID, resourceGroupName, "VirtualNetworkLink", "create_or_update")
	defer func() { metricsCtx.Observe(ctx, err) }()
	ctx, endSpan := runtime.StartSpan(ctx, CreateOrUpdateOperationName, client.tracer, nil)
	defer endSpan(err)
	resp, err := utils.NewPollerWrapper(client.VirtualNetworkLinksClient.BeginCreateOrUpdate(ctx, resourceGroupName, resourceName, parentResourceName, resource, nil)).WaitforPollerResp(ctx)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		return &resp.VirtualNetworkLink, nil
	}
	return nil, nil
}

const DeleteOperationName = "VirtualNetworkLinksClient.Delete"

// Delete deletes a VirtualNetworkLink by name.
func (client *Client) Delete(ctx context.Context, resourceGroupName string, parentResourceName string, resourceName string) (err error) {
	metricsCtx := metrics.BeginARMRequest(client.subscriptionID, resourceGroupName, "VirtualNetworkLink", "delete")
	defer func() { metricsCtx.Observe(ctx, err) }()
	ctx, endSpan := runtime.StartSpan(ctx, DeleteOperationName, client.tracer, nil)
	defer endSpan(err)
	_, err = utils.NewPollerWrapper(client.BeginDelete(ctx, resourceGroupName, parentResourceName, resourceName, nil)).WaitforPollerResp(ctx)
	return err
}
