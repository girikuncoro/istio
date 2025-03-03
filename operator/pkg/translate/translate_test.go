// Copyright Istio Authors
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

package translate

import (
	"fmt"
	"testing"

	"istio.io/api/operator/v1alpha1"
	"istio.io/istio/operator/pkg/name"
	"istio.io/istio/operator/pkg/util"
	"istio.io/istio/pkg/test/util/assert"
)

func Test_skipReplicaCountWithAutoscaleEnabled(t *testing.T) {
	const valuesWithHPAndReplicaCountFormat = `
values:
  pilot:
    autoscaleEnabled: %t
  gateways:
    istio-ingressgateway:
      autoscaleEnabled: %t
    istio-egressgateway:
      autoscaleEnabled: %t
components:
  pilot:
    k8s:
      replicaCount: 2
  ingressGateways:
    - name: istio-ingressgateway
      enabled: true
      k8s:
        replicaCount: 2
  egressGateways:
    - name: istio-egressgateway
      enabled: true
      k8s:
        replicaCount: 2
`

	cases := []struct {
		name       string
		component  name.ComponentName
		values     string
		expectSkip bool
	}{
		{
			name:       "hpa enabled for pilot without replicas",
			component:  name.PilotComponentName,
			values:     fmt.Sprintf(valuesWithHPAndReplicaCountFormat, false, false, false),
			expectSkip: false,
		},
		{
			name:       "hpa enabled for ingressgateway without replica",
			component:  name.IngressComponentName,
			values:     fmt.Sprintf(valuesWithHPAndReplicaCountFormat, false, false, false),
			expectSkip: false,
		},
		{
			name:       "hpa enabled for pilot without replicas",
			component:  name.EgressComponentName,
			values:     fmt.Sprintf(valuesWithHPAndReplicaCountFormat, false, false, false),
			expectSkip: false,
		},
		{
			name:       "hpa enabled for pilot with replicas",
			component:  name.PilotComponentName,
			values:     fmt.Sprintf(valuesWithHPAndReplicaCountFormat, true, false, false),
			expectSkip: true,
		},
		{
			name:       "hpa enabled for ingressgateway with replicass",
			component:  name.IngressComponentName,
			values:     fmt.Sprintf(valuesWithHPAndReplicaCountFormat, false, true, false),
			expectSkip: true,
		},
		{
			name:       "hpa enabled for egressgateway with replicas",
			component:  name.EgressComponentName,
			values:     fmt.Sprintf(valuesWithHPAndReplicaCountFormat, true, false, true),
			expectSkip: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var iop *v1alpha1.IstioOperatorSpec
			if tt.values != "" {
				iop = &v1alpha1.IstioOperatorSpec{}
				if err := util.UnmarshalWithJSONPB(tt.values, iop, true); err != nil {
					t.Fatal(err)
				}
			}

			got := skipReplicaCountWithAutoscaleEnabled(iop, tt.component)
			assert.Equal(t, tt.expectSkip, got)
		})
	}
}
