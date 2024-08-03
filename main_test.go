package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"advanced-mock-server/internal/handlers"
	"advanced-mock-server/mockdb"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetResource(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/resource", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GetResource(nil))
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestCreateResource(t *testing.T) {
	resource := mockdb.Resource{ID: "1", Name: "Test Resource"}
	body, _ := json.Marshal(resource)
	req, err := http.NewRequest("POST", "/api/v1/resource", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateResource(nil))
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Resource created successfully")
}

func TestUpdateResource(t *testing.T) {
	// Create the resource first
	resource := mockdb.Resource{ID: "1", Name: "Test Resource"}
	body, _ := json.Marshal(resource)
	req, err := http.NewRequest("POST", "/api/v1/resource", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateResource(nil))
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Now update the resource
	updatedResource := mockdb.Resource{Name: "Updated Resource"}
	updateBody, _ := json.Marshal(updatedResource)
	updateReq, err := http.NewRequest("PUT", "/api/v1/resource/1", bytes.NewBuffer(updateBody))
	if err != nil {
		t.Fatal(err)
	}

	updateRr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/resource/{id}", handlers.UpdateResource(nil)).Methods("PUT")
	router.ServeHTTP(updateRr, updateReq)

	assert.Equal(t, http.StatusOK, updateRr.Code)
	assert.Contains(t, updateRr.Body.String(), "Resource updated successfully")
}

func TestDeleteResource(t *testing.T) {
	// Create the resource first
	resource := mockdb.Resource{ID: "1", Name: "Test Resource"}
	body, _ := json.Marshal(resource)
	req, err := http.NewRequest("POST", "/api/v1/resource", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateResource(nil))
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Now delete the resource
	deleteReq, err := http.NewRequest("DELETE", "/api/v1/resource/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	deleteRr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/resource/{id}", handlers.DeleteResource(nil)).Methods("DELETE")
	router.ServeHTTP(deleteRr, deleteReq)

	assert.Equal(t, http.StatusOK, deleteRr.Code)
	assert.Contains(t, deleteRr.Body.String(), "Resource deleted successfully")
}
