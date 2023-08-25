package main

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}

//////////////////////////////////////////////
/// MODEL DECLARATION
//////////////////////////////////////////////

// data stored to DB
type Classroom struct {
	gorm.Model
	Name                   string
	GoogleClassroomId      string
	Link                   string
	Status                 string
	IsPublic               bool
	PassingGrade           decimal.Decimal
	Capacity               int
	CreatedBy              uint `gorm:"not null"`
	Section                string
	DescriptionHeading     string
	Description            string
	IsDisabled             bool
	CertificateTemplateId  uint                     `gorm:"default:null"`
	DepartmentId           uint                     `gorm:"default:null"`
	ActivePeriodId         uint                     `gorm:"default:null"`
	Topic                  []Topic                  `gorm:"foreignKey:ClassroomId"`
	ClassroomPeriod        []ClassroomPeriod        `gorm:"foreignKey:ClassroomId"`
	Assignment             []Assignment             `gorm:"foreignKey:ClassroomId"`
	VoteExisting           []VoteExisting           `gorm:"foreignKey:ClassroomId"`
	TeacherClassroom       []TeacherClassroom       `gorm:"foreignKey:ClassroomId"`
	QualificationClassroom []QualificationClassroom `gorm:"foreignKey:ClassroomId"`
	ClassroomCourseTag     []ClassroomCourseTag     `gorm:"foreignKey:ClassroomId"`
}

func (Classroom) TableName() string {
	return "classroom"
}

type Department struct {
	gorm.Model
	Name      string `json:"Name" gorm:"column:name"`
	Classroom []Classroom
}

func (Department) TableName() string {
	return "department"
}

type ClassroomPeriod struct {
	gorm.Model
	ClassroomId        uint                 `gorm:"not null"`
	StartDate          time.Time            `json:"StartDate"`
	EndDate            time.Time            `json:"EndDate"`
	CertExpiredDate    time.Time            `json:"CertExpiredDate"`
	RegistrationPeriod []RegistrationPeriod `gorm:"foreignKey:ClassroomPeriodId"`
	StudentClassroom   []StudentClassroom   `gorm:"foreignKey:ClassroomPeriodId"`
}

func (ClassroomPeriod) TableName() string {
	return "classroom_period"
}

type RegistrationPeriod struct {
	gorm.Model
	ClassroomPeriodId uint      `json:"ClassroomPeriodId"`
	StartDate         time.Time `json:"StartDate"`
	EndDate           time.Time `json:"EndDate"`
}

func (RegistrationPeriod) TableName() string {
	return "registration_period"
}

type Assignment struct {
	gorm.Model
	ClassroomId        uint `gorm:"not null"`
	GoogleCourseworkId time.Time
	StudentSubmission  []StudentSubmission `gorm:"foreignKey:AssignmentId"`
}

func (Assignment) TableName() string {
	return "assignment"
}

type Topic struct {
	gorm.Model
	Name          string
	GoogleTopicId string
	ClassroomId   uint `gorm:"not null"`
}

func (Topic) TableName() string {
	return "topic"
}

type EmploymentHistory struct {
	gorm.Model
	UserId            uint `gorm:"not null"`
	Position          string
	CompanyName       string
	CompanyPublicName string
	StartDate         time.Time
	EndDate           time.Time
	City              string
	Country           string
	Description       string
	IsArchipelago     bool
	HotelId           string
	ContactName       string
	ContactPosition   string
	ContactEmail      string
	ContactPhone      string
}

func (EmploymentHistory) TableName() string {
	return "employment_history"
}

type EducationHistory struct {
	gorm.Model
	UserId         uint `gorm:"not null"`
	EducationLevel string
	SchoolName     string
	StartYear      time.Time
	EndYear        time.Time
	GPA            float64 `gorm:"column:gpa" gorm:"type:numeric(5,2)"`
}

func (EducationHistory) TableName() string {
	return "education_history"
}

type Skill struct {
	gorm.Model
	UserId      string `gorm:"not null"`
	Name        string
	Description string
	Level       int
}

func (Skill) TableName() string {
	return "skill"
}

type Qualification struct {
	gorm.Model
	Name                   string                   `json:"Name"`
	Description            string                   `json:"Description"`
	CertificateTemplateId  uint                     `gorm:"default:null"`
	QualificationClassroom []QualificationClassroom `gorm:"foreignKey:QualificationId"`
	JobQualification       []JobQualification       `gorm:"foreignKey:QualificationId"`
}

