package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
)

///////////////////////////////////
//////// Get all teacher //////////
///////////////////////////////////
type ReturnedTeacher struct {
	ID                 string    `json:"ID"`
	FullName           string    `json:"FullName"`
	NickName           string    `json:"NickName"`
	Email              string    `json:"Email"`
	Phone              string    `json:"Phone"`
	Birth              time.Time `json:"Birth"`
	Gender             string    `json:"Gender"`
	Address            string    `json:"Address"`
	City               string    `json:"City"`
	Country            string    `json:"Country"`
	Nationality        string    `json:"Nationality"`
	VerificationStatus string    `json:"VerificationStatus"`
	IsVerified         bool      `json:"IsVerified"`
	TeacherStatus      string    `json:"TeacherStatus"`
	IsTeacher          bool      `json:"IsTeacher"`
	EmployeeID         string    `json:"EmployeeID"`
	IDCardNumber       string    `json:"IDCardNumber"`
	IsHRManager        bool      `json:"IsHRManager"`
	IsBanned           bool      `json:"IsBanned"`
}

// select all user with is_teacher true
func allTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []Users
	db.Where("is_teacher = ?", true).Find(&users)
	teacher := []ReturnedTeacher{}
	copier.Copy(&teacher, &users)
	json.NewEncoder(w).Encode(teacher)
}

///////////////////////////////////
////// Register new teacher ///////
///////////////////////////////////
// set is_teacher to true
func registerNewTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	db.Model(&Users{}).Where("id = ?", params["id"]).Update("is_teacher", true)
	json.NewEncoder(w).Encode("Successfully register a teacher.")
}

///////////////////////////////////
//// Add a teacher to a class /////
///////////////////////////////////
type NewTeacherClass struct {
	UserId      uint `json:"UserID"`
	ClassroomId uint `json:"ClassroomID"`
}

func addTeacher(w http.ResponseWriter, r *http.Request) {
	var teacherclassroom []TeacherClassroom
	var requestedData NewTeacherClass
	json.NewDecoder(r.Body).Decode(&requestedData)
	db.Where("user_id = ? AND classroom_id = ?", requestedData.UserId, requestedData.ClassroomId).Find(&teacherclassroom)
	//if not found, generate ID
	if len(teacherclassroom) == 0 {
		//save to databse
		var newData = TeacherClassroom{
			UserId:      requestedData.UserId,
			ClassroomId: requestedData.ClassroomId,
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

//////////////////////////////////
//// testing preload /////////////
//////////////////////////////////

func testPreload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var resultClass []Classroom
	err := db.Order("id").Preload("ClassroomPeriod").Preload("ClassroomPeriod.StudentClassroom").Preload("ClassroomPeriod.StudentClassroom.StudentSubmission").Find(&resultClass)
	if err.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resultClass)
	}
}
