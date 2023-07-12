package v1alpha1

import (
	crhelperTypes "github.com/pluralsh/controller-reconcile-helper/pkg/types"
	prom_v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MimirLimits struct {
	// Distributor enforced limits.
	// +kubebuilder:validation:Optional
	RequestRate *float64 `yaml:"request_rate,omitempty" json:"request_rate,omitempty" category:"experimental"`
	// +kubebuilder:validation:Optional
	RequestBurstSize *int `yaml:"request_burst_size,omitempty" json:"request_burst_size,omitempty" category:"experimental"`
	// +kubebuilder:validation:Optional
	IngestionRate *float64 `yaml:"ingestion_rate,omitempty" json:"ingestion_rate,omitempty"`
	// +kubebuilder:validation:Optional
	IngestionBurstSize *int `yaml:"ingestion_burst_size,omitempty" json:"ingestion_burst_size,omitempty"`
	// +kubebuilder:validation:Optional
	AcceptHASamples *bool `yaml:"accept_ha_samples,omitempty" json:"accept_ha_samples,omitempty"`
	// +kubebuilder:validation:Optional
	HAClusterLabel *string `yaml:"ha_cluster_label,omitempty" json:"ha_cluster_label,omitempty"`
	// +kubebuilder:validation:Optional
	HAReplicaLabel *string `yaml:"ha_replica_label,omitempty" json:"ha_replica_label,omitempty"`
	// +kubebuilder:validation:Optional
	HAMaxClusters *int `yaml:"ha_max_clusters,omitempty" json:"ha_max_clusters,omitempty"`
	// +kubebuilder:validation:Optional
	DropLabels []*string `yaml:"drop_labels,omitempty" json:"drop_labels,omitempty" category:"advanced"`
	// +kubebuilder:validation:Optional
	MaxLabelNameLength *int `yaml:"max_label_name_length,omitempty" json:"max_label_name_length,omitempty"`
	// +kubebuilder:validation:Optional
	MaxLabelValueLength *int `yaml:"max_label_value_length,omitempty" json:"max_label_value_length,omitempty"`
	// +kubebuilder:validation:Optional
	MaxLabelNamesPerSeries *int `yaml:"max_label_names_per_series,omitempty" json:"max_label_names_per_series,omitempty"`
	// +kubebuilder:validation:Optional
	MaxMetadataLength *int `yaml:"max_metadata_length,omitempty" json:"max_metadata_length,omitempty"`
	// +kubebuilder:validation:Optional
	MaxNativeHistogramBuckets *int `yaml:"max_native_histogram_buckets,omitempty" json:"max_native_histogram_buckets,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	CreationGracePeriod *metav1.Duration `yaml:"creation_grace_period,omitempty" json:"creation_grace_period,omitempty" category:"advanced"`
	// +kubebuilder:validation:Optional
	EnforceMetadataMetricName *bool `yaml:"enforce_metadata_metric_name,omitempty" json:"enforce_metadata_metric_name,omitempty" category:"advanced"`
	// +kubebuilder:validation:Optional
	IngestionTenantShardSize *int `yaml:"ingestion_tenant_shard_size,omitempty" json:"ingestion_tenant_shard_size,omitempty"`
	// +kubebuilder:validation:Optional
	MetricRelabelConfigs []prom_v1.RelabelConfig `yaml:"metric_relabel_configs,omitempty" json:"metric_relabel_configs,omitempty" doc:"nocli|description=List of metric relabel configurations. Note that in most situations, it is more effective to use metrics relabeling directly in the Prometheus server, e.g. remote_write.write_relabel_configs." category:"experimental"`

	// Ingester enforced limits.
	// Series
	// +kubebuilder:validation:Optional
	MaxGlobalSeriesPerUser *int `yaml:"max_global_series_per_user,omitempty" json:"max_global_series_per_user,omitempty"`
	// +kubebuilder:validation:Optional
	MaxGlobalSeriesPerMetric *int `yaml:"max_global_series_per_metric,omitempty" json:"max_global_series_per_metric,omitempty"`
	// Metadata
	// +kubebuilder:validation:Optional
	MaxGlobalMetricsWithMetadataPerUser *int `yaml:"max_global_metadata_per_user,omitempty" json:"max_global_metadata_per_user,omitempty"`
	// +kubebuilder:validation:Optional
	MaxGlobalMetadataPerMetric *int `yaml:"max_global_metadata_per_metric,omitempty" json:"max_global_metadata_per_metric,omitempty"`
	// +kubebuilder:validation:Optional
	// Exemplars
	// +kubebuilder:validation:Optional
	MaxGlobalExemplarsPerUser *int `yaml:"max_global_exemplars_per_user,omitempty" json:"max_global_exemplars_per_user,omitempty" category:"experimental"`
	// Native histograms
	// +kubebuilder:validation:Optional
	NativeHistogramsIngestionEnabled *bool `yaml:"native_histograms_ingestion_enabled,omitempty" json:"native_histograms_ingestion_enabled,omitempty" category:"experimental"`
	// Active series custom trackers
	// +kubebuilder:validation:Optional
	ActiveSeriesCustomTrackersConfig map[string]string `yaml:"active_series_custom_trackers,omitempty" json:"active_series_custom_trackers,omitempty" doc:"description=Additional custom trackers for active metrics. If there are active series matching a provided matcher (map value), the count will be exposed in the custom trackers metric labeled using the tracker name (map key). Zero valued counts are not exposed (and removed when they go back to zero)." category:"advanced"`
	// Max allowed time window for out-of-order samples.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	OutOfOrderTimeWindow *metav1.Duration `yaml:"out_of_order_time_window,omitempty" json:"out_of_order_time_window,omitempty" category:"experimental"`
	// +kubebuilder:validation:Optional
	OutOfOrderBlocksExternalLabelEnabled *bool `yaml:"out_of_order_blocks_external_label_enabled,omitempty" json:"out_of_order_blocks_external_label_enabled,omitempty" category:"experimental"`
	// +kubebuilder:validation:Optional

	// User defined label to give the option of subdividing specific metrics by another label
	// +kubebuilder:validation:Optional
	SeparateMetricsGroupLabel *string `yaml:"separate_metrics_group_label,omitempty" json:"separate_metrics_group_label,omitempty" category:"experimental"`

	// Querier enforced limits.
	// +kubebuilder:validation:Optional
	MaxChunksPerQuery *int `yaml:"max_fetched_chunks_per_query,omitempty" json:"max_fetched_chunks_per_query,omitempty"`
	// +kubebuilder:validation:Optional
	MaxFetchedSeriesPerQuery *int `yaml:"max_fetched_series_per_query,omitempty" json:"max_fetched_series_per_query,omitempty"`
	// +kubebuilder:validation:Optional
	MaxFetchedChunkBytesPerQuery *int `yaml:"max_fetched_chunk_bytes_per_query,omitempty" json:"max_fetched_chunk_bytes_per_query,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MaxQueryLookback *metav1.Duration `yaml:"max_query_lookback,omitempty" json:"max_query_lookback,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MaxPartialQueryLength *metav1.Duration `yaml:"max_partial_query_length,omitempty" json:"max_partial_query_length,omitempty"`
	// +kubebuilder:validation:Optional
	MaxQueryParallelism *int `yaml:"max_query_parallelism,omitempty" json:"max_query_parallelism,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MaxLabelsQueryLength *metav1.Duration `yaml:"max_labels_query_length,omitempty" json:"max_labels_query_length,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MaxCacheFreshness *metav1.Duration `yaml:"max_cache_freshness,omitempty" json:"max_cache_freshness,omitempty" category:"advanced"`
	// +kubebuilder:validation:Optional
	MaxQueriersPerTenant *int `yaml:"max_queriers_per_tenant,omitempty" json:"max_queriers_per_tenant,omitempty"`
	// +kubebuilder:validation:Optional
	QueryShardingTotalShards *int `yaml:"query_sharding_total_shards,omitempty" json:"query_sharding_total_shards,omitempty"`
	// +kubebuilder:validation:Optional
	QueryShardingMaxShardedQueries *int `yaml:"query_sharding_max_sharded_queries,omitempty" json:"query_sharding_max_sharded_queries,omitempty"`
	// +kubebuilder:validation:Optional
	QueryShardingMaxRegexpSizeBytes *int `yaml:"query_sharding_max_regexp_size_bytes,omitempty" json:"query_sharding_max_regexp_size_bytes,omitempty" category:"experimental"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	SplitInstantQueriesByInterval *metav1.Duration `yaml:"split_instant_queries_by_interval,omitempty" json:"split_instant_queries_by_interval,omitempty" category:"experimental"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	QueryIngestersWithin *metav1.Duration `yaml:"query_ingesters_within,omitempty" json:"query_ingesters_within,omitempty" category:"advanced"`

	// Query-frontend limits.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MaxTotalQueryLength *metav1.Duration `yaml:"max_total_query_length,omitempty" json:"max_total_query_length,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	ResultsCacheTTL *metav1.Duration `yaml:"results_cache_ttl,omitempty" json:"results_cache_ttl,omitempty" category:"experimental"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	ResultsCacheTTLForOutOfOrderTimeWindow *metav1.Duration `yaml:"results_cache_ttl_for_out_of_order_time_window,omitempty" json:"results_cache_ttl_for_out_of_order_time_window,omitempty" category:"experimental"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	ResultsCacheTTLForCardinalityQuery *metav1.Duration `yaml:"results_cache_ttl_for_cardinality_query,omitempty" json:"results_cache_ttl_for_cardinality_query,omitempty" category:"experimental"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	ResultsCacheTTLForLabelsQuery *metav1.Duration `yaml:"results_cache_ttl_for_labels_query,omitempty" json:"results_cache_ttl_for_labels_query,omitempty" category:"experimental"`
	// +kubebuilder:validation:Optional
	ResultsCacheForUnalignedQueryEnabled *bool `yaml:"cache_unaligned_requests,omitempty" json:"cache_unaligned_requests,omitempty" category:"advanced"`
	// +kubebuilder:validation:Optional
	MaxQueryExpressionSizeBytes *int `yaml:"max_query_expression_size_bytes,omitempty" json:"max_query_expression_size_bytes,omitempty" category:"experimental"`

	// Cardinality
	// +kubebuilder:validation:Optional
	CardinalityAnalysisEnabled *bool `yaml:"cardinality_analysis_enabled,omitempty" json:"cardinality_analysis_enabled,omitempty"`
	// +kubebuilder:validation:Optional
	LabelNamesAndValuesResultsMaxSizeBytes *int `yaml:"label_names_and_values_results_max_size_bytes,omitempty" json:"label_names_and_values_results_max_size_bytes,omitempty"`
	// +kubebuilder:validation:Optional
	LabelValuesMaxCardinalityLabelNamesPerRequest *int `yaml:"label_values_max_cardinality_label_names_per_request,omitempty" json:"label_values_max_cardinality_label_names_per_request,omitempty"`

	// Ruler defaults and limits.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	RulerEvaluationDelay *metav1.Duration `yaml:"ruler_evaluation_delay_duration,omitempty" json:"ruler_evaluation_delay_duration,omitempty"`
	// +kubebuilder:validation:Optional
	RulerTenantShardSize *int `yaml:"ruler_tenant_shard_size,omitempty" json:"ruler_tenant_shard_size,omitempty"`
	// +kubebuilder:validation:Optional
	RulerMaxRulesPerRuleGroup *int `yaml:"ruler_max_rules_per_rule_group,omitempty" json:"ruler_max_rules_per_rule_group,omitempty"`
	// +kubebuilder:validation:Optional
	RulerMaxRuleGroupsPerTenant *int `yaml:"ruler_max_rule_groups_per_tenant,omitempty" json:"ruler_max_rule_groups_per_tenant,omitempty"`
	// +kubebuilder:validation:Optional
	RulerRecordingRulesEvaluationEnabled *bool `yaml:"ruler_recording_rules_evaluation_enabled,omitempty" json:"ruler_recording_rules_evaluation_enabled,omitempty" category:"experimental"`
	// +kubebuilder:validation:Optional
	RulerAlertingRulesEvaluationEnabled *bool `yaml:"ruler_alerting_rules_evaluation_enabled,omitempty" json:"ruler_alerting_rules_evaluation_enabled,omitempty" category:"experimental"`
	// +kubebuilder:validation:Optional
	RulerSyncRulesOnChangesEnabled *bool `yaml:"ruler_sync_rules_on_changes_enabled,omitempty" json:"ruler_sync_rules_on_changes_enabled,omitempty" category:"advanced"`

	// Store-gateway.
	// +kubebuilder:validation:Optional
	StoreGatewayTenantShardSize *int `yaml:"store_gateway_tenant_shard_size,omitempty" json:"store_gateway_tenant_shard_size,omitempty"`

	// Compactor.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	CompactorBlocksRetentionPeriod *metav1.Duration `yaml:"compactor_blocks_retention_period,omitempty" json:"compactor_blocks_retention_period,omitempty"`
	// +kubebuilder:validation:Optional
	CompactorSplitAndMergeShards *int `yaml:"compactor_split_and_merge_shards,omitempty" json:"compactor_split_and_merge_shards,omitempty"`
	// +kubebuilder:validation:Optional
	CompactorSplitGroups *int `yaml:"compactor_split_groups,omitempty" json:"compactor_split_groups,omitempty"`
	// +kubebuilder:validation:Optional
	CompactorTenantShardSize *int `yaml:"compactor_tenant_shard_size,omitempty" json:"compactor_tenant_shard_size,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	CompactorPartialBlockDeletionDelay *metav1.Duration `yaml:"compactor_partial_block_deletion_delay,omitempty" json:"compactor_partial_block_deletion_delay,omitempty"`
	// +kubebuilder:validation:Optional
	CompactorBlockUploadEnabled *bool `yaml:"compactor_block_upload_enabled,omitempty" json:"compactor_block_upload_enabled,omitempty"`
	// +kubebuilder:validation:Optional
	CompactorBlockUploadValidationEnabled *bool `yaml:"compactor_block_upload_validation_enabled,omitempty" json:"compactor_block_upload_validation_enabled,omitempty"`
	// +kubebuilder:validation:Optional
	CompactorBlockUploadVerifyChunks *bool `yaml:"compactor_block_upload_verify_chunks,omitempty" json:"compactor_block_upload_verify_chunks,omitempty"`
	// +kubebuilder:validation:Optional
	CompactorBlockUploadMaxBlockSizeBytes *int64 `yaml:"compactor_block_upload_max_block_size_bytes,omitempty" json:"compactor_block_upload_max_block_size_bytes,omitempty" category:"advanced"`

	// This config doesn't have a CLI flag registered here because they're registered in
	// their own original config struct.
	// +kubebuilder:validation:Optional
	S3SSEType *string `yaml:"s3_sse_type,omitempty" json:"s3_sse_type,omitempty" doc:"nocli|description=S3 server-side encryption type. Required to enable server-side encryption overrides for a specific tenant. If not set, the default S3 client settings are used."`
	// +kubebuilder:validation:Optional
	S3SSEKMSKeyID *string `yaml:"s3_sse_kms_key_id,omitempty" json:"s3_sse_kms_key_id,omitempty" doc:"nocli|description=S3 server-side encryption KMS Key ID. Ignored if the SSE type override is not set."`
	// +kubebuilder:validation:Optional
	S3SSEKMSEncryptionContext *string `yaml:"s3_sse_kms_encryption_context,omitempty" json:"s3_sse_kms_encryption_context,omitempty" doc:"nocli|description=S3 server-side encryption KMS encryption context. If unset and the key ID override is set, the encryption context will not be provided to S3. Ignored if the SSE type override is not set."`

	// Alertmanager.
	// Comma-separated list of network CIDRs to block in Alertmanager receiver
	// +kubebuilder:validation:Optional
	AlertmanagerReceiversBlockCIDRNetworks *string `yaml:"alertmanager_receivers_firewall_block_cidr_networks,omitempty" json:"alertmanager_receivers_firewall_block_cidr_networks,omitempty"`
	// +kubebuilder:validation:Optional
	AlertmanagerReceiversBlockPrivateAddresses *bool `yaml:"alertmanager_receivers_firewall_block_private_addresses,omitempty" json:"alertmanager_receivers_firewall_block_private_addresses,omitempty"`

	// +kubebuilder:validation:Optional
	NotificationRateLimit *float64 `yaml:"alertmanager_notification_rate_limit,omitempty" json:"alertmanager_notification_rate_limit,omitempty"`
	// +kubebuilder:validation:Optional
	NotificationRateLimitPerIntegration map[string]*float64 `yaml:"alertmanager_notification_rate_limit_per_integration,omitempty" json:"alertmanager_notification_rate_limit_per_integration,omitempty"`

	// +kubebuilder:validation:Optional
	AlertmanagerMaxConfigSizeBytes *int `yaml:"alertmanager_max_config_size_bytes,omitempty" json:"alertmanager_max_config_size_bytes,omitempty"`
	// +kubebuilder:validation:Optional
	AlertmanagerMaxTemplatesCount *int `yaml:"alertmanager_max_templates_count,omitempty" json:"alertmanager_max_templates_count,omitempty"`
	// +kubebuilder:validation:Optional
	AlertmanagerMaxTemplateSizeBytes *int `yaml:"alertmanager_max_template_size_bytes,omitempty" json:"alertmanager_max_template_size_bytes,omitempty"`
	// +kubebuilder:validation:Optional
	AlertmanagerMaxDispatcherAggregationGroups *int `yaml:"alertmanager_max_dispatcher_aggregation_groups,omitempty" json:"alertmanager_max_dispatcher_aggregation_groups,omitempty"`
	// +kubebuilder:validation:Optional
	AlertmanagerMaxAlertsCount *int `yaml:"alertmanager_max_alerts_count,omitempty" json:"alertmanager_max_alerts_count,omitempty"`
	// +kubebuilder:validation:Optional
	AlertmanagerMaxAlertsSizeBytes *int `yaml:"alertmanager_max_alerts_size_bytes,omitempty" json:"alertmanager_max_alerts_size_bytes,omitempty"`
}

type MimirLimitsInput MimirLimits

type ForwardingRule struct {
	// Ingest defines whether a metric should still be pushed to the Ingesters despite it being forwarded.
	Ingest *bool `yaml:"ingest,omitempty" json:"ingest,omitempty"`
}

const (
	// TenantReadyCondition reports on current status of the Tenant. Ready indicates the tenant has been created and the limits have been applied.
	TenantReadyCondition crhelperTypes.ConditionType = "TenantReady"

	// MimirLimitsValidationErrorReason used when the limits configured for Mimir are invalid.
	MimirLimitsValidationErrorReason = "MimirLimitsValidationError"
)
