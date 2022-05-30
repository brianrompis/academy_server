package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Homepage Endpoint Hit")
}

func handleRequests() {

	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/", homePage)

	myRouter.HandleFunc("/classroom", allClassrooms).Methods("GET")
	//myRouter.HandleFunc("/classroom", createClass).Methods("POST")

	myRouter.HandleFunc("/classes", allClassroomsDB).Methods("GET")
	myRouter.HandleFunc("/classes", postClassesDB).Methods("POST")
	myRouter.HandleFunc("/classes/{id}", getClass).Methods("GET")
	myRouter.HandleFunc("/classes/{id}", editClass).Methods("PUT")
	myRouter.HandleFunc("/classes/{id}", removeClass).Methods("DELETE")

	myRouter.HandleFunc("/schedule", allSchedule).Methods("GET")
	myRouter.HandleFunc("/schedule", addSchedule).Methods("POST")
	myRouter.HandleFunc("/schedule/{id}", getSchedule).Methods("GET")
	myRouter.HandleFunc("/schedule/{id}", editSchedule).Methods("PUT")
	myRouter.HandleFunc("/schedule/{id}", removeSchedule).Methods("DELETE")

	myRouter.HandleFunc("/certificate", allCertificate).Methods("GET")
	myRouter.HandleFunc("/certificate", addCertificate).Methods("POST")
	myRouter.HandleFunc("/certificate/{id}", getCertificate).Methods("GET")
	myRouter.HandleFunc("/certificate/{id}", editCertificate).Methods("PUT")
	myRouter.HandleFunc("/certificate/{id}", removeCertificate).Methods("DELETE")

	myRouter.HandleFunc("/certification", allCertification).Methods("GET")
	myRouter.HandleFunc("/certification", addCertification).Methods("POST")
	myRouter.HandleFunc("/certification/{id}", getCertification).Methods("GET")
	myRouter.HandleFunc("/certification/{id}", editCertification).Methods("PUT")
	myRouter.HandleFunc("/certification/{id}", removeCertification).Methods("DELETE")

	myRouter.HandleFunc("/employment_history", allEmploymentHistory).Methods("GET")
	myRouter.HandleFunc("/employment_history", addEmploymentHistory).Methods("POST")
	myRouter.HandleFunc("/employment_history/{id}", getEmploymentHistory).Methods("GET")
	myRouter.HandleFunc("/employment_history/{id}", editEmploymentHistory).Methods("PUT")
	myRouter.HandleFunc("/employment_history/{id}", removeEmploymentHistory).Methods("DELETE")

	myRouter.HandleFunc("/student", allStudent).Methods("GET")
	myRouter.HandleFunc("/student", addStudent).Methods("POST")
	myRouter.HandleFunc("/student/{id}", getStudent).Methods("GET")
	myRouter.HandleFunc("/student/{id}", editStudent).Methods("PUT")
	myRouter.HandleFunc("/student/{id}", removeStudent).Methods("DELETE")

	myRouter.HandleFunc("/teacher", allTeacher).Methods("GET")
	myRouter.HandleFunc("/teacher", addTeacher).Methods("POST")
	myRouter.HandleFunc("/teacher/{id}", getTeacher).Methods("GET")
	myRouter.HandleFunc("/teacher/{id}", editTeacher).Methods("PUT")
	myRouter.HandleFunc("/teacher/{id}", removeTeacher).Methods("DELETE")

	myRouter.HandleFunc("/user", allUser).Methods("GET")
	myRouter.HandleFunc("/user", addUser).Methods("POST")
	myRouter.HandleFunc("/user/{id}", getUser).Methods("GET")
	myRouter.HandleFunc("/user/{id}", editUser).Methods("PUT")
	myRouter.HandleFunc("/user/{id}", removeUser).Methods("DELETE")

	fmt.Println("handling request")
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"PUT","GET", "HEAD", "POST", "OPTIONS"})
	// ttl := handlers.MaxAge(3600)
	origins := handlers.AllowedOrigins([]string{"http://localhost:8081"})
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(credentials, methods, origins)(myRouter)))
}

func main() {
	InitialMigration()
	fmt.Println("initial migration finished, run handle request")
	handleRequests()
}
