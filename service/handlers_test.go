package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var (
	formatter = render.New(render.Options{
		IndentJSON: true,
	})
)

func MakeTestServer() *negroni.Negroni {
	server := negroni.New()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	server.UseHandler(mx)
	return server
}

// TestStatsHandler just ensures the stats endpoint is reporting the package-wide variable
func TestStatsHandler(t *testing.T) {
	Stats.AlertEventCount = 9
	Stats.PositionEventCount = 10
	Stats.TelemetryEventCount = 11

	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	server := MakeTestServer()
	recorder = httptest.NewRecorder()

	request, _ = http.NewRequest("GET", "/api/stats", nil)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected stats endpoint to return HTTP 200, got %d", recorder.Code)
	}

	var statsResponse ProcessingStats
	payload := recorder.Body.Bytes()
	err := json.Unmarshal(payload, &statsResponse)
	if err != nil {
		t.Errorf("Got an error unmarshaling stats response, %s", err)
	}
	if statsResponse.AlertEventCount != 9 ||
		statsResponse.PositionEventCount != 10 ||
		statsResponse.TelemetryEventCount != 11 {
		t.Errorf("Got incorrect stats reply from endpoint, %+v", statsResponse)
	}
}
