package healthz

// Error the structure of the Error object
type Error struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

// Service struct reprecenting a healthz provider and it's status
type Service struct {
	Name         string `json:"name"`
	Healthy      bool   `json:"healthy"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// Response type, we return a json object with {healthy:bool, errors:[]}
type HealthzResponse struct {
	Services []Service `json:"services,omitempty"`
	Healthy  bool      `json:"healthy"`
}
