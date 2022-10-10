package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Certificate struct {
	ID          string `json:"ID"`
	Template    string `json:"Template"`
	ClassroomID string `json:"ClassroomID"`
}

func (Certificate) TableName() string {
	return "certificates"
}

func allCertificate(w http.ResponseWriter, r *http.Request) {
	var certificate []Certificate
	db.Find(&certificate)
	json.NewEncoder(w).Encode(certificate)
}

func addCertificate(w http.ResponseWriter, r *http.Request) {
	var certificate []Certificate
	json.NewDecoder(r.Body).Decode(&certificate)
	db.Create(&certificate)
}

func getCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var certificate []Certificate
	db.First(&certificate, "id = ?", params["id"])
	json.NewEncoder(w).Encode(certificate)
}

func editCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var certificate []Certificate
	db.First(&certificate, "id = ?", params["id"])
	json.NewDecoder(r.Body).Decode(&certificate)
	db.Save(&certificate)
	json.NewEncoder(w).Encode("Successfully edit the certificate.")
}

func removeCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var certificate []Certificate
	db.First(&certificate, "id = ?", params["id"])
	db.Delete(&certificate)
	json.NewEncoder(w).Encode("The certificate is deleted successfully!")
}

func getStudentCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var certificates []Certificate
	db.Where("student_id = ?", params["student_id"]).Find(&certificates)
	json.NewEncoder(w).Encode(certificates)
}

func removeStudentCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var certificates []Certificate
	db.Where("student_id = ?", params["student_id"]).Delete(&certificates)
	json.NewEncoder(w).Encode("The entire certificates for the student are deleted successfully!")
}
