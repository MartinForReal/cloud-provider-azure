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
package fileservicepropertiesclient

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	armstorage "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/accountclient"
	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/utils"
)

var (
	storageClientFactory *armstorage.ClientFactory
	storageaccountClient *armstorage.AccountsClient
	storageAccount       *armstorage.Account
)

func init() {
	additionalTestCases = func() {
		When("get requests are raised", func() {
			It("should not return error", func(ctx context.Context) {
				newResource, err := realClient.Get(ctx, resourceGroupName, resourceName)
				Expect(err).NotTo(HaveOccurred())
				Expect(newResource).NotTo(BeNil())
			})
		})
		When("invalid get requests are raised", func() {
			It("should return 404 error", func(ctx context.Context) {
				newResource, err := realClient.Get(ctx, resourceGroupName, resourceName+"notfound")
				Expect(err).To(HaveOccurred())
				Expect(newResource).To(BeNil())
			})
		})
		When("set requests are raised", func() {
			It("should not return error", func(ctx context.Context) {
				newResource, err := realClient.Set(ctx, resourceGroupName, resourceName, armstorage.FileServiceProperties{
					FileServiceProperties: &armstorage.FileServicePropertiesProperties{
						ShareDeleteRetentionPolicy: &armstorage.DeleteRetentionPolicy{
							Enabled: to.Ptr(true),
							Days:    to.Ptr(int32(1)),
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(newResource).NotTo(BeNil())
			})
		})
		When("invalid set requests are raised", func() {
			It("should return 404 error", func(ctx context.Context) {
				newResource, err := realClient.Set(ctx, resourceGroupName, resourceName+"notfound", armstorage.FileServiceProperties{
					FileServiceProperties: &armstorage.FileServicePropertiesProperties{
						ShareDeleteRetentionPolicy: &armstorage.DeleteRetentionPolicy{
							Enabled: to.Ptr(true),
							Days:    to.Ptr(int32(1)),
						},
					},
				})
				Expect(err).To(HaveOccurred())
				Expect(newResource).To(BeNil())
			})
		})
	}

	beforeAllFunc = func(ctx context.Context) {
		storageClientOption := clientOption
		storageClientOption.Telemetry.ApplicationID = "ccm-storage-client"
		if location == "chinaeast2" {
			storageClientOption.APIVersion = accountclient.MooncakeApiVersion
		}
		storageClientFactory, err = armstorage.NewClientFactory(subscriptionID, recorder.TokenCredential(), &arm.ClientOptions{
			ClientOptions: storageClientOption,
		})
		Expect(err).NotTo(HaveOccurred())
		storageaccountClient = storageClientFactory.NewAccountsClient()
		resourceName = "akscitfilesdktest"
		storageAccount, err := utils.NewPollerWrapper(storageaccountClient.BeginCreate(ctx, resourceGroupName, resourceName, armstorage.AccountCreateParameters{
			Location: to.Ptr(location),
			Kind:     to.Ptr(armstorage.KindStorageV2),
			Properties: &armstorage.AccountPropertiesCreateParameters{
				DNSEndpointType:              to.Ptr(armstorage.DNSEndpointTypeStandard),
				DefaultToOAuthAuthentication: to.Ptr(false),
				AllowBlobPublicAccess:        to.Ptr(false),
				AllowCrossTenantReplication:  to.Ptr(false),
				IsHnsEnabled:                 to.Ptr(true),
				MinimumTLSVersion:            to.Ptr(armstorage.MinimumTLSVersionTLS12),
				AllowSharedKeyAccess:         to.Ptr(true),
				PublicNetworkAccess:          to.Ptr(armstorage.PublicNetworkAccessDisabled),
				IsLocalUserEnabled:           to.Ptr(true),
				LargeFileSharesState:         to.Ptr(armstorage.LargeFileSharesStateEnabled),
				IsSftpEnabled:                to.Ptr(true),
				EnableNfsV3:                  to.Ptr(true),
				EnableHTTPSTrafficOnly:       to.Ptr(true),
				NetworkRuleSet: &armstorage.NetworkRuleSet{
					Bypass:        to.Ptr(armstorage.BypassAzureServices),
					DefaultAction: to.Ptr(armstorage.DefaultActionDeny),
					IPRules:       []*armstorage.IPRule{},
				},
				Encryption: &armstorage.Encryption{
					RequireInfrastructureEncryption: to.Ptr(false),
					KeySource:                       to.Ptr(armstorage.KeySourceMicrosoftStorage),
					Services: &armstorage.EncryptionServices{
						File: &armstorage.EncryptionService{
							KeyType: to.Ptr(armstorage.KeyTypeAccount),
							Enabled: to.Ptr(true),
						},
						Blob: &armstorage.EncryptionService{
							KeyType: to.Ptr(armstorage.KeyTypeAccount),
							Enabled: to.Ptr(true),
						},
					},
				},
				AccessTier: to.Ptr(armstorage.AccessTierCool),
			},
			SKU: &armstorage.SKU{
				Name: to.Ptr(armstorage.SKUNameStandardLRS),
				Tier: to.Ptr(armstorage.SKUTierStandard),
			},
		}, nil)).WaitforPollerResp(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(storageAccount).NotTo(BeNil())
		Expect(*storageAccount.Name).To(Equal(resourceName))
	}
	afterAllFunc = func(ctx context.Context) {
		_, err = storageaccountClient.Delete(ctx, resourceGroupName, resourceName, nil)
		Expect(err).NotTo(HaveOccurred())
	}
}
