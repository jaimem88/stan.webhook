package webhook

import (
	"encoding/json"
	"net/http"
)

// Healthcheck handler
func (s *Service) Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{}`))
}

// NotFound handler
func (s *Service) NotFound(w http.ResponseWriter, r *http.Request) {
	writeError(w, &Error{Code: http.StatusNotFound, ErrorMessage: r.RequestURI + " not found"})
}

// HandleStanWebhook expects StanRequest and replies with error or StanResponse
func (s *Service) HandleStanWebhook(w http.ResponseWriter, r *http.Request) {

	var stanRequest *StanRequest
	err := json.NewDecoder(r.Body).Decode(&stanRequest)
	if err != nil {
		writeError(w, errBadRequest.msg("decode HandleCreateLead "+err.Error()))
		return
	}
	res := StanResponse{}
	writeJSON(w, res)
}
