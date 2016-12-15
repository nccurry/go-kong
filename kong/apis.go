package kong


// ApisService handles communication with the apis resource of Kong
type ApisService service

// Apis represents a list of Kong Apis
type Apis struct {
	Data []Api `json:"data,omitempty"`
	Total int `json:"total,omitempty"`
	Next string `json:"next,omitempty"`
}

type Api struct {

}