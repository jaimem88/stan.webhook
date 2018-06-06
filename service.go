package webhook

// NewService initialises webhook service
func NewService(environmet, hostname string) *Service {
	return &Service{
		environment: environmet,
		hostname:    hostname,
	}
}
