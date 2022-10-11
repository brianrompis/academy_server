package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Qualification struct {
	ID                     string `json:"ID"`
	Name                   string `json:"Name"`
	Description            string `json:"Description"`
	CertificateTemplateID  string `json:"CertificateTemplateID"`
	QualificationClassroom []QualificationClassroom
	JobQualification       []JobQualification
}

func (Qualification) TableName() string {
	return "qualification"
}

type QualificationClassroom struct {
	ID              string `json:"ID"`
	QualificationID string `json:"QualificationID"`
	ClassroomID     string `json:"ClassroomID"`
}

func (QualificationClassroom) TableName() string {
	return "qualification_classroom"
}

func allQualification(w http.ResponseWriter, r *http.Request) {
	var qualification []Qualification
	db.Find(&qualification)
	json.NewEncoder(w).Encode(qualification)
}

type StudentQualification struct {
	QualificationID          string `json:"QualificationID"`
	Qualification            string `json:"Qualification"`
	QualificationClassroomID string `json:"QualificationClassroomID"`
	Classroom                string `json:"Classroom"`
	ClassroomID              string `json:"ClassroomID"`
	UserID                   string `json:"UserID"`
	IsCompleted              bool   `json:"IsCompleted"`
}

type ResultStudentQual struct {
	ClassroomID string `json:"ClassroomID"`
	Classroom   string `json:"Classroom"`
	IsCompleted bool   `json:"IsCompleted"`
}

type ResultQual struct {
	Name    string              `json:"Name"`
	Classes []ResultStudentQual `json:"Classes"`
}

///////////////////////////////////
//// get student Qualification ////
///////////////////////////////////
func getStudentQualification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var studentQualification []StudentQualification
	db.Raw(`select "qualification".id as "qualification_id", "qualification".name as "qualification"
	, "qualification_classroom".id as "qualification_classroom_id" 
	, "classroom".id as "classroom_id", "classroom"."name" as "classroom"
	, "user".id as "user_id",
		case 
			 WHEN "user".id = ? and "student_classroom".has_certificate = true THEN true
		end "is_completed"
	from "qualification_classroom" 
	right join "qualification" on "qualification_classroom".qualification_id  = "qualification".id
	left join "classroom" on "classroom".id = "qualification_classroom".classroom_id
	left join "classroom_period" on "classroom".id = "classroom_period".classroom_id 
	left join "student_classroom" on "classroom_period".id = "student_classroom".classroom_period_id 
	left join "user" on "student_classroom".user_id = "user".id
	order by "qualification" asc`, params["user_id"]).Scan(&studentQualification)

	m := make(map[string]*ResultQual)
	for _, result := range studentQualification {
		qid := result.QualificationID
		_, idPresent := m[qid]

		if !idPresent {
			// record a new qualification
			var resultStudentQual []ResultStudentQual
			theRes := ResultStudentQual{
				ClassroomID: result.ClassroomID,
				Classroom:   result.Classroom,
				IsCompleted: result.IsCompleted,
			}
			resultStudentQual = append(resultStudentQual, theRes)

			m[qid] = &ResultQual{
				Name:    result.Qualification,
				Classes: resultStudentQual,
			}
		} else {
			// check if classroom not yet recorded
			recorded := false
			newVarStatus := false
			index := -1
			for _, newVar := range m[qid].Classes {
				index++
				if result.ClassroomID == newVar.ClassroomID {
					recorded = true
					newVarStatus = newVar.IsCompleted
					break
				}
			}
			// if not recorded, record it
			if !recorded {
				theRes := ResultStudentQual{
					ClassroomID: result.ClassroomID,
					Classroom:   result.Classroom,
					IsCompleted: result.IsCompleted,
				}

				m[qid].Classes = append(m[qid].Classes, theRes)
			} else {
				// if recorded, check if the user id same with param
				if result.UserID == params["user_id"] {
					// if user id same, check if status in record is not completed
					if !newVarStatus {
						// if status in record is not completed, and the new record is completed, update it
						if result.IsCompleted {
							theRes := ResultStudentQual{
								ClassroomID: result.ClassroomID,
								Classroom:   result.Classroom,
								IsCompleted: result.IsCompleted,
							}
							m[qid].Classes[index] = theRes
						}
					}
				}
			}
		}
	}
	// fmt.Print(m)
	json.NewEncoder(w).Encode(m)
}

func addQualification(w http.ResponseWriter, r *http.Request) {
	var qualification []Qualification
	json.NewDecoder(r.Body).Decode(&qualification)
	db.Create(&qualification)
}

func getQualification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var qualification []Qualification
	db.First(&qualification, "id = ?", params["id"])
	json.NewEncoder(w).Encode(qualification)
}

func editQualification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var qualification []Qualification
	db.First(&qualification, "id = ?", params["id"])
	json.NewDecoder(r.Body).Decode(&qualification)
	db.Save(&qualification)
	json.NewEncoder(w).Encode("Successfully edit the qualification.")
}

func removeQualification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var qualification []Qualification
	db.First(&qualification, "id = ?", params["id"])
	db.Delete(&qualification)
	json.NewEncoder(w).Encode("The qualification is deleted successfully!")
}
