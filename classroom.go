package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/classroom/v1"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

// data from client
type NewClass struct {
	ID           string `json:"ID"`
	Name         string `json:"Name"`
	DepartmentID string `json:"DepartmentID"`
	CreatedBy    string `json:"CreatedBy"`
}

// data stored to DB
type Classroom struct {
	ID                 string          `json:"ID"`
	Name               string          `json:"Name" gorm:"column:name"`
	DepartmentID       string          `json:"DepartmentID" gorm:"column:department_id"`
	GoogleClassroomId  string          `json:"GoogleClassroomId" gorm:"column:google_classroom_id"`
	Link               string          `json:"Link" gorm:"column:link"`
	Status             string          `json:"Status" gorm:"column:status"`
	Public             bool            `json:"Public" gorm:"column:public"`
	PassingGrade       decimal.Decimal `json:"PassingGrade" gorm:"column:passing_grade"`
	Capacity           int             `json:"Capacity" gorm:"column:capacity"`
	ClassStart         time.Time       `json:"ClassStart" gorm:"column:class_start"`
	ClassEnd           time.Time       `json:"ClassEnd" gorm:"column:class_end"`
	RegistrationStart  time.Time       `json:"RegistrationStart" gorm:"column:registration_start"`
	RegistrationEnd    time.Time       `json:"RegistrationEnd" gorm:"column:registration_end"`
	CreatedBy          string          `json:"CreatedBy" gorm:"column:created_by"`
	Section            string          `json:"Section"`
	DescriptionHeading string          `json:"DescriptionHeading"`
	Description        string          `json:"Description"`
}

func (Classroom) TableName() string {
	return "classroom"
}

// data from server(here) to display to front end
type Course struct {
	Id                 string           `json:"Id"`
	Name               string           `json:"Name"`
	DepartmentID       string           `json:"DepartmentID"`
	GoogleClassroomId  string           `json:"GoogleClassroomId"`
	AlternateLink      string           `json:"AlternateLink"`
	Status             string           `json:"Status"`
	Public             bool             `json:"Public"`
	PassingGrade       decimal.Decimal  `json:"PassingGrade"`
	Capacity           int              `json:"Capacity"`
	ClassStart         time.Time        `json:"ClassStart"`
	ClassEnd           time.Time        `json:"ClassEnd"`
	RegistrationStart  time.Time        `json:"RegistrationStart"`
	RegistrationEnd    time.Time        `json:"RegistrationEnd"`
	Section            string           `json:"Section"`
	DescriptionHeading string           `json:"DescriptionHeading"`
	Description        string           `json:"Description"`
	Topics             SimplifiedTopics `json:"Topics"`
}

type Courses []Course

type Topic struct {
	Id          string `json:"Id"`
	Name        string `json:"Name"`
	TopicID     string `json:"TopicID"`
	ClassroomID string `json:"ClassroomID"`
	CourseID    string `json:"CourseID"`
}

type Topics []Topic

type SimplifiedTopic struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}

type SimplifiedTopics []SimplifiedTopic

//////////////////////////////////////
/////////// FUNCTIONS ////////////////
//////////////////////////////////////

// get all classes from Google Classroom

// func allClassrooms(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	classes := getList()

// 	fmt.Println("Endpoint Hit: All Classes Endpoint")
// 	json.NewEncoder(w).Encode(classes)
// }

// create a Classroom Client Service
func classroomClient() *classroom.Service {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, classroom.ClassroomCoursesScope, classroom.ClassroomCoursesReadonlyScope, classroom.ClassroomTopicsReadonlyScope, classroom.ClassroomRostersScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	//Create a Classroom Client service
	srv, err := classroom.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create classroom Client %v", err)
	}
	return srv
}

// type Invitation struct {
// 	ID   string	`json:"id,omitempty"`
// 	CourseId	string	`json:"courseId,omitempty"`
// 	Role	string `json:"role,omitempty"`
// }

// retrieve classroom data from Classroom API

// func getList() Courses {

// 	srv := classroomClient()

// 	// get class list from database

