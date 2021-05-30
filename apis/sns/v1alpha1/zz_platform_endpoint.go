/*
Copyright 2021 The Crossplane Authors.

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

// Code generated by ack-generate. DO NOT EDIT.

package v1alpha1

import (
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// PlatformEndpointParameters defines the desired state of PlatformEndpoint
type PlatformEndpointParameters struct {
	// Region is which region the PlatformEndpoint will be created.
	// +kubebuilder:validation:Required
	Region string `json:"region"`
	// For a list of attributes, see SetEndpointAttributes (https://docs.aws.amazon.com/sns/latest/api/API_SetEndpointAttributes.html).
	Attributes map[string]*string `json:"attributes,omitempty"`
	// Arbitrary user data to associate with the endpoint. Amazon SNS does not use
	// this data. The data must be in UTF-8 format and less than 2KB.
	CustomUserData *string `json:"customUserData,omitempty"`
	// Unique identifier created by the notification service for an app on a device.
	// The specific name for Token will vary, depending on which notification service
	// is being used. For example, when using APNS as the notification service,
	// you need the device token. Alternatively, when using GCM (Firebase Cloud
	// Messaging) or ADM, the device token equivalent is called the registration
	// ID.
	// +kubebuilder:validation:Required
	Token                            *string `json:"token"`
	CustomPlatformEndpointParameters `json:",inline"`
}

// PlatformEndpointSpec defines the desired state of PlatformEndpoint
type PlatformEndpointSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       PlatformEndpointParameters `json:"forProvider"`
}

// PlatformEndpointObservation defines the observed state of PlatformEndpoint
type PlatformEndpointObservation struct {
	// EndpointArn returned from CreateEndpoint action.
	EndpointARN *string `json:"endpointARN,omitempty"`
}

// PlatformEndpointStatus defines the observed state of PlatformEndpoint.
type PlatformEndpointStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          PlatformEndpointObservation `json:"atProvider"`
}

// +kubebuilder:object:root=true

// PlatformEndpoint is the Schema for the PlatformEndpoints API
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,aws}
type PlatformEndpoint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PlatformEndpointSpec   `json:"spec,omitempty"`
	Status            PlatformEndpointStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PlatformEndpointList contains a list of PlatformEndpoints
type PlatformEndpointList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PlatformEndpoint `json:"items"`
}

// Repository type metadata.
var (
	PlatformEndpointKind             = "PlatformEndpoint"
	PlatformEndpointGroupKind        = schema.GroupKind{Group: Group, Kind: PlatformEndpointKind}.String()
	PlatformEndpointKindAPIVersion   = PlatformEndpointKind + "." + GroupVersion.String()
	PlatformEndpointGroupVersionKind = GroupVersion.WithKind(PlatformEndpointKind)
)

func init() {
	SchemeBuilder.Register(&PlatformEndpoint{}, &PlatformEndpointList{})
}