func (Qualification) TableName() string {
	return "qualification"
}

type QualificationClassroom struct {
	gorm.Model
	QualificationId uint `gorm:"not null"`
	ClassroomId     uint `gorm:"not null"`
}

func (QualificationClassroom) TableName() string {
	return "qualification_classroom"
}

type Schedule struct {
	gorm.Model
	Name      string    `json:"Name"`
	StartDate time.Time `json:"StartDate"`
	EndDate   time.Time `json:"EndDate"`
}

// TableName overrides the table name used by Schedule to `schedule`
func (Schedule) TableName() string {
	return "schedule"
}

type StudentClassroom struct {
	gorm.Model
	UserId                uint `gorm:"not null"`
	ClassroomPeriodId     uint `gorm:"not null"`
	Status                string
	Grade                 float64 `gorm:"type:numeric(5,2)"`
	HasCertificate        bool
	CertificateIssuedDate time.Time
	StudentSubmission     []StudentSubmission `gorm:"foreignKey:StudentClassroomId"`
}

func (StudentClassroom) TableName() string {
	return "student_classroom"
}

type StudentSubmission struct {
	gorm.Model
	StudentClassroomId uint    `gorm:"not null"`
	AssignmentId       uint    `gorm:"not null"`
	DraftGrade         float64 `gorm:"type:numeric(5,2)"`
	AssignedGrade      float64 `gorm:"type:numeric(5,2)"`
	LastUpdate         time.Time
	LastAssignedBy     uint                 `gorm:"default:null"`
	GradeAssignHistory []GradeAssignHistory `gorm:"foreignKey:StudentSubmissionId"`
}

func (StudentSubmission) TableName() string {
	return "student_submission"
}

type TeacherClassroom struct {
	gorm.Model
	UserId             uint                 `gorm:"not null"`
	ClassroomId        uint                 `gorm:"not null"`
	StudentSubmission  []StudentSubmission  `gorm:"foreignKey:LastAssignedBy"`
	GradeAssignHistory []GradeAssignHistory `gorm:"foreignKey:TeacherId"`
}

func (TeacherClassroom) TableName() string {
	return "teacher_classroom"
}

type GradeAssignHistory struct {
	gorm.Model
	StudentSubmissionId uint    `gorm:"not null"`
	TeacherId           uint    `gorm:"not null"`
	AssignedGrade       float64 `gorm:"type:numeric(5,2)"`
	ChangeDate          time.Time
}

func (GradeAssignHistory) TableName() string {
	return "grade_assign_history"
}

type DocumentTemplate struct {
	gorm.Model
	Name          string
	Url           string
	TemplateType  string
	Classroom     []Classroom     `gorm:"foreignKey:CertificateTemplateId"`
	Qualification []Qualification `gorm:"foreignKey:CertificateTemplateId"`
}

func (DocumentTemplate) TableName() string {
	return "document_template"
}

type Users struct {
	gorm.Model
	FullName            string
	NickName            string
	Email               string
	Phone               string
	Birth               time.Time
	Gender              string
	Address             string
	City                string
	Country             string
	Nationality         string
	VerificationStatus  string
	IsVerified          bool
	TeacherStatus       string
	IsTeacher           bool
	EmployeeId          string
	IdCardNumber        string
	IsHrManager         bool
	IsBanned            bool
	IsAdmin             bool
	TeacherClassroom    []TeacherClassroom    `gorm:"foreignKey:UserId"`
	Classroom           []Classroom           `gorm:"foreignKey:CreatedBy"`
	StudentClassroom    []StudentClassroom    `gorm:"foreignKey:UserId"`
	VoteNew             []VoteNew             `gorm:"foreignKey:UserId"`
	VoteExisting        []VoteExisting        `gorm:"foreignKey:UserId"`
	SuggestedClassroom  []SuggestedClassroom  `gorm:"foreignKey:UserSuggestedId"`
	JobVacancy          []JobVacancy          `gorm:"foreignKey:CreatedBy"`
	UserApplication     []UserApplication     `gorm:"foreignKey:UserId"`
	EmploymentHistory   []EmploymentHistory   `gorm:"foreignKey:UserId"`
	EducationHistory    []EducationHistory    `gorm:"foreignKey:UserId"`
	Skill               []Skill               `gorm:"foreignKey:UserId"`
	OpenCandidate       []OpenCandidate       `gorm:"foreignKey:UserId"`
	ExternalCertificate []ExternalCertificate `gorm:"foreignKey:UserId"`
	News                []News                `gorm:"foreignKey:IssuedBy"`
	Testimony           []Testimony           `gorm:"foreignKey:UserId"`
}