// 	// get class data from Google Classroom
// 	r, err := srv.Courses.List().PageSize(150).Do()
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve courses. %v", err)
// 	}
// 	var classes Courses
// 	if len(r.Courses) > 0 {
// 		fmt.Print("Courses:\n")
// 		for _, c := range r.Courses {
// 			// fmt.Printf("%s (%s)\n", c.Name, c.Id)
// 			class := Course{
// 				Id:                 c.Id,
// 				Name:               c.Name,
// 				Section:            c.Section,
// 				DescriptionHeading: c.DescriptionHeading,
// 				Description:        c.Description,
// 				AlternateLink:      c.AlternateLink,
// 				Topics:             getTopicList(id, c, srv),
// 			}
// 			classes = append(classes, class)

// 		}
// 	} else {
// 		fmt.Print("No courses found.")
// 	}
// 	return classes
// }

func refreshData(w http.ResponseWriter, r *http.Request) {

	updateClassroom()

	var classroom []Classroom
	db.Find(&classroom)

	var courses Courses
	for _, c := range classroom {
		resCourse := Course{
			Id:                 c.ID,
			Name:               c.Name,
			DepartmentID:       c.DepartmentID,
			GoogleClassroomId:  c.GoogleClassroomId,
			AlternateLink:      c.Link,
			Status:             c.Status,
			Public:             c.Public,
			PassingGrade:       c.PassingGrade,
			Capacity:           c.Capacity,
			ClassStart:         c.ClassStart,
			ClassEnd:           c.ClassEnd,
			RegistrationStart:  c.RegistrationStart,
			RegistrationEnd:    c.RegistrationEnd,
			Section:            c.Section,
			DescriptionHeading: c.DescriptionHeading,
			Description:        c.Description,
			Topics:             getTopicFromDB(c.ID),
		}
		courses = append(courses, resCourse)
	}
	json.NewEncoder(w).Encode(&courses)
}

//update classroom data from Classroom API to database
func updateClassroom() {
	srv := classroomClient()

	// get class list from database
	// use INNER JOIN with teacher's class but for now get all classrooms
	var classes []Classroom
	db.Where("").Find(&classes)

	// get class data from Google Classroom
	// METHOD #1: update each classroom individually
	for _, c := range classes {
		// get the class data from Google Classroom
		log.Printf("Get the course data for %s\n", c.Name)
		r, err := srv.Courses.Get(c.GoogleClassroomId).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve course %s. %v", c.Name, err)
		}
		changes := false
		if c.Name != r.Name {
			changes = true
			c.Name = r.Name
		}
		if c.Link != r.AlternateLink {
			changes = true
			c.Link = r.AlternateLink
		}
		if c.Section != r.Section {
			changes = true
			c.Section = r.Section
		}
		if c.DescriptionHeading != r.DescriptionHeading {
			changes = true
			c.DescriptionHeading = r.DescriptionHeading
		}
		if c.Description != r.Description {
			changes = true
			c.Description = r.Description
			fmt.Println(c.Description)
		}

		// update data to database
		if changes {
			fmt.Printf("Updating course %s\n", c.Name)
			errFind := db.First(&c, "id = ?", c.ID).Error
			errors.Is(errFind, gorm.ErrRecordNotFound)
			db.Save(&c)
		}

		//sync the course topics
		syncTopics(c.ID, r.Id, r.Name, srv)
	}

	// get class data from Google Classroom
	// METHOD #2: get all classroom first
}

// get list of topics for a couse from Classroom API
func getTopicList(classroomID string, courseId string, courseName string, srv *classroom.Service) Topics {
	res, err2 := srv.Courses.Topics.List(courseId).PageSize(30).Do()
	if err2 != nil {
		log.Fatalf("Unable to retrieve topics for classroom: %s", courseName)
	}
	var topics Topics
	for _, t := range res.Topic {
		topicUniqueID := uuid.New()
		resTopic := Topic{
			Id:          topicUniqueID.String(),
			TopicID:     t.TopicId,
			Name:        t.Name,
			ClassroomID: classroomID,
			CourseID:    t.CourseId,
		}
		topics = append(topics, resTopic)
	}
	return topics
}

