package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Student struct {
	ID                     string `json:"ID"`
	UserID                 string `json:"UserID"`
	ClassroomApplied       string `json:"ClassroomApplied"`
	ClassroomMember        string `json:"ClassroomMember"`
	ClassroomCompleted     string `json:"ClassroomCompleted"`
	CertificationApplied   string `json:"CertificationApplied"`
	CertificationMember    string `json:"CertificationMember"`
	CertificationCompleted string `json:"CertificationCompleted"`
	CertificateID          string `json:"CertificateID"`
}

func (Student) TableName() string {
	return "students"
}

func allStudent(w http.ResponseWriter, r *http.Request) {
	var student []Student
	db.Find(&student)
	json.NewEncoder(w).Encode(student)
}

func addStudent(w http.ResponseWriter, r *http.Request) {
	var student []Student
	json.NewDecoder(r.Body).Decode(&student)
	db.Create(&student)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var student []Student
	db.First(&student, "id = ?", params["id"])
	json.NewEncoder(w).Encode(student)
}

func editStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var student []Student
	db.First(&student, "id = ?", params["id"])
	json.NewDecoder(r.Body).Decode(&student)
	db.Save(&student)
	json.NewEncoder(w).Encode("Successfully edit the student.")
}

func removeStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var student []Student
	db.First(&student, "id = ?", params["id"])
	db.Delete(&student)
	json.NewEncoder(w).Encode("The student is deleted successfully!")
}
