package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
)

func allUser(w http.ResponseWriter, r *http.Request) {
	var user []Users
	db.Find(&user)
	json.NewEncoder(w).Encode(user)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user []Users
	json.NewDecoder(r.Body).Decode(&user)
	if err := db.Create(&user).Error; err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode("Added successfully.")
	}
}

func countUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userCount int64
	if err := db.Model(&Users{}).Count(&userCount).Error; err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(userCount)
	}
}

type ResultUser struct {
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
	IsAdmin            bool      `json:"IsAdmin"`
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		w.Write([]byte("allowed"))
		return
	}
	params := mux.Vars(r)
	var user Users
	db.First(&user, "id = ?", params["id"])
	resUser := ResultUser{}
	copier.Copy(&resUser, &user)
	json.NewEncoder(w).Encode(resUser)
}

func editUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user []Users
	db.First(&user, "id = ?", params["id"])
	json.NewDecoder(r.Body).Decode(&user)
	db.Save(&user)
	json.NewEncoder(w).Encode("Successfully edited the user.")
}

func removeUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user []Users
	db.First(&user, "id = ?", params["id"])
	db.Delete(&user)
	json.NewEncoder(w).Encode("The user is deleted successfully!")
}

type Status struct {
	IsTeacher   bool
	IsHRManager bool
	IsAdmin     bool
}

func getRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var status Status
	db.Model(&Users{}).Where("id = ?", params["id"]).Find(&status)
	json.NewEncoder(w).Encode(status)
}

func isArchipelagoEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user Users
	params := mux.Vars(r)
	db.First(&user, "id = ?", params["id"])
	json.NewEncoder(w).Encode(user.EmployeeId)
}
