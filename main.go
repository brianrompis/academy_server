package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"crypto/sha256"
	"crypto/subtle"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type application struct {
	auth struct {
		username string
		password string
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Homepage Endpoint Hit")
}

func handleRequests() {
	app := new(application)

	app.auth.username = os.Getenv("AUTH_USERNAME")
	app.auth.password = os.Getenv("AUTH_PASSWORD")
	if app.auth.username == "" {
		log.Fatal("Basic auth username must be provided")
	}
	if app.auth.password == "" {
		log.Fatal("Basic auth password must be provided")
	}

	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/", homePage)

	myRouter.HandleFunc("/classroom", createClass).Methods("POST")

	myRouter.HandleFunc("/invitation", createInvitation).Methods("POST")

	myRouter.HandleFunc("/sync", refreshData).Methods("GET")

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

	// list all certificate
	myRouter.HandleFunc("/certificate", allCertificate).Methods("GET")

	myRouter.HandleFunc("/certificate", addCertificate).Methods("POST")
	myRouter.HandleFunc("/certificate/{id}", getCertificate).Methods("GET")
	myRouter.HandleFunc("/certificate/{id}", editCertificate).Methods("PUT")
	myRouter.HandleFunc("/certificate/{id}", removeCertificate).Methods("DELETE")
	myRouter.HandleFunc("/certificate/student/{student_id}", getStudentCertificates).Methods("GET")
	myRouter.HandleFunc("/certificate/student/{student_id}", removeStudentCertificates).Methods("DELETE")

	myRouter.HandleFunc("/employment_history", allEmploymentHistory).Methods("GET")
	myRouter.HandleFunc("/employment_history", addEmploymentHistory).Methods("POST")
	myRouter.HandleFunc("/employment_history/{id}", getEmploymentHistory).Methods("GET")
	myRouter.HandleFunc("/employment_history/{id}", editEmploymentHistory).Methods("PUT")
	myRouter.HandleFunc("/employment_history/{id}", removeEmploymentHistory).Methods("DELETE")
	myRouter.HandleFunc("/employment_history/user/{user_id}", getUserEmploymentHistory).Methods("GET")
	myRouter.HandleFunc("/employment_history/user/{user_id}", removeUserEmploymentHistory).Methods("DELETE")
	myRouter.HandleFunc("/employment_history/delete", deleteMultipleEmploymentHistory).Methods("POST")
	myRouter.HandleFunc("/employment_history/edit", editMultipleEmploymentHistory).Methods("PUT")

	myRouter.HandleFunc("/teacher", allTeacher).Methods("GET")
	myRouter.HandleFunc("/teacher", addTeacher).Methods("POST")

	myRouter.HandleFunc("/user", allUser).Methods("GET")
	myRouter.HandleFunc("/user", addUser).Methods("POST")
	myRouter.HandleFunc("/user/{id}", getUser).Methods("GET")
	myRouter.HandleFunc("/user/{id}", editUser).Methods("PUT")
	myRouter.HandleFunc("/user/{id}", removeUser).Methods("DELETE")

	myRouter.HandleFunc("/education", allEducationHistory).Methods("GET")
	myRouter.HandleFunc("/education", addEducationHistory).Methods("POST")
	myRouter.HandleFunc("/education/{id}", getEducationHistory).Methods("GET")
	myRouter.HandleFunc("/education/{id}", editEducationHistory).Methods("PUT")
	myRouter.HandleFunc("/education/{id}", removeEducationHistory).Methods("DELETE")
	myRouter.HandleFunc("/education/user/{user_id}", getUserEducationHistory).Methods("GET")
	myRouter.HandleFunc("/education/user/{user_id}", removeUserEducationHistory).Methods("DELETE")

	// get all active class
	myRouter.HandleFunc("/active_class", allActiveClassroom).Methods("GET")
	// get all student from a class
	myRouter.HandleFunc("/student/{class_id}", classroomStudent).Methods("GET")
	// get all classroom from a teacher
	myRouter.HandleFunc("/teacher/class/{teacher_id}", getTeacherClass).Methods("GET")
	// user's qualification
	myRouter.HandleFunc("/student/qualification/{user_id}", getStudentQualification).Methods("GET")

	// apply middleware
	var handler http.Handler = myRouter
	handler = app.Auth(handler)

	fmt.Println("handling request")
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"PUT", "GET", "HEAD", "POST", "OPTIONS", "DELETE"})
	headers := handlers.AllowedHeaders([]string{"Authorization"})
	// ttl := handlers.MaxAge(3600)
	origins := handlers.AllowedOrigins([]string{"http://localhost:8081", "http://localhost:8082", "http://localhost:8080"})
	log.Fatal(http.ListenAndServe(":8088", handlers.CORS(credentials, methods, origins, headers)(handler)))
}

func main() {

	InitialMigration()
	fmt.Println("initial migration finished, run handle request")
	handleRequests()

}

func (app *application) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(app.auth.username))
			expectedPasswordHash := sha256.Sum256([]byte(app.auth.password))
			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)
			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
			} else {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				http.Error(w, "Wrong username/password", http.StatusUnauthorized)
				return
			}

		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})
}
