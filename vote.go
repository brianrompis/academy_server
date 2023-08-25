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
	if err := db.Where("classroom_id = ?", params["classroom_id"]).Find(&user_vote).Error; err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(user_vote)
	}
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
	if err := db.Where("suggested_classroom_id = ?", params["classroom_id"]).Find(&user_vote).Error; err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(user_vote)
	}
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
	if err := db.Create(&suggestedClassroom).Error; err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode("Successfully added!")
	}
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
	if err := db.Where("id = ?", params["id"]).Find(&suggestedClassroom).Error; err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		resClassroom := []ResultSuggestedClassroom{}
		copier.Copy(&resClassroom, &suggestedClassroom)
		json.NewEncoder(w).Encode(resClassroom)
	}
}

type VoteSuggestedList struct {
	ID                   string `json:"ID"`
	SuggestedClassroomID string `json:"SuggestedClassroomID"`
	VoteType             string `json:"VoteType"`
}

func getUserVoteSuggestedClassroom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var vote_new []VoteNew
	if err := db.Where("user_id = ?", params["user_id"]).Find(&vote_new).Error; err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		voteList := []VoteSuggestedList{}
		copier.Copy(&voteList, &vote_new)
		json.NewEncoder(w).Encode(voteList)
	}
}

type VoteExistingList struct {
	ID          string `json:"ID"`
	ClassroomID string `json:"ClassroomID"`
	VoteType    string `json:"VoteType"`
}

func getUserVoteExistingClassroom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var vote_existing []VoteExisting
	if err := db.Where("user_id = ?", params["user_id"]).Find(&vote_existing).Error; err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		voteList := []VoteExistingList{}
		copier.Copy(&voteList, &vote_existing)
		json.NewEncoder(w).Encode(voteList)
	}
}
