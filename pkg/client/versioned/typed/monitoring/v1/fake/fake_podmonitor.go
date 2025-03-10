// Copyright The prometheus-operator Authors
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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/client/applyconfiguration/monitoring/v1"
	typedmonitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned/typed/monitoring/v1"
	gentype "k8s.io/client-go/gentype"
)

// fakePodMonitors implements PodMonitorInterface
type fakePodMonitors struct {
	*gentype.FakeClientWithListAndApply[*v1.PodMonitor, *v1.PodMonitorList, *monitoringv1.PodMonitorApplyConfiguration]
	Fake *FakeMonitoringV1
}

func newFakePodMonitors(fake *FakeMonitoringV1, namespace string) typedmonitoringv1.PodMonitorInterface {
	return &fakePodMonitors{
		gentype.NewFakeClientWithListAndApply[*v1.PodMonitor, *v1.PodMonitorList, *monitoringv1.PodMonitorApplyConfiguration](
			fake.Fake,
			namespace,
			v1.SchemeGroupVersion.WithResource("podmonitors"),
			v1.SchemeGroupVersion.WithKind("PodMonitor"),
			func() *v1.PodMonitor { return &v1.PodMonitor{} },
			func() *v1.PodMonitorList { return &v1.PodMonitorList{} },
			func(dst, src *v1.PodMonitorList) { dst.ListMeta = src.ListMeta },
			func(list *v1.PodMonitorList) []*v1.PodMonitor { return gentype.ToPointerSlice(list.Items) },
			func(list *v1.PodMonitorList, items []*v1.PodMonitor) { list.Items = gentype.FromPointerSlice(items) },
		),
		fake,
	}
}
