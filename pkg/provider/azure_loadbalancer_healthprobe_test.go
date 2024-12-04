/*
Copyright 2023 The Kubernetes Authors.

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

package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2022-07-01/network"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/cloud-provider-azure/pkg/consts"
)

// getTestProbes returns dualStack probes.
func getTestProbes(protocol, path string, interval, servicePort, probePort, numOfProbe *int32) map[bool][]network.Probe {
	return map[bool][]network.Probe{
		consts.IPVersionIPv4: {getTestProbe(protocol, path, interval, servicePort, probePort, numOfProbe, consts.IPVersionIPv4)},
		consts.IPVersionIPv6: {getTestProbe(protocol, path, interval, servicePort, probePort, numOfProbe, consts.IPVersionIPv6)},
	}
}

func getTestProbe(protocol, path string, interval, servicePort, probePort, numOfProbe *int32, isIPv6 bool) network.Probe {
	suffix := ""
	if isIPv6 {
		suffix = "-" + consts.IPVersionIPv6String
	}
	expectedProbes := network.Probe{
		Name: to.Ptr(fmt.Sprintf("atest1-TCP-%d", *servicePort) + suffix),
		ProbePropertiesFormat: &network.ProbePropertiesFormat{
			Protocol:          network.ProbeProtocol(protocol),
			Port:              probePort,
			IntervalInSeconds: interval,
			ProbeThreshold:    numOfProbe,
		},
	}
	if (strings.EqualFold(protocol, "Http") || strings.EqualFold(protocol, "Https")) && len(strings.TrimSpace(path)) > 0 {
		expectedProbes.RequestPath = to.Ptr(path)
	}
	return expectedProbes
}

// getDefaultTestProbes returns dualStack probes.
func getDefaultTestProbes(protocol, path string) map[bool][]network.Probe {
	return getTestProbes(protocol, path, to.Ptr(int32(5)), to.Ptr(int32(80)), to.Ptr(int32(10080)), to.Ptr(int32(2)))
}

func TestFindProbe(t *testing.T) {
	tests := []struct {
		msg           string
		existingProbe []network.Probe
		curProbe      network.Probe
		expected      bool
	}{
		{
			msg:      "empty existing probes should return false",
			expected: false,
		},
		{
			msg: "probe names match while ports don't should return false",
			existingProbe: []network.Probe{
				{
					Name: to.Ptr("httpProbe"),
					ProbePropertiesFormat: &network.ProbePropertiesFormat{
						Port: to.Ptr(int32(1)),
					},
				},
			},
			curProbe: network.Probe{
				Name: to.Ptr("httpProbe"),
				ProbePropertiesFormat: &network.ProbePropertiesFormat{
					Port: to.Ptr(int32(2)),
				},
			},
			expected: false,
		},
		{
			msg: "probe ports match while names don't should return false",
			existingProbe: []network.Probe{
				{
					Name: to.Ptr("probe1"),
					ProbePropertiesFormat: &network.ProbePropertiesFormat{
						Port: to.Ptr(int32(1)),
					},
				},
			},
			curProbe: network.Probe{
				Name: to.Ptr("probe2"),
				ProbePropertiesFormat: &network.ProbePropertiesFormat{
					Port: to.Ptr(int32(1)),
				},
			},
			expected: false,
		},
		{
			msg: "probe protocol don't match should return false",
			existingProbe: []network.Probe{
				{
					Name: to.Ptr("probe1"),
					ProbePropertiesFormat: &network.ProbePropertiesFormat{
						Port:     to.Ptr(int32(1)),
						Protocol: network.ProbeProtocolHTTP,
					},
				},
			},
			curProbe: network.Probe{
				Name: to.Ptr("probe1"),
				ProbePropertiesFormat: &network.ProbePropertiesFormat{
					Port:     to.Ptr(int32(1)),
					Protocol: network.ProbeProtocolTCP,
				},
			},
			expected: false,
		},
		{
			msg: "probe path don't match should return false",
			existingProbe: []network.Probe{
				{
					Name: to.Ptr("probe1"),
					ProbePropertiesFormat: &network.ProbePropertiesFormat{
						Port:        to.Ptr(int32(1)),
						RequestPath: to.Ptr("/path1"),
					},
				},
			},
			curProbe: network.Probe{
				Name: to.Ptr("probe1"),
				ProbePropertiesFormat: &network.ProbePropertiesFormat{
					Port:        to.Ptr(int32(1)),
					RequestPath: to.Ptr("/path2"),
				},
			},
			expected: false,
		},
		{
			msg: "probe interval don't match should return false",
			existingProbe: []network.Probe{
				{
					Name: to.Ptr("probe1"),
					ProbePropertiesFormat: &network.ProbePropertiesFormat{
						Port:              to.Ptr(int32(1)),
						RequestPath:       to.Ptr("/path"),
						IntervalInSeconds: to.Ptr(int32(5)),
					},
				},
			},
			curProbe: network.Probe{
				Name: to.Ptr("probe1"),
				ProbePropertiesFormat: &network.ProbePropertiesFormat{
					Port:              to.Ptr(int32(1)),
					RequestPath:       to.Ptr("/path"),
					IntervalInSeconds: to.Ptr(int32(10)),
				},
			},
			expected: false,
		},
		{
			msg: "probe match should return true",
			existingProbe: []network.Probe{
				{
					Name: to.Ptr("matchName"),
					ProbePropertiesFormat: &network.ProbePropertiesFormat{
						Port: to.Ptr(int32(1)),
					},
				},
			},
			curProbe: network.Probe{
				Name: to.Ptr("matchName"),
				ProbePropertiesFormat: &network.ProbePropertiesFormat{
					Port: to.Ptr(int32(1)),
				},
			},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			findResult := findProbe(test.existingProbe, test.curProbe)
			assert.Equal(t, test.expected, findResult)
		})
	}
}

func TestShouldKeepSharedProbe(t *testing.T) {
	testCases := []struct {
		desc        string
		service     *v1.Service
		lb          network.LoadBalancer
		wantLB      bool
		expected    bool
		expectedErr error
	}{
		{
			desc:     "When the lb.Probes is nil",
			service:  &v1.Service{},
			lb:       network.LoadBalancer{},
			expected: false,
		},
		{
			desc:    "When the lb.Probes is not nil but does not contain a probe with the name consts.SharedProbeName",
			service: &v1.Service{},
			lb: network.LoadBalancer{
				LoadBalancerPropertiesFormat: &network.LoadBalancerPropertiesFormat{
					Probes: &[]network.Probe{
						{
							Name: to.Ptr("notSharedProbe"),
						},
					},
				},
			},
			expected: false,
		},
		{
			desc:    "When the lb.Probes contains a probe with the name consts.SharedProbeName, but none of the LoadBalancingRules in the probe matches the service",
			service: &v1.Service{},
			lb: network.LoadBalancer{
				LoadBalancerPropertiesFormat: &network.LoadBalancerPropertiesFormat{
					Probes: &[]network.Probe{
						{
							Name: to.Ptr(consts.SharedProbeName),
							ProbePropertiesFormat: &network.ProbePropertiesFormat{
								LoadBalancingRules: &[]network.SubResource{},
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			desc: "When the lb.Probes contains a probe with the name consts.SharedProbeName, and at least one of the LoadBalancingRules in the probe does not match the service",
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: types.UID("uid"),
				},
			},
			lb: network.LoadBalancer{
				LoadBalancerPropertiesFormat: &network.LoadBalancerPropertiesFormat{
					Probes: &[]network.Probe{
						{
							Name: to.Ptr(consts.SharedProbeName),
							ID:   to.Ptr("id"),
							ProbePropertiesFormat: &network.ProbePropertiesFormat{
								LoadBalancingRules: &[]network.SubResource{
									{
										ID: to.Ptr("other"),
									},
									{
										ID: to.Ptr("auid"),
									},
								},
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			desc: "When wantLB is true and the shared probe mode is not turned on",
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: types.UID("uid"),
				},
			},
			lb: network.LoadBalancer{
				LoadBalancerPropertiesFormat: &network.LoadBalancerPropertiesFormat{
					Probes: &[]network.Probe{
						{
							Name: to.Ptr(consts.SharedProbeName),
							ID:   to.Ptr("id"),
							ProbePropertiesFormat: &network.ProbePropertiesFormat{
								LoadBalancingRules: &[]network.SubResource{
									{
										ID: to.Ptr("other"),
									},
									{
										ID: to.Ptr("auid"),
									},
								},
							},
						},
					},
				},
			},
			wantLB: true,
		},
		{
			desc: "When the lb.Probes contains a probe with the name consts.SharedProbeName, and all of the LoadBalancingRules in the probe match the service",
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: types.UID("uid"),
				},
			},
			lb: network.LoadBalancer{
				LoadBalancerPropertiesFormat: &network.LoadBalancerPropertiesFormat{
					Probes: &[]network.Probe{
						{
							Name: to.Ptr(consts.SharedProbeName),
							ID:   to.Ptr("id"),
							ProbePropertiesFormat: &network.ProbePropertiesFormat{
								LoadBalancingRules: &[]network.SubResource{
									{
										ID: to.Ptr("auid"),
									},
								},
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			desc:     "Edge cases such as when the service or LoadBalancer is nil",
			service:  nil,
			lb:       network.LoadBalancer{},
			expected: false,
		},
		{
			desc:    "Case: Invalid LoadBalancingRule ID format causing getLastSegment to return an error",
			service: &v1.Service{},
			lb: network.LoadBalancer{
				LoadBalancerPropertiesFormat: &network.LoadBalancerPropertiesFormat{
					Probes: &[]network.Probe{
						{
							Name: to.Ptr(consts.SharedProbeName),
							ID:   to.Ptr("id"),
							ProbePropertiesFormat: &network.ProbePropertiesFormat{
								LoadBalancingRules: &[]network.SubResource{
									{
										ID: to.Ptr(""),
									},
								},
							},
						},
					},
				},
			},
			expected:    false,
			expectedErr: fmt.Errorf("resource name was missing from identifier"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			az := GetTestCloud(gomock.NewController(t))
			var expectedProbes []network.Probe
			result, err := az.keepSharedProbe(tc.service, tc.lb, expectedProbes, tc.wantLB)
			assert.Equal(t, tc.expectedErr, err)
			if tc.expected {
				assert.Equal(t, 1, len(result))
			}
		})
	}
}
