package service

import (
	"net/http"

	"github.com/unrolled/render"
)

func queryStatsHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.Text(w, http.StatusOK, "TBD")
	}
}