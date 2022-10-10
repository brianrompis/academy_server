package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type DocumentTemplate struct {
	ID            string          `json:"ID"`
	Name          string          `json:"Name"`
	Url           string          `json:"Url"`
	TemplateType  string          `json:"type"`
	Classroom     []Classroom     `gorm:"foreignKey:CertificateTemplateID"`
	Qualification []Qualification `gorm:"foreignKey:CertificateTemplateID"`
}

func (DocumentTemplate) TableName() string {
	return "document_template"
}

func allDocumentTemplate(w http.ResponseWriter, r *http.Request) {
	var documenttemplate []DocumentTemplate
	db.Find(&documenttemplate)
	json.NewEncoder(w).Encode(documenttemplate)
}

func addDocumentTemplate(w http.ResponseWriter, r *http.Request) {
	var documenttemplate []DocumentTemplate
	json.NewDecoder(r.Body).Decode(&documenttemplate)
	db.Create(&documenttemplate)
}

func getDocumentTemplate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var documenttemplate []DocumentTemplate
	db.First(&documenttemplate, "id = ?", params["id"])
	json.NewEncoder(w).Encode(documenttemplate)
}

func editDocumentTemplate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var documenttemplate []DocumentTemplate
	db.First(&documenttemplate, "id = ?", params["id"])
	json.NewDecoder(r.Body).Decode(&documenttemplate)
	db.Save(&documenttemplate)
	json.NewEncoder(w).Encode("Successfully edit the document template.")
}

func removeDocumentTemplate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var documenttemplate []DocumentTemplate
	db.First(&documenttemplate, "id = ?", params["id"])
	db.Delete(&documenttemplate)
	json.NewEncoder(w).Encode("The document template is deleted successfully!")
}
