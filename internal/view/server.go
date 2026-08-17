package view

import (
	"encoding/json"
	"net/http"

	"mon-pi/internal/collector"
)

func jsonOK(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

func SetupRoutes(col *collector.Collector) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/snapshot", func(w http.ResponseWriter, r *http.Request) {
		snap := col.Latest()
		if snap == nil {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			return
		}
		jsonOK(w, snap)
	})

	mux.HandleFunc("/api/history", func(w http.ResponseWriter, r *http.Request) {
		jsonOK(w, map[string]any{
			"cpu": col.CPUHistory(),
			"mem": col.MemHistory(),
		})
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(dashboardHTML))
	})

	return mux
}
