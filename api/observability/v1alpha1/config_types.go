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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ConfigSpec defines the desired state of Config
type ConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	Mimir MimirSpec `json:"mimir"`

	// +kubebuilder:validation:Required
	Loki LokiSpec `json:"loki"`

	// +kubebuilder:validation:Required
	Tempo TempoSpec `json:"tempo"`
}

type MimirSpec struct {
	// +kubebuilder:validation:Required
	ConfigMap ConfigMapSelector `json:"configMap"`

	// +kubebuilder:validation:Optional
	Config MimirConfigSpec `json:"config,omitempty"`
}

type LokiSpec struct {
	// +kubebuilder:validation:Required
	ConfigMap ConfigMapSelector `json:"configMap"`

	// +kubebuilder:validation:Optional
	Config LokiConfigSpec `json:"config,omitempty"`
}

type TempoSpec struct {
	// +kubebuilder:validation:Required
	ConfigMap ConfigMapSelector `json:"configMap"`
}

type MimirConfigSpec struct {
	// +kubebuilder:validation:Optional
	Multi MultiRuntimeConfig `json:"multi_kv_config,omitempty"`

	// +kubebuilder:validation:Optional
	IngesterChunkStreaming *bool `json:"ingester_stream_chunks_when_using_blocks,omitempty"`

	// +kubebuilder:validation:Optional
	IngesterLimits *MimirIngesterInstanceLimits `json:"ingester_limits,omitempty"`
	// +kubebuilder:validation:Optional
	DistributorLimits *MimirDistributorInstanceLimits `json:"distributor_limits,omitempty"`
}

type LokiConfigSpec struct {
	// +kubebuilder:validation:Optional
	Multi MultiRuntimeConfig `json:"multi_kv_config,omitempty"`
	// +kubebuilder:validation:Optional
	TenantConfig map[string]*LokiRuntimeConfig `json:"configs,omitempty"`
}

type LokiRuntimeConfig struct {
	// +kubebuilder:validation:Optional
	LogStreamCreation bool `json:"log_stream_creation,omitempty"`
	// +kubebuilder:validation:Optional
	LogPushRequest bool `json:"log_push_request,omitempty"`
	// +kubebuilder:validation:Optional
	LogPushRequestStreams bool `json:"log_push_request_streams,omitempty"`
	// LimitedLogPushErrors is to be implemented and will allow logging push failures at a controlled pace.
	// +kubebuilder:validation:Optional
	LimitedLogPushErrors bool `json:"limited_log_push_errors,omitempty"`
}

type MultiRuntimeConfig struct {
	// Primary store used by MultiClient. Can be updated in runtime to switch to a different store (eg. consul -> etcd,
	// or to gossip). Doing this allows nice migration between stores. Empty values are ignored.
	PrimaryStore string `json:"primary"`

	// Mirroring enabled or not. Nil = no change.
	Mirroring *bool `json:"mirror_enabled"`
}

type MimirIngesterInstanceLimits struct {
	// +kubebuilder:validation:Optional
	MaxIngestionRate float64 `json:"max_ingestion_rate,omitempty"`
	// +kubebuilder:validation:Optional
	MaxInMemoryTenants int64 `json:"max_tenants,omitempty"`
	// +kubebuilder:validation:Optional
	MaxInMemorySeries int64 `json:"max_series,omitempty"`
	// +kubebuilder:validation:Optional
	MaxInflightPushRequests int64 `json:"max_inflight_push_requests,omitempty"`
}

type MimirDistributorInstanceLimits struct {
	// +kubebuilder:validation:Optional
	MaxIngestionRate float64 `json:"max_ingestion_rate,omitempty"`
	// +kubebuilder:validation:Optional
	MaxInflightPushRequests int `json:"max_inflight_push_requests,omitempty"`
	// +kubebuilder:validation:Optional
	MaxInflightPushRequestsBytes int `json:"max_inflight_push_requests_bytes,omitempty"`
}

type ConfigMapSelector struct {
	// +kubebuilder:default:="mimir-runtime"
	Name string `json:"name"`

	// +kubebuilder:default:="mimir"
	Namespace string `json:"namespace"`

	// +kubebuilder:default:="runtime.yaml"
	Key string `json:"key"`
}

// type KetoConfig struct {
// 	ReadRemote string `json:"readRemote,omitempty"`

// 	WriteRemote string `json:"writeRemote,omitempty"`
// }

// ConfigStatus defines the observed state of Config
type ConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:path=configs,scope=Cluster

// +genclient:nonNamespaced
// Config is the Schema for the configs API
type Config struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigSpec   `json:"spec,omitempty"`
	Status ConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConfigList contains a list of Config
type ConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Config `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Config{}, &ConfigList{})
}
