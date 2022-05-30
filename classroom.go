package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/classroom/v1"
	"google.golang.org/api/option"
)

////////////////////////////////////////////////
///////////    Google Classroom    /////////////
////////////////////////////////////////////////
type Classroom struct {
	Id                 string `json:"Id"`
	Name               string `json:"Name"`
	Section            string `json:"Section"`
	DescriptionHeading string `json:"DescriptionHeading"`
	Description        string `json:"Description"`
}

type Classrooms []Classroom

func allClassrooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	classes := showList()

	fmt.Println("Endpoint Hit: All Classes Endpoint")
	json.NewEncoder(w).Encode(classes)
}

func createClassroom(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test POST endpoint worked")
}

func showList() Classrooms {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, classroom.ClassroomCoursesReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	//Create a Classroom Client service
	srv, err := classroom.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create classroom Client %v", err)
	}

	//displaying all classes
	r, err := srv.Courses.List().PageSize(150).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve courses. %v", err)
	}
	var classes Classrooms
	if len(r.Courses) > 0 {
		fmt.Print("Courses:\n")
		for _, c := range r.Courses {
			// fmt.Printf("%s (%s)\n", c.Name, c.Id)
			class := Classroom{
				Id:                 c.Id,
				Name:               c.Name,
				Section:            c.Section,
				DescriptionHeading: c.DescriptionHeading,
				Description:        c.Description,
			}
			classes = append(classes, class)
		}
	} else {
		fmt.Print("No courses found.")
	}
	return classes
}

func createClass() {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/classroom.courses.readonly https://www.googleapis.com/auth/classroom.courses")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	//Create a Classroom Client service
	srv, err := classroom.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create classroom Client %v", err)
	}

	// create a class
	c := &classroom.Course{
		Name:               "New Test Classes",
		Section:            "Period 2",
		DescriptionHeading: "It works",
		Description:        "We'll be learning about about the structure of living creatures from a combination of textbooks, guest lectures, and lab work. Expect to be excited!",
		OwnerId:            "me",
		CourseState:        "PROVISIONED",
	}
	course, err := srv.Courses.Create(c).Do()
	if err != nil {
		log.Fatalf("Course unable to be created %v", err)
	}

	fmt.Println(course.Id)

}

///////////////////////////////////////////////////
////////////////  Database   //////////////////////
///////////////////////////////////////////////////

type ClassDB struct {
	ID                string `json:"ID"`
	Name              string `json:"Name" gorm:"column:name"`
	ScheduleId        string `json:"ScheduleId" gorm:"column:schedule_id"`
	GoogleClassroomId string `json:"GoogleClassroomId" gorm:"column:google_classroom_id"`
}

func (ClassDB) TableName() string {
	return "classroom"
}

func allClassroomsDB(w http.ResponseWriter, r *http.Request) {
	var classes []ClassDB
	db.Find(&classes)
	json.NewEncoder(w).Encode(classes)
}

func postClassesDB(w http.ResponseWriter, r *http.Request) {
	var classroom ClassDB
	json.NewDecoder(r.Body).Decode(&classroom)
	db.Create(&classroom)
}

func getClass(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var class ClassDB
	db.First(&class, "id = ?", params["id"])
	json.NewEncoder(w).Encode(class)
}

func editClass(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var class ClassDB
	db.First(&class, "id = ?", params["id"])
	json.NewDecoder(r.Body).Decode(&class)
	db.Save(&class)
	json.NewEncoder(w).Encode("Successfully edit the class.")
}

func removeClass(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var class ClassDB
	db.First(&class, "id = ?", params["id"])
	db.Delete(&class)
	json.NewEncoder(w).Encode("The class is deleted successfully!")
}
