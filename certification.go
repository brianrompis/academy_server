package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Certification struct {
	ID            string    `json:"ID"`
	Name          string    `json:"Name"`
	IssuedDate    time.Time `json:"IssuedDate"`
	ExpiredDate   time.Time `json:"ExpiredDate"`
	CertificateID string    `json:"CertificateID"`
	ClassroomId   string    `json:"ClassroomId"`
}

func (Certification) TableName() string {
	return "certifications"
}

func allCertification(w http.ResponseWriter, r *http.Request) {
	var certification []Certification
	db.Find(&certification)
	json.NewEncoder(w).Encode(certification)
}

func addCertification(w http.ResponseWriter, r *http.Request) {
	var certification []Certification
	json.NewDecoder(r.Body).Decode(&certification)
	db.Create(&certification)
}

func getCertification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var certification []Certification
	db.First(&certification, "id = ?", params["id"])
	json.NewEncoder(w).Encode(certification)
}

func editCertification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var certification []Certification
	db.First(&certification, "id = ?", params["id"])
	json.NewDecoder(r.Body).Decode(&certification)
	db.Save(&certification)
	json.NewEncoder(w).Encode("Successfully edit the certification.")
}

func removeCertification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var certification []Certification
	db.First(&certification, "id = ?", params["id"])
	db.Delete(&certification)
	json.NewEncoder(w).Encode("The certification is deleted successfully!")
}