func (Users) TableName() string {
	return "users"
}

type VoteNew struct {
	gorm.Model
	UserId               uint `gorm:"not null"`
	SuggestedClassroomId uint `gorm:"not null"`
	VoteType             string
}

func (VoteNew) TableName() string {
	return "vote_new"
}

type VoteExisting struct {
	gorm.Model
	UserId      uint `gorm:"not null"`
	ClassroomId uint `gorm:"not null"`
	VoteType    string
}

func (VoteExisting) TableName() string {
	return "vote_existing"
}

type SuggestedClassroom struct {
	gorm.Model
	Name            string
	Description     string
	UserSuggestedId uint `gorm:"not null"`
	Status          string
	VoteNew         []VoteNew `gorm:"foreignKey:SuggestedClassroomId"`
}

func (SuggestedClassroom) TableName() string {
	return "suggested_classroom"
}

type Job struct {
	gorm.Model
	Name             string
	Level            string
	Department       string
	JobQualification []JobQualification `gorm:"foreignKey:JobId"`
	JobVacancy       []JobVacancy       `gorm:"foreignKey:JobId"`
	OpenCandidate    []OpenCandidate    `gorm:"foreignKey:JobId"`
}

func (Job) TableName() string {
	return "job"
}

type JobQualification struct {
	gorm.Model
	JobId           uint `gorm:"not null"`
	QualificationId uint `gorm:"not null"`
}

func (JobQualification) TableName() string {
	return "job_qualification"
}

type JobVacancy struct {
	gorm.Model
	JobId             uint `gorm:"not null"`
	HotelId           string
	HideHotel         bool
	IssuedDate        time.Time
	ExpiredDate       string
	CreatedBy         uint `gorm:"not null"`
	Description       string
	Salary            int
	Status            string
	Type              string
	CandidateSelected string
	UserApplication   []UserApplication
}

func (JobVacancy) TableName() string {
	return "job_vacancy"
}

type UserApplication struct {
	gorm.Model
	UserId       uint `gorm:"not null"`
	JobVacancyId uint `gorm:"not null"`
	ApplyDate    time.Time
	Status       string
}

func (UserApplication) TableName() string {
	return "user_application"
}

type OpenCandidate struct {
	gorm.Model
	UserId           uint `gorm:"not null"`
	Start            time.Time
	UserIntroduction string
	IsOpen           bool
	JobId            uint           `gorm:"not null"`
	OpenLocation     []OpenLocation `gorm:"foreignKey:OpenCandidateId"`
}

func (OpenCandidate) TableName() string {
	return "open_candidate"
}

type OpenLocation struct {
	gorm.Model
	OpenCandidateId uint `gorm:"not null"`
	Location        string
}

func (OpenLocation) TableName() string {
	return "open_location"
}

type ExternalCertificate struct {
	gorm.Model
	Name        string
	FileUrl     string
	ExternalUrl string
	UserId      uint `gorm:"not null"`
	IssuedDate  time.Time
	ExpiredDate time.Time
	Type        string
}

func (ExternalCertificate) TableName() string {
	return "external_certificate"
}

type News struct {
	gorm.Model
	Title      string
	Content    string
	ImgUrl     string
	IssuedBy   uint `gorm:"not null"`
	IssuedDate time.Time
	Type       string
	NewsImage  []NewsImage `gorm:"foreignKey:NewsId"`
}

func (News) TableName() string {
	return "news"
}

type Image struct {
	gorm.Model
	Title     string
	Url       string
	NewsImage []NewsImage `gorm:"foreignKey:ImageId"`
}

func (Image) TableName() string {
	return "image"
}

type NewsImage struct {
	gorm.Model
	NewsId  uint `gorm:"not null"`
	ImageId uint `gorm:"not null"`
}

func (NewsImage) TableName() string {
	return "news_image"
}

type Testimony struct {
	gorm.Model
	Title      string
	Content    string
	UserId     uint `gorm:"not null"`
	PostedDate time.Time
}

