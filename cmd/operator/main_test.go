// Copyright 2016 The prometheus-operator Authors
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

package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
)

var runs = 0

type mockKubernetesClient struct {
	kubernetes.Interface
	discovery discovery.DiscoveryInterface
}

func newMockKubernetesClient() kubernetes.Interface {
	return &mockKubernetesClient{
		discovery: &mockDiscoveryClient{},
	}
}

func (m *mockKubernetesClient) Discovery() discovery.DiscoveryInterface {
	return m.discovery
}

type mockDiscoveryClient struct {
	discovery.DiscoveryInterface
}

func (m *mockDiscoveryClient) ServerResourcesForGroupVersion(_ string) (*metav1.APIResourceList, error) {
	time.Sleep(10 * time.Second)

	runs++

	if runs >= 3 {
		return &metav1.APIResourceList{
			APIResources: []metav1.APIResource{
				metav1.APIResource{
					Name: "test",
				},
			},
		}, nil
	}

	return &metav1.APIResourceList{
		APIResources: []metav1.APIResource{
			metav1.APIResource{
				Name: "fail",
			},
		},
	}, nil
}

func TestWaitForCRDInstalled(t *testing.T) {
	ctx := context.Background()
	client := newMockKubernetesClient()

	installed, err := checkInstalledWithTimeout(ctx, client, storagev1.SchemeGroupVersion, "test", 5*time.Second)
	require.NoError(t, err)
	require.Equal(t, runs, 1)
	require.False(t, installed)

	installed, err = checkInstalledWithTimeout(ctx, client, storagev1.SchemeGroupVersion, "test", 50*time.Second)
	require.NoError(t, err)
	require.GreaterOrEqual(t, runs, 3)
	require.True(t, installed)
}
