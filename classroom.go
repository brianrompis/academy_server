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
	"github.com/jinzhu/copier"
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
	DepartmentId string `json:"DepartmentId"`
	CreatedBy    string `json:"CreatedBy"`
}

// data stored to DB
type Classroom struct {
	ID                     string          `json:"ID" gorm:"column:id"`
	Name                   string          `json:"Name" gorm:"column:name"`
	GoogleClassroomId      string          `json:"GoogleClassroomId" gorm:"column:google_classroom_id"`
	Link                   string          `json:"Link" gorm:"column:link"`
	Status                 string          `json:"Status" gorm:"column:status"`
	IsPublic               bool            `json:"IsPublic" gorm:"column:is_public"`
	PassingGrade           decimal.Decimal `json:"PassingGrade" gorm:"column:passing_grade"`
	Capacity               int             `json:"Capacity" gorm:"column:capacity"`
	CreatedBy              string          `json:"CreatedBy" gorm:"column:created_by"`
	Section                string          `json:"Section" gorm:"column:section"`
	DescriptionHeading     string          `json:"DescriptionHeading" gorm:"column:description_heading"`
	Description            string          `json:"Description" gorm:"column:description"`
	CertificateTemplateID  string          `json:"CertificateTemplateID" gorm:"column:certificate_template_id"`
	DepartmentId           string          `json:"DepartmentId" gorm:"column:department_id"`
	ActivePeriodID         string          `json:"ActivePeriodID" gorm:"column:active_period_id"`
	Topic                  []Topic
	ClassroomPeriod        []ClassroomPeriod
	Assignment             []Assignment
	VoteExisting           []VoteExisting
	TeacherClassroom       []TeacherClassroom
	QualificationClassroom []QualificationClassroom
}

func (Classroom) TableName() string {
	return "classroom"
}

type Department struct {
	ID        string `json:"ID" gorm:"column:id"`
	Name      string `json:"Name" gorm:"column:name"`
	Classroom []Classroom
}

func (Department) TableName() string {
	return "department"
}

type ClassroomPeriod struct {
	ID                 string    `json:"ID"`
	ClassroomID        string    `json:"ClassroomID"`
	StartDate          time.Time `json:"StartDate"`
	EndDate            time.Time `json:"EndDate"`
	CertExpiredDate    time.Time `json:"CertExpiredDate"`
	RegistrationPeriod []RegistrationPeriod
	StudentClassroom   []StudentClassroom
	Classroom          Classroom `gorm:"foreignKey:ActivePeriodID"`
}

func (ClassroomPeriod) TableName() string {
	return "classroom_period"
}

// add new classroom period
func addClassroomPeriod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var classroom_period []ClassroomPeriod
	json.NewDecoder(r.Body).Decode(&classroom_period)
	for _, p := range classroom_period {
		//a = append(a, p.ID)
		db.Exec(`with ap as (
			insert into "classroom_period"(id, classroom_id, start_date, end_date, cert_expired_date)
			values(?, ?, ?, ?, ?)
			returning id
		) update classroom set active_period_id = ? where "classroom".id = ?`,
			p.ID, p.ClassroomID, p.StartDate, p.EndDate, p.CertExpiredDate, p.ID, p.ClassroomID)
	}
	json.NewEncoder(w).Encode("Successfully add the classroom period.")
}

type RegistrationPeriod struct {
	ID                string    `json:"ID"`
	ClassroomPeriodID string    `json:"ClassroomPeriodID"`
	StartDate         time.Time `json:"StartDate"`
	EndDate           time.Time `json:"EndDate"`
}

func (RegistrationPeriod) TableName() string {
	return "registration_period"
}

type Assignment struct {
	ID                 string    `json:"ID"`
	ClassroomID        string    `json:"ClassroomID"`
	GoogleCourseworkID time.Time `json:"GoogleCourseworkID"`
	StudentSubmission  []StudentSubmission
}

func (Assignment) TableName() string {
	return "assignment"
}

// data from server(here) to display to front end
type Course struct {
	Id                 string           `json:"Id"`
	Name               string           `json:"Name"`
	DepartmentId       string           `json:"DepartmentId"`
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
	Id            string `json:"Id"`
	Name          string `json:"Name"`
	GoogleTopicID string `json:"GoogleTopicID"`
	ClassroomID   string `json:"ClassroomID"`
}

