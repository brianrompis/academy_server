package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type EmploymentHistory struct {
	ID        string    `json:"ID"`
	UserID	string	`json:"UserID"`
	HotelID   string    `json:"HotelID"`
	Position	string		`json:"Position"`
	CompanyName	string	`json:"CompanyName"`
	StartDate time.Time `json:"StartDate"`
	EndDate   time.Time `json:"EndDate"`
	City	string	`json:"City"`
	Description	string	`json:"Description"`
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
	json.NewEncoder(w).Encode("Successfully add an employment history.")
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

func getUserEmploymentHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var employment_history []EmploymentHistory
	db.Where("user_id = ?", params["user_id"]).Find(&employment_history)
	json.NewEncoder(w).Encode(employment_history)
}

func removeUserEmploymentHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var employment_history []EmploymentHistory
	db.Where("user_id = ?", params["user_id"]).Delete(&employment_history)
	json.NewEncoder(w).Encode("The entire employment histories for the user are deleted successfully!")
}

type DeleteID struct {
	ID	string `json:"ID"`
}

func deleteMultipleEmploymentHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var employment_history []EmploymentHistory
	var deleteID []DeleteID
	a := []string{}
	json.NewDecoder(r.Body).Decode(&deleteID)
	for _, s := range deleteID {
		a = append(a, s.ID)
	}
	db.Delete(&employment_history, a)
	json.NewEncoder(w).Encode("The employment_histories are deleted successfully!")
}

func editMultipleEmploymentHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var employment_history []EmploymentHistory
	json.NewDecoder(r.Body).Decode(&employment_history)
	for _, s := range employment_history {
		db.First(&s, "id = ?", s.ID)
		db.Save(&s)
	}
	json.NewEncoder(w).Encode("The employment_histories are edited successfully!")
}


