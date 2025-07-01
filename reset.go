package main

import "net/http"

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
	cfg.fileServerHits.Store(0)
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error deleting users", err)
		return
	}
	w.Write([]byte("Succesful reset"))
}
