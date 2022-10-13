package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
)

// import (
// 	"encoding/json"
// 	"net/http"

// 	"github.com/gorilla/mux"
// )

type VoteNew struct {
	ID                   string `json:"ID"`
	UserID               string `json:"UserID"`
	SuggestedClassroomID string `json:"SuggestedClassroomID"`
	VoteType             string `json:"VoteType"`
}

func (VoteNew) TableName() string {
	return "vote_new"
}

type VoteExisting struct {
	ID          string `json:"ID"`
	UserID      string `json:"UserID"`
	ClassroomID string `json:"ClassroomID"`
	VoteType    string `json:"VoteType"`
}

func (VoteExisting) TableName() string {
	return "vote_existing"
}

type SuggestedClassroom struct {
	ID              string `json:"ID"`
	Name            string `json:"Name"`
	Description     string `json:"Description"`
	UserSuggestedID string `json:"UserSuggestedID"`
	Status          string `json:"Status"`
	VoteNew         []VoteNew
}

func (SuggestedClassroom) TableName() string {
	return "suggested_classroom"
}

type UserVoteExisting struct {
	ID          string `json:"ID"`
	UserID      string `json:"UserID"`
	ClassroomID string `json:"ClassroomID"`
	VoteType    string `json:"VoteType"`
}

type UserVoteNew struct {
	ID                   string `json:"ID"`
	UserID               string `json:"UserID"`
	SuggestedClassroomID string `json:"SuggestedClassroomID"`
	VoteType             string `json:"VoteType"`
}

func addVoteNew(w http.ResponseWriter, r *http.Request) {
	var vote_new []VoteNew
	json.NewDecoder(r.Body).Decode(&vote_new)
	if err := db.Create(&vote_new).Error; err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode("Added successfully!")
	}
}

func addVoteExisting(w http.ResponseWriter, r *http.Request) {
	var vote_existing []VoteExisting
	json.NewDecoder(r.Body).Decode(&vote_existing)
	if err := db.Create(&vote_existing).Error; err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode("Added successfully!")
	}
}

func getClassroomVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user_vote []VoteExisting
	db.Where("classroom_id = ?", params["classroom_id"]).Find(&user_vote)
	json.NewEncoder(w).Encode(user_vote)
}

func countClassroomVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var voteCount VoteCount
	if err := db.Raw(`select
		sum(case when vote_type = 'up' then 1 else 0 end) as vote_up,
		sum(case when vote_type = 'down' then 1 else 0 end) as vote_down
		FROM vote_existing vn  WHERE classroom_id  = ?`, params["classroom_id"]).Scan(&voteCount).Error; err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(voteCount)
	}
}

func getSuggestedClassroomVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user_vote []VoteNew
	db.Where("suggested_classroom_id = ?", params["classroom_id"]).Find(&user_vote)
	json.NewEncoder(w).Encode(user_vote)
}

type VoteCount struct {
	VoteUp   int `json:"VoteUp"`
	VoteDown int `json:"VoteDown"`
}

func countSuggestedClassroomVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var voteCount VoteCount
	if err := db.Raw(`select
		sum(case when vote_type = 'up' then 1 else 0 end) as vote_up,
		sum(case when vote_type = 'down' then 1 else 0 end) as vote_down
		FROM vote_new vn  WHERE suggested_classroom_id  = ?`, params["classroom_id"]).Scan(&voteCount).Error; err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(voteCount)
	}
}

func addSuggestedClassroom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var suggestedClassroom []SuggestedClassroom
	json.NewDecoder(r.Body).Decode(&suggestedClassroom)
	db.Create(&suggestedClassroom)
	json.NewEncoder(w).Encode("Successfully added!")
}

type ResultSuggestedClassroom struct {
	ID              string `json:"ID"`
	Name            string `json:"Name"`
	Description     string `json:"Description"`
	UserSuggestedID string `json:"UserSuggestedID"`
	Status          string `json:"Status"`
}

func getSuggestedClassroom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var suggestedClassroom []SuggestedClassroom
	db.Where("id = ?", params["id"]).Find(&suggestedClassroom)
	resClassroom := []ResultSuggestedClassroom{}
	copier.Copy(&resClassroom, &suggestedClassroom)
	json.NewEncoder(w).Encode(resClassroom)
}
