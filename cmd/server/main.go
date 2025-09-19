package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sakshamg567/cachy/internal/coordinator"
)

func main() {
	addresses := []string{"localhost:50051", "localhost:50052", "localhost:50053"}

	cd := coordinator.NewCoordinator(addresses)

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		val, err := cd.Get(r.Context(), key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"value": val})
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ok := cd.Set(r.Context(), body.Key, body.Value)
		if !ok {
			http.Error(w, "failed to set", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/add-node", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Address string `json:"address"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		log.Printf("adding new node: %s", body.Address)
		cd.AddNode(body.Address)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Node addition process started"))
	})

	log.Fatal(http.ListenAndServe(":6969", nil))
}
