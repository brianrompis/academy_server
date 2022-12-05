package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Certificate struct {
	ClassroomID string    `json:"ClassroomID"`
	Classroom   string    `json:"Classroom"`
	TemplateID  string    `json:"TemplateID"`
	TemplateUrl string    `json:"TemplateUrl"`
	Grade       float64   `json:"Grade"`
	IssuedDate  time.Time `json:"IssuedDate"`
	ExpiredDate time.Time `json:"ExpiredDate"`
	PeriodStart time.Time `json:"PeriodStart"`
	PeriodEnd   time.Time `json:"PeriodEnd"`
}

func getStudentCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var certificates []Certificate
	if err := db.Raw(`select c.id as "classroom_id" , c."name" as "classroom", c.certificate_template_id as "template_id"
	, dt.url as "template_url"
	, sc.grade, sc.certificate_issued_date as "issued_date"
	, cp.cert_expired_date as "expired_date", cp.start_date as "period_start", cp.end_date as "period_end"
	from student_classroom sc
	inner join classroom_period cp on sc.classroom_period_id = cp.id 
	inner join classroom c on cp.classroom_id = c.id 
	inner join document_template dt on c.certificate_template_id = dt.id 
	where sc.has_certificate = true
	and sc.user_id = ?`, params["user_id"]).Scan(&certificates).Error; err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(certificates)
	}
}
