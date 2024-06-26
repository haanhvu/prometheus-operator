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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	v1alpha1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

// EC2SDConfigApplyConfiguration represents an declarative configuration of the EC2SDConfig type for use
// with apply.
type EC2SDConfigApplyConfiguration struct {
	Region          *string                `json:"region,omitempty"`
	AccessKey       *v1.SecretKeySelector  `json:"accessKey,omitempty"`
	SecretKey       *v1.SecretKeySelector  `json:"secretKey,omitempty"`
	RoleARN         *string                `json:"roleARN,omitempty"`
	RefreshInterval *monitoringv1.Duration `json:"refreshInterval,omitempty"`
	Port            *int                   `json:"port,omitempty"`
	Filters         *v1alpha1.Filters      `json:"filters,omitempty"`
}

// EC2SDConfigApplyConfiguration constructs an declarative configuration of the EC2SDConfig type for use with
// apply.
func EC2SDConfig() *EC2SDConfigApplyConfiguration {
	return &EC2SDConfigApplyConfiguration{}
}

// WithRegion sets the Region field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Region field is set to the value of the last call.
func (b *EC2SDConfigApplyConfiguration) WithRegion(value string) *EC2SDConfigApplyConfiguration {
	b.Region = &value
	return b
}

// WithAccessKey sets the AccessKey field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AccessKey field is set to the value of the last call.
func (b *EC2SDConfigApplyConfiguration) WithAccessKey(value v1.SecretKeySelector) *EC2SDConfigApplyConfiguration {
	b.AccessKey = &value
	return b
}

// WithSecretKey sets the SecretKey field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SecretKey field is set to the value of the last call.
func (b *EC2SDConfigApplyConfiguration) WithSecretKey(value v1.SecretKeySelector) *EC2SDConfigApplyConfiguration {
	b.SecretKey = &value
	return b
}

// WithRoleARN sets the RoleARN field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RoleARN field is set to the value of the last call.
func (b *EC2SDConfigApplyConfiguration) WithRoleARN(value string) *EC2SDConfigApplyConfiguration {
	b.RoleARN = &value
	return b
}

// WithRefreshInterval sets the RefreshInterval field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RefreshInterval field is set to the value of the last call.
func (b *EC2SDConfigApplyConfiguration) WithRefreshInterval(value monitoringv1.Duration) *EC2SDConfigApplyConfiguration {
	b.RefreshInterval = &value
	return b
}

// WithPort sets the Port field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Port field is set to the value of the last call.
func (b *EC2SDConfigApplyConfiguration) WithPort(value int) *EC2SDConfigApplyConfiguration {
	b.Port = &value
	return b
}

// WithFilters sets the Filters field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Filters field is set to the value of the last call.
func (b *EC2SDConfigApplyConfiguration) WithFilters(value v1alpha1.Filters) *EC2SDConfigApplyConfiguration {
	b.Filters = &value
	return b
}
