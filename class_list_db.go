package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

type ClassDB struct {
	gorm.Model
	name          string `json:"Name"`
	schedule_id   string `json:"ScheduleId"`
	gClassroom_id string `json:"GoogleClassroomId"`
}

type Schedule struct {
	gorm.Model
	name       string `json:"Name"`
	start_date string `json:"StartDate"`
	end_date   string `json:"EndDate"`
}

func InitialMigration() {
	dsn := "brian:1q2w3e4r!Q@W#E$R@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Faild to connect to database")
	}

	db.AutoMigrate(&Schedule{})
}

func allClassroomsDB(w http.ResponseWriter, r *http.Request) {
	dsn := "brian:1q2w3e4r!Q@W#E$R@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}

	var classes []ClassDB
	db.Find(&classes)
	json.NewEncoder(w).Encode(classes)
}

func postClassesDB(w http.ResponseWriter, r *http.Request) {
	dsn := "brian:1q2w3e4r!Q@W#E$R@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}

	db.Create(&ClassDB{
		name:          "Physics",
		schedule_id:   "4",
		gClassroom_id: "AD54002F",
	})
}

func allSchedule(w http.ResponseWriter, r *http.Request) {
	dsn := "brian:1q2w3e4r!Q@W#E$R@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}

	var schedule []Schedule
	db.Find(&schedule)
	json.NewEncoder(w).Encode(schedule)
}

func addSchedule(w http.ResponseWriter, r *http.Request) {
	dsn := "brian:1q2w3e4r!Q@W#E$R@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}

	format := "2006-01-02"
	// startDate, _ := time.Parse(format, "2019-07-10")
	// endDate, _ := time.Parse(format, "2019-07-11")

	user := Schedule{
		// id:         9,
		name:       "1st Semester 2022",
		start_date: format,
		end_date:   format,
	}

	db.Create(&user)
}
