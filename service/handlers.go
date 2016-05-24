package service

import (
	"net/http"

	"github.com/unrolled/render"
)

func homeHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.Text(w, http.StatusOK, "Drones Sample - Event Processor see http://github.com/cloudnativego/drones-events")
	}
}
