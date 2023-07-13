package v1alpha1

import (
	prom_v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	prom_config "github.com/prometheus/common/config"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RemoteWriteSpec struct {
	URL string `yaml:"url" json:"url"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	RemoteTimeout *metav1.Duration `yaml:"remote_timeout,omitempty" json:"remote_timeout,omitempty"`
	// +kubebuilder:validation:Optional
	Headers map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	// +kubebuilder:validation:Optional
	WriteRelabelConfigs []*RelabelConfig `yaml:"write_relabel_configs,omitempty" json:"write_relabel_configs,omitempty"`
	// +kubebuilder:validation:Optional
	Name *string `yaml:"name,omitempty" json:"name,omitempty"`
	// +kubebuilder:validation:Optional
	SendExemplars *bool `yaml:"send_exemplars,omitempty" json:"send_exemplars,omitempty"`
	// +kubebuilder:validation:Optional
	SendNativeHistograms *bool `yaml:"send_native_histograms,omitempty" json:"send_native_histograms,omitempty"`

	// We cannot do proper Go type embedding below as the parser will then parse
	// values arbitrarily into the overflow maps of further-down types.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XPreserveUnknownFields
	// +kubebuilder:pruning:PreserveUnknownFields
	HTTPClientConfig *HTTPClientConfig `yaml:",inline" json:",inline"`
	// +kubebuilder:validation:Optional
	QueueConfig *QueueConfig `yaml:"queue_config,omitempty" json:"queue_config,omitempty"`
	// +kubebuilder:validation:Optional
	MetadataConfig *MetadataConfig `yaml:"metadata_config,omitempty" json:"metadata_config,omitempty"`
	// +kubebuilder:validation:Optional
	SigV4Config *SigV4Config `yaml:"sigv4,omitempty" json:"sigv4,omitempty"`
}

// MetadataConfig is the configuration for sending metadata to remote
// storage.
type MetadataConfig struct {
	// Send controls whether we send metric metadata to remote storage.
	// +kubebuilder:validation:Optional
	Send *bool `yaml:"send,omitempty" json:"send,omitempty"`
	// SendInterval controls how frequently we send metric metadata.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	SendInterval *metav1.Duration `yaml:"send_interval,omitempty" json:"send_interval,omitempty"`
	// Maximum number of samples per send.
	// +kubebuilder:validation:Optional
	MaxSamplesPerSend *int `yaml:"max_samples_per_send,omitempty" json:"max_samples_per_send,omitempty"`
}

// SigV4Config is the configuration for signing remote write requests with
// AWS's SigV4 verification process. Empty values will be retrieved using the
// AWS default credentials chain.
type SigV4Config struct {
	// +kubebuilder:validation:Optional
	Region *string `yaml:"region,omitempty" json:"region,omitempty"`
	// +kubebuilder:validation:Optional
	AccessKey *string `yaml:"access_key,omitempty" json:"access_key,omitempty"`
	// +kubebuilder:validation:Optional
	SecretKey *prom_config.Secret `yaml:"secret_key,omitempty" json:"secret_key,omitempty"`
	// +kubebuilder:validation:Optional
	Profile *string `yaml:"profile,omitempty" json:"profile,omitempty"`
	// +kubebuilder:validation:Optional
	RoleARN *string `yaml:"role_arn,omitempty" json:"role_arn,omitempty"`
}

type QueueConfig struct {
	// Number of samples to buffer per shard before we block. Defaults to
	// MaxSamplesPerSend.
	// +kubebuilder:validation:Optional
	Capacity *int `yaml:"capacity,omitempty" json:"capacity,omitempty"`

	// Max number of shards, i.e. amount of concurrency.
	// +kubebuilder:validation:Optional
	MaxShards *int `yaml:"max_shards,omitempty" json:"max_shards,omitempty"`

	// Min number of shards, i.e. amount of concurrency.
	// +kubebuilder:validation:Optional
	MinShards *int `yaml:"min_shards,omitempty" json:"min_shards,omitempty"`

	// Maximum number of samples per send.
	// +kubebuilder:validation:Optional
	MaxSamplesPerSend *int `yaml:"max_samples_per_send,omitempty" json:"max_samples_per_send,omitempty"`

	// Maximum time sample will wait in buffer.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	BatchSendDeadline *metav1.Duration `yaml:"batch_send_deadline,omitempty" json:"batch_send_deadline,omitempty"`

	// On recoverable errors, backoff exponentially.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MinBackoff *metav1.Duration `yaml:"min_backoff,omitempty" json:"min_backoff,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	MaxBackoff *metav1.Duration `yaml:"max_backoff,omitempty" json:"max_backoff,omitempty"`
	// +kubebuilder:validation:Optional
	RetryOnRateLimit *bool `yaml:"retry_on_http_429,omitempty" json:"retry_on_http_429,omitempty"`
}

// HTTPClientConfig configures an HTTP client.
type HTTPClientConfig struct {
	// The HTTP basic authentication credentials for the targets.
	// +kubebuilder:validation:Optional
	BasicAuth *prom_config.BasicAuth `yaml:"basic_auth,omitempty" json:"basic_auth,omitempty"`
	// The HTTP authorization credentials for the targets.
	// +kubebuilder:validation:Optional
	Authorization *prom_config.Authorization `yaml:"authorization,omitempty" json:"authorization,omitempty"`
	// The OAuth2 client credentials used to fetch a token for the targets.
	// +kubebuilder:validation:Optional
	OAuth2 *OAuth2 `yaml:"oauth2,omitempty" json:"oauth2,omitempty"`
	// TLSConfig to use to connect to the targets.
	// +kubebuilder:validation:Optional
	TLSConfig *prom_config.TLSConfig `yaml:"tls_config,omitempty" json:"tls_config,omitempty"`
	// FollowRedirects specifies whether the client should follow HTTP 3xx redirects.
	// The omitempty flag is not set, because it would be hidden from the
	// marshalled configuration when set to false.
	// +kubebuilder:validation:Optional
	FollowRedirects *bool `yaml:"follow_redirects" json:"follow_redirects"`
	// EnableHTTP2 specifies whether the client should configure HTTP2.
	// The omitempty flag is not set, because it would be hidden from the
	// marshalled configuration when set to false.
	// +kubebuilder:validation:Optional
	EnableHTTP2 *bool `yaml:"enable_http2" json:"enable_http2"`
	// Proxy configuration.
	// +kubebuilder:validation:Optional
	ProxyConfig `yaml:",inline" json:",inline"`
}

// OAuth2 is the oauth2 client configuration.
type OAuth2 struct {
	ClientID string `yaml:"client_id" json:"client_id"`
	// +kubebuilder:validation:Optional
	ClientSecret *prom_config.Secret `yaml:"client_secret" json:"client_secret"`
	// +kubebuilder:validation:Optional
	ClientSecretFile *string `yaml:"client_secret_file" json:"client_secret_file"`
	// +kubebuilder:validation:Optional
	Scopes []*string `yaml:"scopes,omitempty" json:"scopes,omitempty"`
	// +kubebuilder:validation:Optional
	TokenURL *string `yaml:"token_url" json:"token_url"`
	// +kubebuilder:validation:Optional
	EndpointParams map[string]string `yaml:"endpoint_params,omitempty" json:"endpoint_params,omitempty"`
	// +kubebuilder:validation:Optional
	TLSConfig *prom_config.TLSConfig `yaml:"tls_config,omitempty" json:"tls_config,omitempty"`
	// +kubebuilder:validation:Optional
	ProxyConfig `yaml:",inline" json:",inline"`
}

type ProxyConfig struct {
	// HTTP proxy server to use to connect to the targets.
	// +kubebuilder:validation:Optional
	ProxyURL *string `yaml:"proxy_url,omitempty" json:"proxy_url,omitempty"`
	// NoProxy contains addresses that should not use a proxy.
	// +kubebuilder:validation:Optional
	NoProxy *string `yaml:"no_proxy,omitempty" json:"no_proxy,omitempty"`
	// ProxyFromEnvironment makes use of net/http ProxyFromEnvironment function
	// to determine proxies.
	// +kubebuilder:validation:Optional
	ProxyFromEnvironment *bool `yaml:"proxy_from_environment,omitempty" json:"proxy_from_environment,omitempty"`
	// ProxyConnectHeader optionally specifies headers to send to
	// proxies during CONNECT requests. Assume that at least _some_ of
	// these headers are going to contain secrets and use Secret as the
	// value type instead of string.
	// +kubebuilder:validation:Optional
	ProxyConnectHeader *prom_config.Header `yaml:"proxy_connect_header,omitempty" json:"proxy_connect_header,omitempty"`
}

type RelabelConfig struct {
	// A list of labels from which values are taken and concatenated
	// with the configured separator in order.
	// +kubebuilder:validation:Optional
	SourceLabels []*prom_v1.LabelName `yaml:"source_labels,omitempty" json:"source_labels,omitempty"`
	// Separator is the string between concatenated values from the source labels.
	// +kubebuilder:validation:Optional
	Separator *string `yaml:"separator,omitempty" json:"separator,omitempty"`
	// Regex against which the concatenation is matched.
	// +kubebuilder:validation:Optional
	Regex *string `yaml:"regex,omitempty" json:"regex,omitempty"`
	// Modulus to take of the hash of concatenated values from the source labels.
	// +kubebuilder:validation:Optional
	Modulus *uint64 `yaml:"modulus,omitempty" json:"modulus,omitempty"`
	// TargetLabel is the label to which the resulting string is written in a replacement.
	// Regexp interpolation is allowed for the replace action.
	// +kubebuilder:validation:Optional
	TargetLabel *string `yaml:"target_label,omitempty" json:"target_label,omitempty"`
	// Replacement is the regex replacement pattern to be used.
	// +kubebuilder:validation:Optional
	Replacement *string `yaml:"replacement,omitempty" json:"replacement,omitempty"`
	// Action is the action to be performed for the relabeling.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=replace;Replace;keep;Keep;drop;Drop;hashmod;HashMod;labelmap;LabelMap;labeldrop;LabelDrop;labelkeep;LabelKeep;lowercase;Lowercase;uppercase;Uppercase;keepequal;KeepEqual;dropequal;DropEqual
	// +kubebuilder:default=replace
	Action *string `yaml:"action,omitempty" json:"action,omitempty"`
}
