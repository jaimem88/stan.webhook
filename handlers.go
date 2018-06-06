package webhook

import (
	"encoding/json"
	"net/http"
)

// HandleHealthcheck handler
func HandleHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{}`))
}

// HandleNotFound handler
func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	writeError(w, &Error{Code: http.StatusNotFound, ErrorMessage: r.RequestURI + " not found"})
}

// HandleStanWebhook expects StanRequest and replies with error or StanResponse
func HandleStanWebhook(w http.ResponseWriter, r *http.Request) {

	var stanRequest StanRequest
	err := json.NewDecoder(r.Body).Decode(&stanRequest)
	if err != nil {
		writeError(w, errBadRequest.msg("HandleStanWebhook "+err.Error()))
		return
	}
	if len(stanRequest.Payload) == 0 {
		writeError(w, errBadRequestMalformedPayload.msg("HandleStanWebhook"))
		return
	}
	responses := []*Response{}
	for _, item := range stanRequest.Payload {
		if item.DRM && item.EpisodeCount > 0 {
			res := &Response{
				Image: item.Image.ShowImage,
				Slug:  item.Slug,
				Title: item.Title,
			}
			responses = append(responses, res)
		}
	}
	res := &StanResponse{Response: responses}
	writeJSON(w, res)
}
