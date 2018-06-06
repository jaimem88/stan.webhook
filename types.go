package webhook

// Service ..
type Service struct {
	environment string
	hostname    string
}

// StanRequest describes the payload received in the webhook endpoint
type StanRequest struct{}

// StanResponse describes the expected response to be sent to Stan
type StanResponse struct{}
