package webhook

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func writeError(w http.ResponseWriter, e *Error) {

	log.WithError(e).WithField("error-message", e.message).Error()

	js, _ := json.Marshal(e)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	w.Write(js)
}

func writeJSON(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	js, err := json.Marshal(i)
	if err != nil {
		writeError(w, errInternalServerError.msg("json.Marshal: "+err.Error()))
		return
	}
	w.Write(js)
}
