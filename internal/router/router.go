package router

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/frederickmarvel/supernova/internal/service"
	"github.com/gorilla/mux"
)

func New(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/trend/update", func(w http.ResponseWriter, _ *http.Request) {
		if err := service.UpdateTrends(db); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")

	r.HandleFunc("/trend/check", func(w http.ResponseWriter, _ *http.Request) {
		trends, ts, err := service.GetLatest(db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"timestamp": ts,
			"trend":     trends,
		})
	}).Methods("GET")
	return r
}
