package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
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
	IsHRManager        bool      `json:"IsHRManager" gorm:"column:is_hr_manager"`
	IsBanned           bool      `json:"IsBanned"`
	TeacherClassroom   []TeacherClassroom
	Classroom          []Classroom `gorm:"foreignKey:CreatedBy"`
	StudentClassroom   []StudentClassroom
	UserVote           []UserVote
	SuggestedClassroom []SuggestedClassroom `gorm:"foreignKey:UserSuggestedID"`
	JobVacancy         []JobVacancy         `gorm:"foreignKey:CreatedBy"`
	UserApplication    []UserApplication
	EmploymentHistory  []EmploymentHistory
	EducationHistory   []EducationHistory
	Skill              []Skill
	OpenCandidate      []OpenCandidate
}

func (User) TableName() string {
	return "user"
}

func allUser(w http.ResponseWriter, r *http.Request) {
	var user []User
	db.Find(&user)
	json.NewEncoder(w).Encode(user)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	var user []User
	json.NewDecoder(r.Body).Decode(&user)
	db.Create(&user)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		w.Write([]byte("allowed"))
		return
	}
	params := mux.Vars(r)
	var user []User
	db.First(&user, "id = ?", params["id"])
	json.NewEncoder(w).Encode(user)
}

func editUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user []User
	db.First(&user, "id = ?", params["id"])
	json.NewDecoder(r.Body).Decode(&user)
	db.Save(&user)
	json.NewEncoder(w).Encode("Successfully edit the user.")
}

func removeUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user []User
	db.First(&user, "id = ?", params["id"])
	db.Delete(&user)
	json.NewEncoder(w).Encode("The user is deleted successfully!")
}
