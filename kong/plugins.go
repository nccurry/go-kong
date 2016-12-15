package kong

type PluginsService service

type Plugins struct {
	Data []Plugin `json:"consumer,omitempty"`
	Total int `json:"total,omitempty"`
	Next string `json:"next,omitempty"`
}

type Plugin struct {
	ID string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	CreatedAt int `json:"created_at,omitempty"`
	Enabled bool `json:"enabled,omitempty"`
	APIID string `json:"api_id,omitempty"`
	ConsumerID string `json:"consumer_id,omitempty"`
	Config map[string]interface{} `json:"config,omitempty"`
}