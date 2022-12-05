package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

type Schedule struct {
	ID        string    `json:"ID"`
	Name      string    `json:"Name"`
	StartDate time.Time `json:"StartDate"`
	EndDate   time.Time `json:"EndDate"`
}

// TableName overrides the table name used by Schedule to `schedule`
func (Schedule) TableName() string {
	return "schedule"
}

func allSchedule(w http.ResponseWriter, r *http.Request) {
	var schedule []Schedule
	db.Find(&schedule)
	json.NewEncoder(w).Encode(schedule)
}

func addSchedule(w http.ResponseWriter, r *http.Request) {

	// format := "2006-01-02"
	// startDate, _ := time.Parse(format, "2019-07-10")
	// endDate, _ := time.Parse(format, "2019-07-11")
	// string format for json: "2022-01-01T00:00:00+08:00"

	var schedule Schedule
	json.NewDecoder(r.Body).Decode(&schedule)
	db.Create(&schedule)
}

func getSchedule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var schedule Schedule
	db.First(&schedule, "id = ?", params["id"])
	json.NewEncoder(w).Encode(schedule)
}

func editSchedule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var schedule Schedule
	db.First(&schedule, "id = ?", params["id"])
	json.NewDecoder(r.Body).Decode(&schedule)
	db.Save(&schedule)
	json.NewEncoder(w).Encode("Successfully edited the schedule.")
}

func removeSchedule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var schedule Schedule
	db.First(&schedule, "id = ?", params["id"])
	db.Delete(&schedule)
	json.NewEncoder(w).Encode("The schedule is deleted successfully!")
}
