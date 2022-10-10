package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TeacherClassroom struct {
	ID                 string               `json:"ID"`
	UserID             string               `json:"UserID"`
	ClassroomID        string               `json:"ClassroomID"`
	StudentSubmission  []StudentSubmission  `gorm:"foreignKey:LastAssignedBy"`
	GradeAssignHistory []GradeAssignHistory `gorm:"foreignKey:TeacherID"`
}

func (TeacherClassroom) TableName() string {
	return "teacher_classroom"
}

type GradeAssignHistory struct {
	ID                  string    `json:"ID"`
	StudentSubmissionID string    `json:"StudentSubmissionID"`
	TeacherID           string    `json:"TeacherID"`
	AssignedGrade       float64   `json:"AssignedGrade" gorm:"type:numeric(5,2)"`
	ChangeDate          time.Time `json:"ChangeDate"`
}

func (GradeAssignHistory) TableName() string {
	return "grade_assign_history"
}

///////////////////////////////////
//////// Get all teacher //////////
///////////////////////////////////
// select all user with is_teacher true
func allTeacher(w http.ResponseWriter, r *http.Request) {
	var users []User
	db.Where("is_teacher = ?", true).Find(&users)
	json.NewEncoder(w).Encode(users)
}

///////////////////////////////////
////// Register new teacher ///////
///////////////////////////////////
// set is_teacher to true
func registerNewTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	db.Model(&User{}).Where("id = ?", params["id"]).Update("is_teacher", true)
	json.NewEncoder(w).Encode("Successfully register a teacher.")
}

///////////////////////////////////
//// Add a teacher to a class /////
///////////////////////////////////
type NewTeacherClass struct {
	UserID      string `json:"UserID"`
	ClassroomID string `json:"ClassroomID"`
}

func addTeacher(w http.ResponseWriter, r *http.Request) {
	var teacherclassroom []TeacherClassroom
	var requestedData NewTeacherClass
	json.NewDecoder(r.Body).Decode(&requestedData)
	db.Where("user_id = ? AND classroom_id = ?", requestedData.UserID, requestedData.ClassroomID).Find(&teacherclassroom)
	//if not found, generate ID
	if len(teacherclassroom) == 0 {
		//save to databse
		var newData = TeacherClassroom{
			ID:          uuid.New().String(),
			UserID:      requestedData.UserID,
			ClassroomID: requestedData.ClassroomID,
		}
		db.Create(&newData)
	}
}

///////////////////////////////////
/// get all class for a teacher ///
///////////////////////////////////
func getTeacherClass(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Executing Get All Classs for Teacher function")
	params := mux.Vars(r)
	var resultClass []Classroom
	db.Raw(`select "classroom".*
	from "user"
	inner join "teacher_classroom" on "user".id = "teacher_classroom".user_id
	inner join "classroom" on "teacher_classroom".classroom_id = "classroom".id
	where "user".id = ?
	order by "classroom"."name" asc`, params["teacher_id"]).Scan(&resultClass)

	json.NewEncoder(w).Encode(resultClass)
}

///////////////////////////////////
//// remove teacher from a user ///
///////////////////////////////////
// delete from teacher_classroom and user
