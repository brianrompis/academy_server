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

func inviteStudent() {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, classroom.ClassroomRostersScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	///////////////////////////////////
	//Create a Classroom Client service
	///////////////////////////////////
	srv, err := classroom.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create classroom Client %v", err)
	}

	///////////////////
	// enroll a student
	///////////////////

	courseId := "76215154252"
	studentId := "brian.rompis@gmail.com"

	//admin only
	//ret, err := srv.Courses.Students.Create(courseId, &classroom.Student{UserId: studentId}).Do()

	i := &classroom.Invitation{
		UserId:   studentId,
		CourseId: courseId,
		Role:     "STUDENT",
	}
	ret, err := srv.Invitations.Create(i).Do()

	if err != nil {
		log.Fatalf("Course unable to be created %v", err)
	}
	fmt.Println(ret)
}
