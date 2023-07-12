package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type LokiLimits struct {

	// Distributor enforced limits.
	// +kubebuilder:validation:Optional
	IngestionRateStrategy *string `yaml:"ingestion_rate_strategy,omitempty" json:"ingestion_rate_strategy,omitempty"`
	// +kubebuilder:validation:Optional
	IngestionRateMB *float64 `yaml:"ingestion_rate_mb,omitempty" json:"ingestion_rate_mb,omitempty"`
	// +kubebuilder:validation:Optional
	IngestionBurstSizeMB *float64 `yaml:"ingestion_burst_size_mb,omitempty" json:"ingestion_burst_size_mb,omitempty"`
	// +kubebuilder:validation:Optional
	MaxLabelNameLength *int `yaml:"max_label_name_length,omitempty" json:"max_label_name_length,omitempty"`
	// +kubebuilder:validation:Optional
	MaxLabelValueLength *int `yaml:"max_label_value_length,omitempty" json:"max_label_value_length,omitempty"`
	// +kubebuilder:validation:Optional
	MaxLabelNamesPerSeries *int `yaml:"max_label_names_per_series,omitempty" json:"max_label_names_per_series,omitempty"`
	// +kubebuilder:validation:Optional
	RejectOldSamples *bool `yaml:"reject_old_samples,omitempty" json:"reject_old_samples,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	RejectOldSamplesMaxAge *metav1.Duration `yaml:"reject_old_samples_max_age,omitempty" json:"reject_old_samples_max_age,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	CreationGracePeriod *metav1.Duration `yaml:"creation_grace_period,omitempty" json:"creation_grace_period,omitempty"`
	// +kubebuilder:validation:Optional
	EnforceMetricName *bool `yaml:"enforce_metric_name,omitempty" json:"enforce_metric_name,omitempty"`
	// +kubebuilder:validation:Optional
	MaxLineSize *uint64 `yaml:"max_line_size,omitempty" json:"max_line_size,omitempty"`
	// +kubebuilder:validation:Optional
	MaxLineSizeTruncate *bool `yaml:"max_line_size_truncate,omitempty" json:"max_line_size_truncate,omitempty"`
	// +kubebuilder:validation:Optional
	IncrementDuplicateTimestamp *bool `yaml:"increment_duplicate_timestamp,omitempty" json:"increment_duplicate_timestamp,omitempty"`

	// Ingester enforced limits.
	// +kubebuilder:validation:Optional
	MaxLocalStreamsPerUser *int `yaml:"max_streams_per_user,omitempty" json:"max_streams_per_user,omitempty"`
	// +kubebuilder:validation:Optional
	MaxGlobalStreamsPerUser *int `yaml:"max_global_streams_per_user,omitempty" json:"max_global_streams_per_user,omitempty"`
	// +kubebuilder:validation:Optional
	UnorderedWrites *bool `yaml:"unordered_writes,omitempty" json:"unordered_writes,omitempty"`
	// +kubebuilder:validation:Optional
	PerStreamRateLimit *uint64 `yaml:"per_stream_rate_limit,omitempty" json:"per_stream_rate_limit,omitempty"`
	// +kubebuilder:validation:Optional
	PerStreamRateLimitBurst *uint64 `yaml:"per_stream_rate_limit_burst,omitempty" json:"per_stream_rate_limit_burst,omitempty"`

	// Querier enforced limits.
	// +kubebuilder:validation:Optional
	MaxChunksPerQuery *int `yaml:"max_chunks_per_query,omitempty" json:"max_chunks_per_query,omitempty"`
	// +kubebuilder:validation:Optional
	MaxQuerySeries *int `yaml:"max_query_series,omitempty" json:"max_query_series,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MaxQueryLookback *metav1.Duration `yaml:"max_query_lookback,omitempty" json:"max_query_lookback,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MaxQueryLength *metav1.Duration `yaml:"max_query_length,omitempty" json:"max_query_length,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MaxQueryRange *metav1.Duration `yaml:"max_query_range,omitempty" json:"max_query_range,omitempty"`
	// +kubebuilder:validation:Optional
	MaxQueryParallelism *int `yaml:"max_query_parallelism,omitempty" json:"max_query_parallelism,omitempty"`
	// +kubebuilder:validation:Optional
	TSDBMaxQueryParallelism *int `yaml:"tsdb_max_query_parallelism,omitempty" json:"tsdb_max_query_parallelism,omitempty"`
	// +kubebuilder:validation:Optional
	TSDBMaxBytesPerShard *uint64 `yaml:"tsdb_max_bytes_per_shard,omitempty" json:"tsdb_max_bytes_per_shard,omitempty"`
	// +kubebuilder:validation:Optional
	CardinalityLimit *int `yaml:"cardinality_limit,omitempty" json:"cardinality_limit,omitempty"`
	// +kubebuilder:validation:Optional
	MaxStreamsMatchersPerQuery *int `yaml:"max_streams_matchers_per_query,omitempty" json:"max_streams_matchers_per_query,omitempty"`
	// +kubebuilder:validation:Optional
	MaxConcurrentTailRequests *int `yaml:"max_concurrent_tail_requests,omitempty" json:"max_concurrent_tail_requests,omitempty"`
	// +kubebuilder:validation:Optional
	MaxEntriesLimitPerQuery *int `yaml:"max_entries_limit_per_query,omitempty" json:"max_entries_limit_per_query,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MaxCacheFreshness *metav1.Duration `yaml:"max_cache_freshness_per_query,omitempty" json:"max_cache_freshness_per_query,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MaxStatsCacheFreshness *metav1.Duration `yaml:"max_stats_cache_freshness,omitempty" json:"max_stats_cache_freshness,omitempty"`
	// +kubebuilder:validation:Optional
	MaxQueriersPerTenant *int `yaml:"max_queriers_per_tenant,omitempty" json:"max_queriers_per_tenant,omitempty"`
	// +kubebuilder:validation:Optional
	QueryReadyIndexNumDays *int `yaml:"query_ready_index_num_days,omitempty" json:"query_ready_index_num_days,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	QueryTimeout *metav1.Duration `yaml:"query_timeout,omitempty" json:"query_timeout,omitempty"`

	// Query frontend enforced limits. The default is actually parameterized by the queryrange config.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	QuerySplitDuration *metav1.Duration `yaml:"split_queries_by_interval,omitempty" json:"split_queries_by_interval,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MinShardingLookback *metav1.Duration `yaml:"min_sharding_lookback,omitempty" json:"min_sharding_lookback,omitempty"`
	// +kubebuilder:validation:Optional
	MaxQueryBytesRead *uint64 `yaml:"max_query_bytes_read,omitempty" json:"max_query_bytes_read,omitempty"`
	// +kubebuilder:validation:Optional
	MaxQuerierBytesRead *uint64 `yaml:"max_querier_bytes_read,omitempty" json:"max_querier_bytes_read,omitempty"`
	// +kubebuilder:validation:Optional
	VolumeEnabled *bool `yaml:"volume_enabled,omitempty" json:"volume_enabled,omitempty" doc:"description=Enable log-volume endpoints."`
	// +kubebuilder:validation:Optional
	VolumeMaxSeries *int `yaml:"volume_max_series,omitempty" json:"volume_max_series,omitempty" doc:"description=The maximum number of aggregated series in a log-volume response"`

	// Ruler defaults and limits.

	// TODO(dannyk): this setting is misnamed and probably deprecatable.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	RulerEvaluationDelay *metav1.Duration `yaml:"ruler_evaluation_delay_duration,omitempty" json:"ruler_evaluation_delay_duration,omitempty"`
	// +kubebuilder:validation:Optional
	RulerMaxRulesPerRuleGroup *int `yaml:"ruler_max_rules_per_rule_group,omitempty" json:"ruler_max_rules_per_rule_group,omitempty"`
	// +kubebuilder:validation:Optional
	RulerMaxRuleGroupsPerTenant *int `yaml:"ruler_max_rule_groups_per_tenant,omitempty" json:"ruler_max_rule_groups_per_tenant,omitempty"`
	// +kubebuilder:validation:Optional
	// TODO: fix type RulerAlertManagerConfig     *ruler_config.AlertManagerConfig `yaml:"ruler_alertmanager_config,omitempty" json:"ruler_alertmanager_config,omitempty" doc:"hidden"`
	// +kubebuilder:validation:Optional
	RulerTenantShardSize *int `yaml:"ruler_tenant_shard_size,omitempty" json:"ruler_tenant_shard_size,omitempty"`

	// TODO(dannyk): add HTTP client overrides (basic auth / tls config, etc)
	// Ruler remote-write limits.

	// this field is the inversion of the general remote_write.enabled because the zero value of a boolean is false,
	// and if it were ruler_remote_write_enabled, it would be impossible to know if the value was explicitly set or default
	// +kubebuilder:validation:Optional
	RulerRemoteWriteDisabled *bool `yaml:"ruler_remote_write_disabled,omitempty" json:"ruler_remote_write_disabled,omitempty" doc:"description=Disable recording rules remote-write."`

	// +kubebuilder:validation:Optional
	// TODO: fix type RulerRemoteWriteConfig map[string]config.RemoteWriteConfig `yaml:"ruler_remote_write_config,omitempty" json:"ruler_remote_write_config,omitempty" doc:"description=Configures global and per-tenant limits for remote write clients. A map with remote client id as key."`

	// TODO(dannyk): possible enhancement is to align this with rule group interval
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	RulerRemoteEvaluationTimeout *metav1.Duration `yaml:"ruler_remote_evaluation_timeout,omitempty" json:"ruler_remote_evaluation_timeout,omitempty" doc:"description=Timeout for a remote rule evaluation. Defaults to the value of 'querier.query-timeout'."`
	// +kubebuilder:validation:Optional
	RulerRemoteEvaluationMaxResponseSize *int64 `yaml:"ruler_remote_evaluation_max_response_size,omitempty" json:"ruler_remote_evaluation_max_response_size,omitempty" doc:"description=Maximum size (in bytes) of the allowable response size from a remote rule evaluation. Set to 0 to allow any response size (default)."`

	// Global and per tenant deletion mode
	// +kubebuilder:validation:Optional
	DeletionMode *string `yaml:"deletion_mode,omitempty" json:"deletion_mode,omitempty"`

	// Global and per tenant retention
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	RetentionPeriod *metav1.Duration `yaml:"retention_period,omitempty" json:"retention_period,omitempty"`
	// +kubebuilder:validation:Optional
	// TODO: fix type StreamRetention []StreamRetention `yaml:"retention_stream,omitempty" json:"retention_stream,omitempty" doc:"description=Per-stream retention to apply, if the retention is enable on the compactor side.\nExample:\n retention_stream:\n - selector: '{namespace=\"dev\"}'\n priority: 1\n period: 24h\n- selector: '{container=\"nginx\"}'\n priority: 1\n period: 744h\nSelector is a Prometheus labels matchers that will apply the 'period' retention only if the stream is matching. In case multiple stream are matching, the highest priority will be picked. If no rule is matched the 'retention_period' is used."`

	// +kubebuilder:validation:Optional
	// TODO: fix type ShardStreams *shardstreams.Config `yaml:"shard_streams,omitempty" json:"shard_streams,omitempty"`

	// +kubebuilder:validation:Optional
	// TODO: fix type BlockedQueries []*validation.BlockedQuery `yaml:"blocked_queries,omitempty" json:"blocked_queries,omitempty"`

	// +kubebuilder:validation:Optional
	RequiredLabels []string `yaml:"required_labels,omitempty" json:"required_labels,omitempty" doc:"description=Define a list of required selector labels."`
	// +kubebuilder:validation:Optional
	RequiredNumberLabels *int `yaml:"minimum_labels_number,omitempty" json:"minimum_labels_number,omitempty" doc:"description=Minimum number of label matchers a query should contain."`
	// +kubebuilder:validation:Optional
	IndexGatewayShardSize *int `yaml:"index_gateway_shard_size,omitempty" json:"index_gateway_shard_size,omitempty"`
}
