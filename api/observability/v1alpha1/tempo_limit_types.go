package v1alpha1

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

type TempoLimits struct {
	// Distributor enforced limits.
	// +kubebuilder:validation:Optional
	IngestionRateStrategy *string `yaml:"ingestion_rate_strategy,omitempty" json:"ingestion_rate_strategy,omitempty"`
	// +kubebuilder:validation:Optional
	IngestionRateLimitBytes *int `yaml:"ingestion_rate_limit_bytes,omitempty" json:"ingestion_rate_limit_bytes,omitempty"`
	// +kubebuilder:validation:Optional
	IngestionBurstSizeBytes *int `yaml:"ingestion_burst_size_bytes,omitempty" json:"ingestion_burst_size_bytes,omitempty"`

	// Ingester enforced limits.
	// +kubebuilder:validation:Optional
	MaxLocalTracesPerUser *int `yaml:"max_traces_per_user,omitempty" json:"max_traces_per_user,omitempty"`
	// +kubebuilder:validation:Optional
	MaxGlobalTracesPerUser *int `yaml:"max_global_traces_per_user,omitempty" json:"max_global_traces_per_user,omitempty"`

	// Forwarders
	// +kubebuilder:validation:Optional
	Forwarders []string `yaml:"forwarders,omitempty" json:"forwarders,omitempty"`

	// Metrics-generator config
	// +kubebuilder:validation:Optional
	MetricsGeneratorRingSize *int `yaml:"metrics_generator_ring_size,omitempty" json:"metrics_generator_ring_size,omitempty"`
	// TODO: ensure the list only contains service-graphs and span-metrics
	// +kubebuilder:validation:Optional
	MetricsGeneratorProcessors []string `yaml:"metrics_generator_processors,omitempty" json:"metrics_generator_processors,omitempty"`
	// +kubebuilder:validation:Optional
	MetricsGeneratorMaxActiveSeries *uint32 `yaml:"metrics_generator_max_active_series,omitempty" json:"metrics_generator_max_active_series,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MetricsGeneratorCollectionInterval *metav1.Duration `yaml:"metrics_generator_collection_interval,omitempty" json:"metrics_generator_collection_interval,omitempty"`
	// +kubebuilder:validation:Optional
	MetricsGeneratorDisableCollection *bool `yaml:"metrics_generator_disable_collection,omitempty" json:"metrics_generator_disable_collection,omitempty"`
	// +kubebuilder:validation:Optional
	MetricsGeneratorForwarderQueueSize *int `yaml:"metrics_generator_forwarder_queue_size,omitempty" json:"metrics_generator_forwarder_queue_size,omitempty"`
	// +kubebuilder:validation:Optional
	MetricsGeneratorForwarderWorkers *int `yaml:"metrics_generator_forwarder_workers,omitempty" json:"metrics_generator_forwarder_workers,omitempty"`
	// +kubebuilder:validation:Optional
	MetricsGeneratorProcessorServiceGraphsHistogramBuckets []float64 `yaml:"metrics_generator_processor_service_graphs_histogram_buckets,omitempty" json:"metrics_generator_processor_service_graphs_histogram_buckets,omitempty"`
	// +kubebuilder:validation:Optional
	MetricsGeneratorProcessorServiceGraphsDimensions []string `yaml:"metrics_generator_processor_service_graphs_dimensions,omitempty" json:"metrics_generator_processor_service_graphs_dimensions,omitempty"`
	// +kubebuilder:validation:Optional
	MetricsGeneratorProcessorServiceGraphsPeerAttributes []string `yaml:"metrics_generator_processor_service_graphs_peer_attributes,omitempty" json:"metrics_generator_processor_service_graphs_peer_attributes,omitempty"`
	// +kubebuilder:validation:Optional
	MetricsGeneratorProcessorServiceGraphsEnableClientServerPrefix *bool `yaml:"metrics_generator_processor_service_graphs_enable_client_server_prefix,omitempty" json:"metrics_generator_processor_service_graphs_enable_client_server_prefix,omitempty"`
	// +kubebuilder:validation:Optional
	MetricsGeneratorProcessorSpanMetricsHistogramBuckets []float64 `yaml:"metrics_generator_processor_span_metrics_histogram_buckets,omitempty" json:"metrics_generator_processor_span_metrics_histogram_buckets,omitempty"`
	// +kubebuilder:validation:Optional
	MetricsGeneratorProcessorSpanMetricsDimensions []string `yaml:"metrics_generator_processor_span_metrics_dimensions,omitempty" json:"metrics_generator_processor_span_metrics_dimensions,omitempty"`
	// +kubebuilder:validation:Optional
	MetricsGeneratorProcessorSpanMetricsIntrinsicDimensions map[string]bool `yaml:"metrics_generator_processor_span_metrics_intrinsic_dimensions,omitempty" json:"metrics_generator_processor_span_metrics_intrinsic_dimensions,omitempty"`
	// +kubebuilder:validation:Optional
	MetricsGeneratorProcessorSpanMetricsFilterPolicies []FilterPolicy `yaml:"metrics_generator_processor_span_metrics_filter_policies,omitempty" json:"metrics_generator_processor_span_metrics_filter_policies,omitempty"`
	// +kubebuilder:validation:Optional
	MetricsGeneratorProcessorSpanMetricsDimensionMappings []DimensionMappings `yaml:"metrics_generator_processor_span_metrics_dimension_mappings,omitempty" json:"metrics_generator_processor_span_metrics_dimension_mapings,omitempty"`
	// +kubebuilder:validation:Optional
	MetricsGeneratorProcessorSpanMetricsEnableTargetInfo *bool `yaml:"metrics_generator_processor_span_metrics_enable_target_info,omitempty" json:"metrics_generator_processor_span_metrics_enable_target_info,omitempty"`
	// +kubebuilder:validation:Optional
	MetricsGeneratorProcessorLocalBlocksMaxLiveTraces *uint64 `yaml:"metrics_generator_processor_local_blocks_max_live_traces,omitempty" json:"metrics_generator_processor_local_blocks_max_live_traces,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MetricsGeneratorProcessorLocalBlocksMaxBlockDuration *metav1.Duration `yaml:"metrics_generator_processor_local_blocks_max_block_duration,omitempty" json:"metrics_generator_processor_local_blocks_max_block_duration,omitempty"`
	// +kubebuilder:validation:Optional
	MetricsGeneratorProcessorLocalBlocksMaxBlockBytes *uint64 `yaml:"metrics_generator_processor_local_blocks_max_block_bytes,omitempty" json:"metrics_generator_processor_local_blocks_max_block_bytes,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MetricsGeneratorProcessorLocalBlocksFlushCheckPeriod *metav1.Duration `yaml:"metrics_generator_processor_local_blocks_flush_check_period,omitempty" json:"metrics_generator_processor_local_blocks_flush_check_period,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MetricsGeneratorProcessorLocalBlocksTraceIdlePeriod *metav1.Duration `yaml:"metrics_generator_processor_local_blocks_trace_idle_period,omitempty" json:"metrics_generator_processor_local_blocks_trace_idle_period,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MetricsGeneratorProcessorLocalBlocksCompleteBlockTimeout *metav1.Duration `yaml:"metrics_generator_processor_local_blocks_complete_block_timeout,omitempty" json:"metrics_generator_processor_local_blocks_complete_block_timeout,omitempty"`

	// Compactor enforced limits.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	BlockRetention *metav1.Duration `yaml:"block_retention,omitempty" json:"block_retention,omitempty"`

	// Querier and Ingester enforced limits.
	// +kubebuilder:validation:Optional
	MaxBytesPerTagValuesQuery *int `yaml:"max_bytes_per_tag_values_query,omitempty" json:"max_bytes_per_tag_values_query,omitempty"`
	// +kubebuilder:validation:Optional
	MaxBlocksPerTagValuesQuery *int `yaml:"max_blocks_per_tag_values_query,omitempty" json:"max_blocks_per_tag_values_query,omitempty"`

	// QueryFrontend enforced limits
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MaxSearchDuration *metav1.Duration `yaml:"max_search_duration,omitempty" json:"max_search_duration,omitempty"`

	// MaxBytesPerTrace is enforced in the Ingester, Compactor, Querier (Search) and Serverless (Search). It
	//  is not used when doing a trace by id lookup.
	// +kubebuilder:validation:Optional
	MaxBytesPerTrace *int `yaml:"max_bytes_per_trace,omitempty" json:"max_bytes_per_trace,omitempty"`
}