// get list of topics(simplified) for a couse from Classroom API and update DB accordingly
func syncTopics(classroomID string, courseId string, courseName string, srv *classroom.Service) {
	//get list of topics from classroom
	log.Printf("Get the topics data for %s\n", courseName)
	res, err2 := srv.Courses.Topics.List(courseId).PageSize(30).Do()
	if err2 != nil {
		log.Fatalf("Unable to retrieve topics for classroom: %s", courseName)
	}

	if len(res.Topic) > 0 {

		//get list of topics from DB
		var topicsDB Topics
		db.Where("classroom_id = ?", classroomID).Find(&topicsDB)

		//map topic in DB
		var mappedTopic = map[string]Topic{}
		for _, t := range topicsDB {
			mappedTopic[t.TopicID] = t
		}

		//update topic in DB according to Classroom
		for _, t := range res.Topic {
			if val, ok := mappedTopic[t.TopicId]; ok {
				// check if the topic name still same
				if val.Name != t.Name {
					// update the topic name in database
					fmt.Printf("Updating name for topic %s\n", t.Name)
					db.Model(&Topic{}).Where("id = ?", val.Id).Update("name", t.Name)
				}
				// remove the matched record in map in order delete the rest from DB
				delete(mappedTopic, t.TopicId)
			} else {
				// add new topic to database
				fmt.Printf("Add new topic %s\n", t.Name)
				newTopic := Topic{
					Id:          uuid.New().String(),
					Name:        t.Name,
					TopicID:     t.TopicId,
					ClassroomID: classroomID,
					CourseID:    courseId,
				}
				db.Create(&newTopic)
			}
		}
		//delete topic in DB if no longer exist in Classroom based on records left in the map
		for key, val := range mappedTopic {
			fmt.Printf("Deleting %s from database\n", val.Name)
			db.Where("topic_id = ?", key).Delete(&val)
		}
	} else {
		fmt.Printf("There is no topic for %s\n", courseName)
		// delete all topics from DB if in classroom is not found
		var topics Topics
		db.Where("course_id = ?", courseId).Delete(&topics)
	}
}

// get list of topics(simplified) for a couse from database
func getTopicFromDB(classroomID string) SimplifiedTopics {
	var topics Topics
	db.Where("classroom_id = ?", classroomID).Find(&topics)
	var simplifiedTopics SimplifiedTopics
	for _, t := range topics {
		resTopic := SimplifiedTopic{
			Id:   t.Id,
			Name: t.Name,
		}
		simplifiedTopics = append(simplifiedTopics, resTopic)
	}
	return simplifiedTopics
}

//get specific topic from Classroom API based on id
func getSpecificTopic(topicID string, classroomID string, c *classroom.Course, srv *classroom.Service) Topic {
	res, err := srv.Courses.Topics.Get(c.Id, topicID).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve topic with id: %s", topicID)
	}
	topicUniqueID := uuid.New()
	resTopic := Topic{
		Id:          topicUniqueID.String(),
		TopicID:     res.TopicId,
		Name:        res.Name,
		ClassroomID: classroomID,
		CourseID:    res.CourseId,
	}
	return resTopic
}

