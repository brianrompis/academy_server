package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type EmploymentHistory struct {
	ID        string    `json:"ID"`
	StudentID string    `json:"StudentID"`
	Email     string    `json:"Email"`
	HotelID   string    `json:"HotelID"`
	StartDate time.Time `json:"StartDate"`
	EndDate   time.Time `json:"EndDate"`
}

func (EmploymentHistory) TableName() string {
	return "employment_histories"
}

func allEmploymentHistory(w http.ResponseWriter, r *http.Request) {
	var employment_history []EmploymentHistory
	db.Find(&employment_history)
	json.NewEncoder(w).Encode(employment_history)
}

func addEmploymentHistory(w http.ResponseWriter, r *http.Request) {
	var employment_history []EmploymentHistory
	json.NewDecoder(r.Body).Decode(&employment_history)
	db.Create(&employment_history)
}

func getEmploymentHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var employment_history []EmploymentHistory
	db.First(&employment_history, "id = ?", params["id"])
	json.NewEncoder(w).Encode(employment_history)
}

func editEmploymentHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var employment_history []EmploymentHistory
	db.First(&employment_history, "id = ?", params["id"])
	json.NewDecoder(r.Body).Decode(&employment_history)
	db.Save(&employment_history)
	json.NewEncoder(w).Encode("Successfully edit the employment_history.")
}

func removeEmploymentHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var employment_history []EmploymentHistory
	db.First(&employment_history, "id = ?", params["id"])
	db.Delete(&employment_history)
	json.NewEncoder(w).Encode("The employment_history is deleted successfully!")
}
