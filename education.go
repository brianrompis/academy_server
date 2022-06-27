package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

type EducationHistory struct {
	ID             string          `json:"ID"`
	UserID         string          `json:"UserID"`
	EducationLevel string          `json:"EducationLevel"`
	SchoolName     string          `json:"SchoolName"`
	StartYear      time.Time       `json:"StartYear"`
	EndYear        time.Time       `json:"EndYear"`
	GPA            decimal.Decimal `json:"GPA"`
}

func (EducationHistory) TableName() string {
	return "education_history"
}

func allEducationHistory(w http.ResponseWriter, r *http.Request) {
	var education []EducationHistory
	db.Find(&education)
	json.NewEncoder(w).Encode(education)
}

func addEducationHistory(w http.ResponseWriter, r *http.Request) {
	var education []EducationHistory
	json.NewDecoder(r.Body).Decode(&education)
	db.Create(&education)
	json.NewEncoder(w).Encode("Successfully add the education history.")
}

func getEducationHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var education []EducationHistory
	db.First(&education, "id = ?", params["id"])
	json.NewEncoder(w).Encode(education)
}

func editEducationHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var education []EducationHistory
	db.First(&education, "id = ?", params["id"])
	json.NewDecoder(r.Body).Decode(&education)
	db.Save(&education)
	json.NewEncoder(w).Encode("Successfully edit the education history.")
}

func removeEducationHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var education []EducationHistory
	db.First(&education, "id = ?", params["id"])
	db.Delete(&education)
	json.NewEncoder(w).Encode("The education history is deleted successfully!")
}

func getUserEducationHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var education []EducationHistory
	db.Where("user_id = ?", params["user_id"]).Find(&education)
	json.NewEncoder(w).Encode(education)
}

func removeUserEducationHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var education []EducationHistory
	db.Where("user_id = ?", params["user_id"]).Delete(&education)
	json.NewEncoder(w).Encode("The entire education histories for the user are deleted successfully!")
}