type TempoLimitsInput TempoLimits

type FilterPolicy struct {
	// +kubebuilder:validation:Optional
	Include *PolicyMatch `yaml:"include,omitempty" json:"include,omitempty"`
	// +kubebuilder:validation:Optional
	Exclude *PolicyMatch `yaml:"exclude,omitempty" json:"exclude,omitempty"`
}

type MatchType string

const (
	Strict MatchType = "strict"
	Regex  MatchType = "regex"
)

var AllMatchType = []MatchType{
	Strict,
	Regex,
}

func (e MatchType) IsValid() bool {
	switch e {
	case Strict, Regex:
		return true
	}
	return false
}

func (e MatchType) String() string {
	return string(e)
}

func (e *MatchType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MatchType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MatchType", str)
	}
	return nil
}

func (e MatchType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PolicyMatch struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=strict;regex
	MatchType *MatchType `yaml:"match_type,omitempty" json:"match_type,omitempty"`
	// +kubebuilder:validation:Optional
	Attributes []MatchPolicyAttribute `yaml:"attributes,omitempty" json:"attributes,omitempty"`
}

type MatchPolicyAttribute struct {
	// +kubebuilder:validation:Optional
	Key *string `yaml:"key,omitempty" json:"key,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:validation:Type=object
	Value *WrappedMap `yaml:"value,omitempty" json:"value,omitempty"`
}

type WrappedMap struct {
	Object map[string]interface{} `yaml:",inline" json:",inline"`
}

// MarshalJSON defers JSON encoding to the wrapped map
func (w *WrappedMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.Object)
}

// UnmarshalJSON implements the Unmarshaler interface.
func (w *WrappedMap) UnmarshalJSON(data []byte) error {
	var out map[string]interface{}
	err := json.Unmarshal(data, &out)
	if err != nil {
		return err
	}
	w.Object = out
	return nil
}

func (w *WrappedMap) DeepCopy() *WrappedMap {
	if w != nil && w.Object != nil {
		return &WrappedMap{
			Object: runtime.DeepCopyJSON(w.Object),
		}
	}
	c := NewWrappedMap()
	return &c
}

// DeepCopyInto is an ~autogenerated~ deepcopy function, copying the receiver, writing into out. in must be non-nil.
// Works around https://github.com/kubernetes/code-generator/issues/50
func (w *WrappedMap) DeepCopyInto(out *WrappedMap) {
	bytes, err := json.Marshal(w.Object)
	if err != nil {
		// we assume that it marshals cleanly because otherwise the resource would not have been
		// created in the API server
		panic(err)
	}
	var clone map[string]interface{}
	err = json.Unmarshal(bytes, &clone)
	if err != nil {
		// we assume again optimistically because we just marshalled that the round trip works as well
		panic(err)
	}
	out.Object = clone
}

// NewWrappedMap returns an empty WrappedMap
func NewWrappedMap() WrappedMap {
	return WrappedMap{Object: map[string]interface{}{}}
}

type DimensionMappings struct {
	// +kubebuilder:validation:Optional
	Name *string `yaml:"name,omitempty" json:"name,omitempty"`
	// +kubebuilder:validation:Optional
	SourceLabel []string `yaml:"source_labels,omitempty" json:"source_labels,omitempty"`
	// +kubebuilder:validation:Optional
	Join *string `yaml:"join,omitempty" json:"join,omitempty"`
}
