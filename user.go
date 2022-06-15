package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	ID         string `json:"ID"`
	FirstName       string `json:"FirstName"`
	LastName       string `json:"LastName"`
	Email      string `json:"Email"`
	Phone      string `json:"Phone"`
	Birth      time.Time `json:"Birth"`
	Gender      string `json:"Gender"`
	Address    string `json:"Address"`
	City    string `json:"City"`
	Nationality    string `json:"Nationality"`
	VerificationStatus string   `json:"VerificationStatus"`
	Verified   bool   `json:"Verified"`
	TeacherStatus	string	`json:"TeacherStatus"`
	IsTeacher  bool   `json:"IsTeacher"`
	EmployeeID bool   `json:"EmployeeID"`
}

func (User) TableName() string {
	return "users"
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
