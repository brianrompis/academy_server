package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/classroom/v1"
	"google.golang.org/api/option"
)

///////////////////////////////////////
///////// The Main Function ///////////
///////////////////////////////////////
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
