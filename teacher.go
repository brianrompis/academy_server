package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Teacher struct {
	ID          string `json:"ID"`
	UserID      string `json:"UserID"`
	ClassroomID string `json:"ClassroomID"`
	CourseId    string `json:"CourseID"`
}

func (Teacher) TableName() string {
	return "teachers"
}

func allTeacher(w http.ResponseWriter, r *http.Request) {
	var teacher []Teacher
	db.Find(&teacher)
	json.NewEncoder(w).Encode(teacher)
}

func addTeacher(w http.ResponseWriter, r *http.Request) {
	var teacher []Teacher
	json.NewDecoder(r.Body).Decode(&teacher)
	db.Create(&teacher)
}

func getTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var teacher []Teacher
	db.First(&teacher, "id = ?", params["id"])
	json.NewEncoder(w).Encode(teacher)
}

func editTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var teacher []Teacher
	db.First(&teacher, "id = ?", params["id"])
	json.NewDecoder(r.Body).Decode(&teacher)
	db.Save(&teacher)
	json.NewEncoder(w).Encode("Successfully edit the teacher.")
}

func removeTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var teacher []Teacher
	db.First(&teacher, "id = ?", params["id"])
	db.Delete(&teacher)
	json.NewEncoder(w).Encode("The teacher is deleted successfully!")
}