func (Topic) TableName() string {
	return "topic"
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
			GoogleClassroomId:  c.GoogleClassroomId,
			AlternateLink:      c.Link,
			Status:             c.Status,
			Public:             c.IsPublic,
			PassingGrade:       c.PassingGrade,
			Capacity:           c.Capacity,
			DepartmentId:       c.DepartmentId,
			Section:            c.Section,
			DescriptionHeading: c.DescriptionHeading,
			Description:        c.Description,
			// ClassStart:         c.ClassStart,
			// ClassEnd:           c.ClassEnd,
			// RegistrationStart:  c.RegistrationStart,
			// RegistrationEnd:    c.RegistrationEnd,
			Topics: getTopicFromDB(c.ID),
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
			Id:            topicUniqueID.String(),
			GoogleTopicID: t.TopicId,
			Name:          t.Name,
			ClassroomID:   classroomID,
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
			mappedTopic[t.GoogleTopicID] = t
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
					Id:            uuid.New().String(),
					Name:          t.Name,
					GoogleTopicID: t.TopicId,
					ClassroomID:   classroomID,
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
		Id:            topicUniqueID.String(),
		GoogleTopicID: res.TopicId,
		Name:          res.Name,
		ClassroomID:   classroomID,
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
		// initTime := time.Now()
		defaultGrade := decimal.NewFromInt(0)
		defaultCapacity := 500
		var classroom = Classroom{
			ID:                 c.ID,
			Name:               c.Name,
			DepartmentId:       c.DepartmentId,
			GoogleClassroomId:  course.Id,
			Link:               course.AlternateLink,
			Status:             "REVIEWED",
			IsPublic:           false,
			PassingGrade:       defaultGrade,
			Capacity:           defaultCapacity,
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
	w.Header().Set("Content-Type", "application/json")
	var classroom []Classroom
	db.Find(&classroom)

	var courses Courses
	for _, c := range classroom {
		resCourse := Course{
			Id:                c.ID,
			Name:              c.Name,
			DepartmentId:      c.DepartmentId,
			GoogleClassroomId: c.GoogleClassroomId,
			AlternateLink:     c.Link,
			Status:            c.Status,
			Public:            c.IsPublic,
			PassingGrade:      c.PassingGrade,
			Capacity:          c.Capacity,
			// ClassStart:         c.ClassStart,
			// ClassEnd:           c.ClassEnd,
			// RegistrationStart:  c.RegistrationStart,
			// RegistrationEnd:    c.RegistrationEnd,
			Section:            c.Section,
			DescriptionHeading: c.DescriptionHeading,
			Description:        c.Description,
			Topics:             getTopicFromDB(c.ID),
		}
		courses = append(courses, resCourse)
	}
	json.NewEncoder(w).Encode(&courses)
}

type AvailableClassroom struct {
	ID                    string          `json:"ID"`
	Name                  string          `json:"Name"`
	GoogleClassroomId     string          `json:"GoogleClassroomId"`
	Link                  string          `json:"Link"`
	Status                string          `json:"Status"`
	IsPublic              bool            `json:"IsPublic"`
	PassingGrade          decimal.Decimal `json:"PassingGrade"`
	Capacity              int             `json:"Capacity"`
	CreatedBy             string          `json:"CreatedBy"`
	Section               string          `json:"Section"`
	DescriptionHeading    string          `json:"DescriptionHeading"`
	Description           string          `json:"Description"`
	CertificateTemplateID string          `json:"CertificateTemplateID"`
	DepartmentId          string          `json:"DepartmentId"`
	ActivePeriodID        string          `json:"ActivePeriodID"`
	StartDate             time.Time       `json:"StartDate"`
	EndDate               time.Time       `json:"EndDate"`
	RegistrationStart     time.Time       `json:"RegistrationStart"`
	RegistrationEnd       time.Time       `json:"RegistrationEnd"`
}

// get all available classroom (today within the registration period)
func allAvailableClassroom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var availableClassroom []AvailableClassroom
	db.Raw(`select "classroom".*, "classroom_period".start_date, "classroom_period".end_date, "registration_period".start_date  as "registration_start", "registration_period".end_date as "registration_end"
	from "registration_period"
	inner join "classroom_period" on "classroom_period".id = "registration_period".classroom_period_id
	inner join "classroom" on "classroom_period".classroom_id = "classroom".id 
	where current_date between  "registration_period".start_date and "registration_period".end_date 
	order by "registration_period".start_date`).Scan(&availableClassroom)

	json.NewEncoder(w).Encode(availableClassroom)
}

// get all ongoing classroom
func allOngoingClassroom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var availableClassroom []AvailableClassroom
	db.Raw(`select "classroom".*, "classroom_period".start_date, "classroom_period".end_date, "registration_period".start_date  as "registration_start", "registration_period".end_date as "registration_end"
	from "classroom_period"
	inner join "classroom" on "classroom_period".classroom_id = "classroom".id 
	left join "registration_period" on "classroom_period".id = "registration_period".classroom_period_id  
	where current_date between "classroom_period".start_date and "classroom_period".end_date
	order by "classroom_period".start_date`).Scan(&availableClassroom)

	json.NewEncoder(w).Encode(availableClassroom)
}

// get all reviewed classroom
func allReviewedClassroom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var classroom []Classroom
	db.Where("status = 'REVIEWED'").Find(&classroom)
	resClassroom := []AvailableClassroom{}
	copier.Copy(&resClassroom, &classroom)
	json.NewEncoder(w).Encode(resClassroom)
}

// get all student classroom
func allStudentClassroom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var availableClassroom []AvailableClassroom
	db.Raw(`select "classroom".*, "classroom_period".start_date, "classroom_period".end_date, "registration_period".start_date  as "registration_start", "registration_period".end_date as "registration_end"
	from "classroom_period"
	inner join "classroom" on "classroom_period".classroom_id = "classroom".id 
	left join "registration_period" on "classroom_period".id = "registration_period".classroom_period_id  
	where current_date between "classroom_period".start_date and "classroom_period".end_date
	order by "classroom_period".start_date`).Scan(&availableClassroom)

	json.NewEncoder(w).Encode(availableClassroom)
}

// get all active classroom (public)
func allActivePublicClassroom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
