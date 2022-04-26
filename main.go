package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Classroom struct {
	Id                 string `json:"Id"`
	Name               string `json:"Name"`
	Section            string `json:"Section"`
	DescriptionHeading string `json:"DescriptionHeading"`
	Description        string `json:"Description"`
}

type Classrooms []Classroom

func allClassrooms(w http.ResponseWriter, r *http.Request) {
	classes := showList()

	fmt.Println("Endpoint Hit: All Classes Endpoint")
	json.NewEncoder(w).Encode(classes)
}

func testPostClasses(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test POST endpoint worked")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Homepage Endpoint Hit")
}

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)

	myRouter.HandleFunc("/classes", allClassroomsDB).Methods("GET")
	myRouter.HandleFunc("/classes", postClassesDB).Methods("POST")

	myRouter.HandleFunc("/schedule", allSchedule).Methods("GET")
	myRouter.HandleFunc("/schedule/{id}", getSchedule).Methods("GET")
	myRouter.HandleFunc("/schedule", addSchedule).Methods("POST")
	myRouter.HandleFunc("/schedule/{id}", editSchedule).Methods("PUT")
	myRouter.HandleFunc("/schedule/{id}", removeSchedule).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	InitialMigration()

	handleRequests()
}
