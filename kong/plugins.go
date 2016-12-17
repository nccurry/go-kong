package kong

import (
	"github.com/fatih/structs"
	"net/http"
	"reflect"
	"strings"
)

type PluginsService service

type Plugins struct {
	Data  []Plugin `json:"consumer,omitempty"`
	Total int      `json:"total,omitempty"`
	Next  string   `json:"next,omitempty"`
}

// Plugin defines a generic Kong plugin. Due to the fact that
// Config has multiple structures depending on the plugin
// it is defined as a map[string]interface{} here.
//
// One of the below Plugin definitions can be used if you know
// the structure of Config. Each of the below Plugin structs
// expose a ToGeneric method that will allow you to convert
// them to the generic Plugin type if necessary.
type Plugin struct {
	ID         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	CreatedAt  int         `json:"created_at,omitempty"`
	Enabled    bool        `json:"enabled,omitempty"`
	APIID      string      `json:"api_id,omitempty"`
	ConsumerID string      `json:"consumer_id,omitempty"`
	Config     interface{} `json:"config,omitempty"`
}

// This method allows for a variety of plugin creation behaviors.
// If you want to add a plugin for:
// Every API and Consumer - Don't set api_id and consumer_id in plugin
// Every API and a specific Consumer - Only set consumer_id in plugin
// For every Consumer and a specific API - Only set api_id in plugin
// For a specific Consumer and API - Set both api_id and consumer_id in plugin
func (s *PluginsService) Post(plugin *Plugin) (*http.Response, error) {
	req, err := s.client.NewRequest("POST", "plugins", plugin)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

func (s *PluginsService) GetAll() (*Plugins, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "plugins", nil)
	if err != nil {
		return nil, nil, err
	}

	plugins := new(Plugins)
	resp, err := s.client.Do(req, plugins)
	if err != nil {
		return nil, resp, err
	}

	return plugins, resp, err
}

// https://getkong.org/plugins/acl/
type ACLPlugin struct {
	Plugin
	Config *ACLConfig `json:"config,omitempty"`
}

// ToPlugin converts ACLPlugin to the more generic Plugin.
// This assumes that Config does not contain nested fields
func (c *ACLPlugin) ToPlugin() *Plugin {
	s := structs.New(c.Config)
	fieldNames := s.Names()

	config := make(map[string]interface{})

	// Iterate over Config fields and create a map with keys based on the json tags
	for _, v := range fieldNames {
		tags := strings.Split(s.Field(v).Tag("json"), ",")

		// isZero checks whether the value of the field v is 'zero'
		// If it is, we do not need to add it to our map[string]interface{}
		// Doing so causes errors later when we try to marshal to JSON
		if ok := isZero(reflect.ValueOf(s.Field(v).Value())); !ok {
			config[tags[0]] = s.Field(v).Value()
		}

	}

	plugin := &c.Plugin
	plugin.Config = config
	return plugin
}

type ACLConfig struct {
	Whitelist []string `json:"whitelist,omitempty"`
	Blacklist []string `json:"blacklist,omitempty"`
}

// https://getkong.org/plugins/request-size-limiting/
type RequestSizeLimitingPlugin struct {
	Plugin
	Config RequestSizeLimitingConfig `json:"config,omitempty"`
}

type RequestSizeLimitingConfig struct {
	AllowedPayloadSize int `json:"allowed_payload_size,omitempty"`
}

// https://getkong.org/plugins/correlation-id/
type CorrelationIDPlugin struct {
	Plugin
	Config CorrelationIDConfig `json:"config,omitempty"`
}

type CorrelationIDConfig struct {
	HeaderName     string `json:"header_name,omitempty"`
	Generator      string `json:"generator,omitempty"`
	EchoDownstream bool   `json:"echo_downstream,omitempty"`
}

// https://getkong.org/plugins/rate-limiting/
type RateLimitingPlugin struct {
	Plugin
	Config RateLimitingConfig `json:"config,omitempty"`
}

type RateLimitingConfig struct {
	Second        int    `json:"second,omitempty"`
	Minute        int    `json:"minute,omitempty"`
	Hour          int    `json:"hour,omitempty"`
	Day           int    `json:"day, omitempty"`
	Month         int    `json:"month,omitempty"`
	Year          int    `json:"year, omitempty"`
	LimitBy       string `json:"limit_by,omitempty"`
	Policy        string `json:"policy,omitempty"`
	FaultTolerant bool   `json:"fault_tolerant,omitempty"`
	RedisHost     string `json:"redis_host,omitempty"`
	RedisPort     string `json:"redis_port,omitempty"`
	RedisPassword string `json:"redis_password,omitempty"`
	RedisTimeout  int    `json:"redis_timeout,omitempty"`
}

// https://getkong.org/plugins/jwt/
type JWTPlugin struct {
	Plugin
	Config JWTConfig `json:"config,omitempty"`
}

type JWTConfig struct {
	URIParamNames  []string `json:"uri_param_names,omitempty"`
	ClaimsToVerify []string `json:"claims_to_verify,omitempty"`
	KeyClaimName   string   `json:"key_claim_name,omitempty"`
	SecretIsBase64 bool     `json:"secret_is_base64,omitempty"`
}

// https://getkong.org/plugins/file-log/
type FileLogPlugin struct {
	Plugin
	Config FileLogConfig `json:"config,omitempty"`
}

type FileLogConfig struct {
	Path string `json:"path,omitempty"`
}

// https://getkong.org/plugins/key-authentication/
type KeyAuthenticationPlugin struct {
	Plugin
	Config KeyAuthenticationConfig `json:"config,omitempty"`
}

type KeyAuthenticationConfig struct {
	KeyNames        []string `json:"key_names,omitempty"`
	HideCredentials bool     `json:"hide_credentials,omitempty"`
}
