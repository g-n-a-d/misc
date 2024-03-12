package crud

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"modules/crud/item"
)

// HandleCreate handles HTTP POST requests to create a new item
func HandleCreate(w http.ResponseWriter, r *http.Request) {
	var newItem item.Item
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, "Decode failed", http.StatusBadRequest)
		return
	}

	err = Insert(newItem)
	if err != nil {
		http.Error(w, "Insert failed", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// HandleRead handles HTTP GET requests to retrieve an item by ID
func HandleRead(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	id_, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Query failed", http.StatusBadRequest)
		return
	}

	item, err := Query(id_)
	if err != nil {
		http.Error(w, "Query failed", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

// HandleUpdate handles HTTP PUT requests to update an existing item
func HandleUpdate(w http.ResponseWriter, r *http.Request) {
	var updatedItem item.Item
	err := json.NewDecoder(r.Body).Decode(&updatedItem)
	if err != nil {
		http.Error(w, "Decode failed", http.StatusBadRequest)
		return
	}

	err = Update(updatedItem)
	if err != nil {
		http.Error(w, "Update failed", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleDelete handles HTTP DELETE requests to delete an item by ID
func HandleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	id_, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Query failed", http.StatusBadRequest)
		return
	}

	err = Delete(id_)
	if err != nil {
		http.Error(w, "Delete failed", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Start_app() {
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		HandleCreate(w, r)
	})

	http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
		HandleRead(w, r)
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		HandleUpdate(w, r)
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		HandleDelete(w, r)
	})

	port := 8080
	fmt.Printf("CRUD server is listening on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil))
}