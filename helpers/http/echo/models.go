package helperecho

type Result struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type Err struct {
	Errors []string `json:"errors"`
}
