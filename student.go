package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type StudentClassroom struct {
	ID                    string    `json:"ID"`
	UserID                string    `json:"UserID"`
	ClassroomPeriodID     string    `json:"ClassroomPeriodID"`
	Status                string    `json:"Status"`
	Grade                 float64   `json:"Grade" gorm:"type:numeric(5,2)"`
	HasCertificate        bool      `json:"HasCertificate"`
	CertificateIssuedDate time.Time `json:"CertificateIssuedDate"`
	StudentSubmission     []StudentSubmission
}

func (StudentClassroom) TableName() string {
	return "student_classroom"
}

type StudentSubmission struct {
	ID                 string    `json:"ID"`
	StudentClassroomID string    `json:"StudentClassroomID"`
	AssignmentID       string    `json:"AssignmentID"`
	DraftGrade         float64   `json:"DraftGrade" gorm:"type:numeric(5,2)"`
	AssignedGrade      float64   `json:"AssignedGrade" gorm:"type:numeric(5,2)"`
	LastUpdate         time.Time `json:"LastUpdate"`
	LastAssignedBy     string    `json:"LastAssignedBy"`
	GradeAssignHistory []GradeAssignHistory
}

func (StudentSubmission) TableName() string {
	return "student_submission"
}

////////////////////////////////////
/// get all student from a class ///
////////////////////////////////////
type ResultStudent struct {
	StudentID      string    `json:"StudentID"`
	Name           string    `json:"Name"`
	Email          string    `json:"Email"`
	Classroom      string    `json:"Classroom"`
	ClassStart     time.Time `json:"ClassStart"`
	Grade          float64   `json:"Grade"`
	HasCertificate bool      `json:"HasCertificate"`
}

func classroomStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Executing Get All Student function")

	params := mux.Vars(r)
	var resultStudent []ResultStudent
	db.Raw(`select "student_classroom"."id" as "StudentID", "user"."full_name" as "Name", "user"."email" as "Email", "classroom"."name" as "Classroom", "classroom_period".start_date as "ClassStart", "student_classroom"."grade" as "Grade", "student_classroom".has_certificate as "HasCertificate" 
	from "classroom"
	inner join "classroom_period" on "classroom".active_period_id = "classroom_period".id
	inner join "student_classroom" on "classroom_period".id = "student_classroom".classroom_period_id
	inner join "user" on "student_classroom".user_id = "user".id
	where "classroom".id = ?
	order by "user".full_name asc`, params["class_id"]).Scan(&resultStudent)

	json.NewEncoder(w).Encode(resultStudent)
}

// func addStudent(w http.ResponseWriter, r *http.Request) {
// 	var student []Student
// 	json.NewDecoder(r.Body).Decode(&student)
// 	db.Create(&student)
// }

// func getStudent(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	var student []Student
// 	db.First(&student, "id = ?", params["id"])
// 	json.NewEncoder(w).Encode(student)
// }

// func editStudent(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	var student []Student
// 	db.First(&student, "id = ?", params["id"])
// 	json.NewDecoder(r.Body).Decode(&student)
// 	db.Save(&student)
// 	json.NewEncoder(w).Encode("Successfully edit the student.")
// }

// func removeStudent(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	var student []Student
// 	db.First(&student, "id = ?", params["id"])
// 	db.Delete(&student)
// 	json.NewEncoder(w).Encode("The student is deleted successfully!")
// }
