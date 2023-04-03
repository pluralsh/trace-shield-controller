/*
Copyright 2023.

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

package v1alpha1

import (
	crhelperTypes "github.com/pluralsh/controller-reconcile-helper/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TenantSpec defines the desired state of Tenant
type TenantSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// DisplayName is a human readable name for the tenant
	DisplayName string `json:"displayName,omitempty"`

	// Limits is the set of limits for the tenant
	Limits LimitSpec `json:"limits,omitempty"`
}

// Defines the limits for a tenant
type LimitSpec struct {
	//The set of limits for the tenant on Mimir

	// +kubebuilder:validation:Optional
	Mimir Mimir `json:"mimir,omitempty"`
}

// TenantStatus defines the observed state of Tenant
type TenantStatus struct {
	// Conditions defines current service state of the PacketMachine.
	// +optional
	Conditions crhelperTypes.Conditions `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:path=tenants,scope=Cluster

// Tenant is the Schema for the tenants API
type Tenant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TenantSpec   `json:"spec,omitempty"`
	Status TenantStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TenantList contains a list of Tenant
type TenantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Tenant `json:"items"`
}

// GetConditions returns the list of conditions for a WireGuardServer API object.
func (t *Tenant) GetConditions() crhelperTypes.Conditions {
	return t.Status.Conditions
}

// SetConditions will set the given conditions on a WireGuardServer object.
func (t *Tenant) SetConditions(conditions crhelperTypes.Conditions) {
	t.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&Tenant{}, &TenantList{})
}