func (Testimony) TableName() string {
	return "testimony"
}

type CourseTag struct {
	gorm.Model
	Name               string
	ClassroomCourseTag []ClassroomCourseTag `gorm:"foreignKey:CourseTagId"`
}

func (CourseTag) TableName() string {
	return "course_tag"
}

type ClassroomCourseTag struct {
	gorm.Model
	CourseTagId uint `gorm:"not null"`
	ClassroomId uint `gorm:"not null"`
}

func (ClassroomCourseTag) TableName() string {
	return "classroom_course_tag"
}

//////////////////////////////////////////////
/// MODEL MIGRATION
//////////////////////////////////////////////
func InitialMigration() {
	// dsn := "brian:1q2w3e4r!Q@W#E$R@tcp(127.0.0.1:3306)/classroom?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "host=localhost user=brian password=!Q@W#E$R1q2w3e4r dbname=classroom port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Faild to connect to database")
	}

	// db.Migrator().CreateTable(&Classroom{})
	db.AutoMigrate(&Classroom{})
	db.AutoMigrate(&DocumentTemplate{})
	// db.Migrator().CreateConstraint(&ClassroomPeriod{}, "Classroom")
	// db.Migrator().CreateConstraint(&ClassroomPeriod{}, "fk_classroom_period_classroom")
	db.AutoMigrate(&QualificationClassroom{})
	// db.Migrator().CreateConstraint(&Classroom{}, "QualificationClassroom")
	// db.Migrator().CreateConstraint(&Classroom{}, "fk_classroom_qualification_classroom")
	db.AutoMigrate(&Qualification{})
	// db.Migrator().CreateConstraint(&DocumentTemplate{}, "Qualification")
	// db.Migrator().CreateConstraint(&DocumentTemplate{}, "fk_document_template_qualification")
	// db.Migrator().CreateConstraint(&Qualification{}, "QualificationClassroom")
	// db.Migrator().CreateConstraint(&Qualification{}, "fk_qualification_qualification_classroom")
	db.AutoMigrate(&Users{})
	// db.Migrator().CreateConstraint(&Users{}, "TeacherClassroom")
	// db.Migrator().CreateConstraint(&Users{}, "fk_user_teacher_classroom")
	// db.Migrator().CreateConstraint(&Users{}, "Classroom")
	// db.Migrator().CreateConstraint(&Users{}, "fk_user_classroom")
	// db.Migrator().CreateConstraint(&Users{}, "StudentClassroom")
	// db.Migrator().CreateConstraint(&Users{}, "fk_user_student_classroom")
	// db.Migrator().CreateConstraint(&Users{}, "VoteNew")
	// db.Migrator().CreateConstraint(&Users{}, "fk_user_vote_new")
	// db.Migrator().CreateConstraint(&Users{}, "VoteExisting")
	// db.Migrator().CreateConstraint(&Users{}, "fk_user_vote_existing")
	db.AutoMigrate(&JobQualification{})
	// db.Migrator().CreateConstraint(&Qualification{}, "JobQualification")
	// db.Migrator().CreateConstraint(&Qualification{}, "fk_qualification_job_qualification")
	db.AutoMigrate(&Job{})
	// db.Migrator().CreateConstraint(&Job{}, "JobQualification")
	// db.Migrator().CreateConstraint(&Job{}, "fk_job_job_qualification")
	// db.Migrator().CreateConstraint(&Job{}, "JobVacancy")
	// db.Migrator().CreateConstraint(&Job{}, "fk_job_job_vacancy")
	db.AutoMigrate(&JobVacancy{})
	// db.Migrator().CreateConstraint(&Job{}, "JobVacancy")
	// db.Migrator().CreateConstraint(&Job{}, "fk_job_job_vacancy")
	db.AutoMigrate(&UserApplication{})
	// db.Migrator().CreateConstraint(&Users{}, "UserApplication")
	// db.Migrator().CreateConstraint(&Users{}, "fk_user_user_application")
	db.AutoMigrate(&EmploymentHistory{})
	// db.Migrator().CreateConstraint(&Users{}, "EmploymentHistory")
	// db.Migrator().CreateConstraint(&Users{}, "fk_user_employment_history")
	db.AutoMigrate(&EducationHistory{})
	// db.Migrator().CreateConstraint(&Users{}, "EducationHistory")
	// db.Migrator().CreateConstraint(&Users{}, "fk_user_education_history")
	db.AutoMigrate(&Skill{})
	// db.Migrator().CreateConstraint(&Users{}, "Skill")
	// db.Migrator().CreateConstraint(&Users{}, "fk_user_skill")
	db.AutoMigrate(&Topic{})
	db.AutoMigrate(&Department{})
	// db.Migrator().CreateConstraint(&Department{}, "Classroom")
	// db.Migrator().CreateConstraint(&Department{}, "fk_department_classroom")
	db.AutoMigrate(&ClassroomPeriod{})
	// db.Migrator().CreateConstraint(&Classroom{}, "ClassroomPeriod")
	// db.Migrator().CreateConstraint(&Classroom{}, "fk_classroom_classroom_period")
	db.AutoMigrate(&RegistrationPeriod{})
	// db.Migrator().CreateConstraint(&ClassroomPeriod{}, "RegistrationPeriod")
	// db.Migrator().CreateConstraint(&ClassroomPeriod{}, "fk_classroom_period_registration_period")
	db.AutoMigrate(&Assignment{})
	// db.Migrator().CreateConstraint(&Classroom{}, "Assignment")
	// db.Migrator().CreateConstraint(&Classroom{}, "fk_classroom_assignment")
	db.AutoMigrate(&StudentSubmission{})
	// db.Migrator().CreateConstraint(&Assignment{}, "StudentSubmission")
	// db.Migrator().CreateConstraint(&Assignment{}, "fk_assignment__student_submission")
	db.AutoMigrate(&StudentClassroom{})
	// db.Migrator().CreateConstraint(&StudentClassroom{}, "StudentSubmission")
	// db.Migrator().CreateConstraint(&StudentClassroom{}, "fk_student_classroom__student_submission")
	// db.Migrator().CreateConstraint(&ClassroomPeriod{}, "StudentClassroom")
	// db.Migrator().CreateConstraint(&ClassroomPeriod{}, "fk_classroom_period_student_classroom")
	db.AutoMigrate(&TeacherClassroom{})
	// db.Migrator().CreateConstraint(&Classroom{}, "TeacherClassroom")
	// db.Migrator().CreateConstraint(&Classroom{}, "fk_classroom_teacher_classroom")
	db.AutoMigrate(&SuggestedClassroom{})
	db.AutoMigrate(&VoteNew{})
	// db.Migrator().CreateConstraint(&SuggestedClassroom{}, "VoteNew")
	// db.Migrator().CreateConstraint(&SuggestedClassroom{}, "fk_suggested_classroom_vote_new")
	db.AutoMigrate(&VoteExisting{})
	// db.Migrator().CreateConstraint(&Classroom{}, "VoteExisting")
	// db.Migrator().CreateConstraint(&Classroom{}, "fk_classroom_vote_existing")
	db.AutoMigrate(&GradeAssignHistory{})
	// db.Migrator().CreateConstraint(&StudentSubmission{}, "GradeAssignHistory")
	// db.Migrator().CreateConstraint(&StudentSubmission{}, "fk_student_submission_grade_assign_history")
	// db.Migrator().CreateConstraint(&TeacherClassroom{}, "GradeAssignHistory")
	// db.Migrator().CreateConstraint(&TeacherClassroom{}, "fk_teacher_classroom_grade_assign_history")
	db.AutoMigrate(&OpenCandidate{})
	// db.Migrator().CreateConstraint(&Users{}, "OpenCandidate")
	// db.Migrator().CreateConstraint(&Users{}, "fk_user_open_candidate")
	// db.Migrator().CreateConstraint(&Job{}, "OpenCandidate")
	// db.Migrator().CreateConstraint(&Job{}, "fk_job_open_candidate")
	db.AutoMigrate(&OpenLocation{})
	// db.Migrator().CreateConstraint(&OpenCandidate{}, "OpenLocation")
	// db.Migrator().CreateConstraint(&OpenCandidate{}, "fk_open_candidate_open_location")
	db.AutoMigrate(&ExternalCertificate{})
	db.AutoMigrate(&News{})
	db.AutoMigrate(&Image{})
	db.AutoMigrate(&NewsImage{})
	db.AutoMigrate(&Testimony{})
	db.AutoMigrate(&CourseTag{})
	db.AutoMigrate(&ClassroomCourseTag{})
	db.AutoMigrate(&Schedule{})
}
