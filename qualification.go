package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Qualification struct {
	ID                     string `json:"ID"`
	Name                   string `json:"Name"`
	Description            string `json:"Description"`
	CertificateTemplateID  string `json:"CertificateTemplateID"`
	QualificationClassroom []QualificationClassroom
	JobQualification       []JobQualification
}

func (Qualification) TableName() string {
	return "qualification"
}

type QualificationClassroom struct {
	ID              string `json:"ID"`
	QualificationID string `json:"QualificationID"`
	ClassroomID     string `json:"ClassroomID"`
}

func (QualificationClassroom) TableName() string {
	return "qualification_classroom"
}

func allQualification(w http.ResponseWriter, r *http.Request) {
	var qualification []Qualification
	db.Find(&qualification)
	json.NewEncoder(w).Encode(qualification)
}

type StudentQualification struct {
	UserID        string `json:"UserID"`
	Name          string `json:"Name"`
	Email         string `json:"Email"`
	ClassroomID   string `json:"ClassroomID"`
	Classroom     string `json:"Classroom"`
	Qualification string `json:"Qualification"`
}

///////////////////////////////////
//// get student Qualification ////
///////////////////////////////////
func getStudentQualification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var studentQualification []StudentQualification
	db.Raw(`select "user".id as "user_id", "user".full_name as "name", "user".email, "classroom".id as "classroom_id", "classroom".name as "classroom", "qualification".id as "qualification_id", "qualification".name as "qualification"  
	from "user" 
	inner join "student_classroom" on "user".id = "student_classroom".user_id
	and "student_classroom".has_certificate  = true
	inner join "classroom_period" on "student_classroom".classroom_period_id = "classroom_period".id  
	inner join "classroom" on "classroom_period".classroom_id = "classroom".id 
	inner join "qualification_classroom" on "qualification_classroom".classroom_id = "classroom".id 
	right join "qualification" on "qualification_classroom".qualification_id = "qualification".id
	where "user".id = ?
	order by "classroom" asc`, params["user_id"]).Scan(&studentQualification)

	json.NewEncoder(w).Encode(studentQualification)
}

func addQualification(w http.ResponseWriter, r *http.Request) {
	var qualification []Qualification
	json.NewDecoder(r.Body).Decode(&qualification)
	db.Create(&qualification)
}

func getQualification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var qualification []Qualification
	db.First(&qualification, "id = ?", params["id"])
	json.NewEncoder(w).Encode(qualification)
}

func editQualification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var qualification []Qualification
	db.First(&qualification, "id = ?", params["id"])
	json.NewDecoder(r.Body).Decode(&qualification)
	db.Save(&qualification)
	json.NewEncoder(w).Encode("Successfully edit the qualification.")
}

func removeQualification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var qualification []Qualification
	db.First(&qualification, "id = ?", params["id"])
	db.Delete(&qualification)
	json.NewEncoder(w).Encode("The qualification is deleted successfully!")
}
