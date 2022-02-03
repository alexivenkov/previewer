package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		r.Header
		return
	})

	server := &http.Server{Addr: ":8080", Handler: mux}
	server.ListenAndServe()
}