// create a new class course
func createClass(w http.ResponseWriter, r *http.Request) {

	srv := classroomClient()

	//extract data from rest api
	var classroomDetail []NewClass
	json.NewDecoder(r.Body).Decode(&classroomDetail)

	for _, c := range classroomDetail {
		fmt.Println(c.Name)

		// create a class in Google Classroom
		newClass := &classroom.Course{
			Name:               c.Name,
			Section:            "",
			DescriptionHeading: "",
			Description:        "",
			OwnerId:            "me",
			CourseState:        "PROVISIONED",
		}
		course, err := srv.Courses.Create(newClass).Do()
		if err != nil {
			log.Fatalf("Course unable to be created %v", err)
		}
		fmt.Println(course.Id)

		//send invitation to teacher
		sendTeacherInvitation(course.Id, "brian.rompis@gmail.com", srv)

		// save to Database
		initTime := time.Now()
		defaultGrade := decimal.NewFromInt(0)
		defaultCapacity := 500
		var classroom = Classroom{
			ID:                 c.ID,
			Name:               c.Name,
			DepartmentID:       c.DepartmentID,
			GoogleClassroomId:  course.Id,
			Link:               course.AlternateLink,
			Status:             "REVIEWED",
			Public:             false,
			PassingGrade:       defaultGrade,
			Capacity:           defaultCapacity,
			ClassStart:         initTime,
			ClassEnd:           initTime,
			RegistrationStart:  initTime,
			RegistrationEnd:    initTime,
			CreatedBy:          c.CreatedBy,
			Section:            "",
			DescriptionHeading: "",
			Description:        "",
		}
		db.Create(&classroom)
		json.NewEncoder(w).Encode("Successfully stored the new Classroom to database.")
	}
}

func sendTeacherInvitation(id string, user string, srv *classroom.Service) {
	body := &classroom.Invitation{
		CourseId: id,
		Role:     "TEACHER",
		UserId:   user,
	}
	_, err0 := srv.Invitations.Create(body).Do()
	if err0 != nil {
		log.Fatalf("Unable to send teacher invitation to %s %v", user, err0)
	}
}

type Invitation struct {
	Email       string `json:"Email"`
	ClassroomID string `json:"ClassroomID"`
	Role        string `json:"Role"`
}

func createInvitation(w http.ResponseWriter, r *http.Request) {
	srv := classroomClient()

	var invitations []Invitation
	json.NewDecoder(r.Body).Decode(&invitations)
	for _, i := range invitations {
		sendInvitation(i.ClassroomID, i.Email, i.Role, srv)
	}
}

func sendInvitation(id string, user string, role string, srv *classroom.Service) {
	body := &classroom.Invitation{
		CourseId: id,
		Role:     role,
		UserId:   user,
	}
	_, err := srv.Invitations.Create(body).Do()
	if err != nil {
		log.Fatalf("Unable to send %s invitation to %s %v", role, user, err)
	}
}

///////////////////////////////////////////////////
///////////  Classroom to Database   //////////////
///////////////////////////////////////////////////

func allClassroomsDB(w http.ResponseWriter, r *http.Request) {
	var classroom []Classroom
	db.Find(&classroom)

	var courses Courses
	for _, c := range classroom {
		resCourse := Course{
			Id:                 c.ID,
			Name:               c.Name,
			DepartmentID:       c.DepartmentID,
			GoogleClassroomId:  c.GoogleClassroomId,
			AlternateLink:      c.Link,
			Status:             c.Status,
			Public:             c.Public,
			PassingGrade:       c.PassingGrade,
			Capacity:           c.Capacity,
			ClassStart:         c.ClassStart,
			ClassEnd:           c.ClassEnd,
			RegistrationStart:  c.RegistrationStart,
			RegistrationEnd:    c.RegistrationEnd,
			Section:            c.Section,
			DescriptionHeading: c.DescriptionHeading,
			Description:        c.Description,
			Topics:             getTopicFromDB(c.ID),
		}
		courses = append(courses, resCourse)
	}
	json.NewEncoder(w).Encode(&courses)
}

func postClassesDB(w http.ResponseWriter, r *http.Request) {
	var classroom Classroom
	json.NewDecoder(r.Body).Decode(&classroom)
	db.Create(&classroom)
}

func getClass(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var class Classroom
	db.First(&class, "id = ?", params["id"])
	json.NewEncoder(w).Encode(class)
}

func editClass(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var class Classroom
	db.First(&class, "id = ?", params["id"])
	json.NewDecoder(r.Body).Decode(&class)
	db.Save(&class)
	json.NewEncoder(w).Encode("Successfully edit the class.")
}

func removeClass(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var class Classroom
	db.First(&class, "id = ?", params["id"])
	db.Delete(&class)
	json.NewEncoder(w).Encode("The class is deleted successfully!")
}
