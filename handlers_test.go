package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	jsonBadBody = bytes.NewBuffer([]byte(`{"bad":"request"}`))
	jsonBody    = bytes.NewBuffer(testPayload)
)

func TestService_HandleStanWebhook(t *testing.T) {

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	type want struct {
		code         int
		errorMessage string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"handle Webhook - 400 on missing Payload key in request body",
			args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "http://localhost:8080/v1xstan-webhook", jsonBadBody),
			},
			want{
				code:         400,
				errorMessage: errBadRequestMalformedPayload.ErrorMessage,
			},
		},
		{
			"handle Webhook - 400 on empty body",
			args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "http://localhost:8080/v1xstan-webhook", nil),
			},
			want{
				code:         400,
				errorMessage: errBadRequest.ErrorMessage,
			},
		},
		{
			"handle Webhook - 400 on empty body",
			args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "http://localhost:8080/v1xstan-webhook", jsonBody),
			},
			want{
				code:         200,
				errorMessage: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleStanWebhook(tt.args.w, tt.args.r)
			resp := tt.args.w.Result()
			var e Error
			json.NewDecoder(resp.Body).Decode(&e)

			if got := resp.StatusCode; got != tt.want.code {
				t.Errorf("HandleStanWebhook got = %v, want %v", got, tt.want.code)
			}
			if got := e.ErrorMessage; got != tt.want.errorMessage {
				t.Errorf("HandleStanWebhook got = %+v, want %+v", got, tt.want.errorMessage)
			}
		})
	}
}
