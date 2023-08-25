package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

/////////////////////////////////////////
/////////// Employment History //////////
/////////////////////////////////////////

func allEmploymentHistory(w http.ResponseWriter, r *http.Request) {
	var employment_history []EmploymentHistory
	db.Find(&employment_history)
	json.NewEncoder(w).Encode(employment_history)
}

func addEmploymentHistory(w http.ResponseWriter, r *http.Request) {
	var employment_history []EmploymentHistory
	json.NewDecoder(r.Body).Decode(&employment_history)
	db.Create(&employment_history)
	json.NewEncoder(w).Encode("Successfully added an employment history.")
}

func getEmploymentHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var employment_history []EmploymentHistory
	db.First(&employment_history, "id = ?", params["id"])
	json.NewEncoder(w).Encode(employment_history)
}

func editEmploymentHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var employment_history []EmploymentHistory
	db.First(&employment_history, "id = ?", params["id"])
	json.NewDecoder(r.Body).Decode(&employment_history)
	db.Save(&employment_history)
	json.NewEncoder(w).Encode("Successfully edited the employment_history.")
}

func removeEmploymentHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var employment_history []EmploymentHistory
	db.First(&employment_history, "id = ?", params["id"])
	db.Delete(&employment_history)
	json.NewEncoder(w).Encode("The employment_history is deleted successfully!")
}

func getUserEmploymentHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var employment_history []EmploymentHistory
	db.Where("user_id = ?", params["user_id"]).Find(&employment_history)
	json.NewEncoder(w).Encode(employment_history)
}

func removeUserEmploymentHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var employment_history []EmploymentHistory
	db.Where("user_id = ?", params["user_id"]).Delete(&employment_history)
	json.NewEncoder(w).Encode("The entire employment histories for the user are deleted successfully!")
}

type DeleteID struct {
	ID string `json:"ID"`
}

func deleteMultipleEmploymentHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var employment_history []EmploymentHistory
	var deleteID []DeleteID
	a := []string{}
	json.NewDecoder(r.Body).Decode(&deleteID)
	for _, s := range deleteID {
		a = append(a, s.ID)
	}
	db.Delete(&employment_history, a)
	json.NewEncoder(w).Encode("The employment_histories are deleted successfully!")
}

func editMultipleEmploymentHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var employment_history []EmploymentHistory
	json.NewDecoder(r.Body).Decode(&employment_history)
	for _, s := range employment_history {
		db.First(&s, "id = ?", s.ID)
		db.Save(&s)
	}
	json.NewEncoder(w).Encode("The employment_histories are edited successfully!")
}

/////////////////////////////////////////
/////////// Education History ///////////
/////////////////////////////////////////

func allEducationHistory(w http.ResponseWriter, r *http.Request) {
	var education []EducationHistory
	db.Find(&education)
	json.NewEncoder(w).Encode(education)
}

func addEducationHistory(w http.ResponseWriter, r *http.Request) {
	var education []EducationHistory
	json.NewDecoder(r.Body).Decode(&education)
	db.Create(&education)
	json.NewEncoder(w).Encode("Successfully added the education history.")
}

func getEducationHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var education []EducationHistory
	db.First(&education, "id = ?", params["id"])
	json.NewEncoder(w).Encode(education)
}

func editEducationHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var education []EducationHistory
	db.First(&education, "id = ?", params["id"])
	json.NewDecoder(r.Body).Decode(&education)
	db.Save(&education)
	json.NewEncoder(w).Encode("Successfully edited the education history.")
}

func removeEducationHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var education []EducationHistory
	db.First(&education, "id = ?", params["id"])
	db.Delete(&education)
	json.NewEncoder(w).Encode("The education history is deleted successfully!")
}

func getUserEducationHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var education []EducationHistory
	db.Where("user_id = ?", params["user_id"]).Find(&education)
	json.NewEncoder(w).Encode(education)
}

func removeUserEducationHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var education []EducationHistory
	db.Where("user_id = ?", params["user_id"]).Delete(&education)
	json.NewEncoder(w).Encode("The entire education histories for the user are deleted successfully!")
}

/////////////////////////////////////////
/////////////// Skill ///////////////////
/////////////////////////////////////////

func addSkill(w http.ResponseWriter, r *http.Request) {
	var skill []Skill
	json.NewDecoder(r.Body).Decode(&skill)
	db.Create(&skill)
	json.NewEncoder(w).Encode("Successfully added the skill.")
}

func getSkill(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var skill []Skill
	db.First(&skill, "id = ?", params["id"])
	json.NewEncoder(w).Encode(skill)
}

func editSkill(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var skill []Skill
	db.First(&skill, "id = ?", params["id"])
	json.NewDecoder(r.Body).Decode(&skill)
	db.Save(&skill)
	json.NewEncoder(w).Encode("Successfully edited the skill.")
}

func removeSkill(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var skill []Skill
	db.First(&skill, "id = ?", params["id"])
	db.Delete(&skill)
	json.NewEncoder(w).Encode("The skill is deleted successfully!")
}

func getUserSkill(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var skill []Skill
	db.Where("user_id = ?", params["user_id"]).Find(&skill)
	json.NewEncoder(w).Encode(skill)
}

func removeUserSkill(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var skill []Skill
	db.Where("user_id = ?", params["user_id"]).Delete(&skill)
	json.NewEncoder(w).Encode("The entire skills for the user are deleted successfully!")
}
