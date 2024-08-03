package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"advanced-mock-server/mockdb"
)

var validate = validator.New()

type Response struct {
	Message string `json:"message"`
}

func GetResource(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resources := mockdb.GetAllResources(db)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resources)
	}
}

func CreateResource(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resource mockdb.Resource
		if err := json.NewDecoder(r.Body).Decode(&resource); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := validate.Struct(resource); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mockdb.AddResource(db, resource)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Message: "Resource created successfully"})
	}
}

func UpdateResource(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var resource mockdb.Resource
		if err := json.NewDecoder(r.Body).Decode(&resource); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := validate.Struct(resource); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := mockdb.UpdateResource(db, id, resource); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Message: "Resource updated successfully"})
	}
}

func DeleteResource(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		if err := mockdb.DeleteResource(db, id); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Message: "Resource deleted successfully"})
	}
}

func ResetDatabaseHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mockdb.ResetDatabase(db)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Message: "Database reset successfully"})
	}
}
