package main

import "time"

type Job struct {
	ID               string `json:"ID"`
	Name             string `json:"Name"`
	Level            string `json:"Level"`
	Department       string `json:"Department"`
	JobQualification []JobQualification
	JobVacancy       []JobVacancy
	OpenCandidate    []OpenCandidate
}

func (Job) TableName() string {
	return "job"
}

type JobQualification struct {
	ID              string `json:"ID"`
	JobID           string `json:"JobID"`
	QualificationID string `json:"QualificationID"`
}

func (JobQualification) TableName() string {
	return "job_qualification"
}

type JobVacancy struct {
	ID                string    `json:"ID"`
	JobID             string    `json:"JobID"`
	HotelID           string    `json:"HotelID"`
	HideHotel         bool      `json:"HideHotel"`
	IssuedDate        time.Time `json:"IssuedDate"`
	ExpiredDate       string    `json:"ExpiredDate"`
	CreatedBy         string    `json:"CreatedBy"`
	Description       string    `json:"Description"`
	Salary            int       `json:"Salary"`
	Status            string    `json:"Status"`
	Type              string    `json:"Type"`
	CandidateSelected string    `json:"CandidateSelected"`
	UserApplication   []UserApplication
}

func (JobVacancy) TableName() string {
	return "job_vacancy"
}

type UserApplication struct {
	ID           string `json:"ID"`
	UserID       string `json:"UserID"`
	JobVacancyID string `json:"JobVacancyID"`
	ApplyDate    string `json:"ApplyDate"`
	Status       string `json:"Status"`
}

func (UserApplication) TableName() string {
	return "user_application"
}

type OpenCandidate struct {
	ID               string    `json:"ID"`
	UserID           string    `json:"UserID"`
	Start            time.Time `json:"Start"`
	UserIntroduction string    `json:"UserIntroduction"`
	IsOpen           bool      `json:"IsOpen"`
	JobID            string    `json:"JobID"`
	OpenLocation     []OpenLocation
}

func (OpenCandidate) TableName() string {
	return "open_candidate"
}

type OpenLocation struct {
	ID              string `json:"ID"`
	OpenCandidateID string `json:"OpenCandidateID"`
	Location        string `json:"Location"`
}

func (OpenLocation) TableName() string {
	return "open_location"
}
