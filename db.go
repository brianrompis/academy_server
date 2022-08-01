package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}

func InitialMigration() {
	// dsn := "brian:1q2w3e4r!Q@W#E$R@tcp(127.0.0.1:3306)/classroom?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "host=127.0.0.1 user=brian password=!Q@W#E$R1q2w3e4r dbname=classroom port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Faild to connect to database")
	}

	db.AutoMigrate(&Schedule{})
	db.AutoMigrate(&Classroom{})
	db.AutoMigrate(&Topic{})
	db.AutoMigrate(&Certificate{})
	db.AutoMigrate(&Certification{})
	db.AutoMigrate(&EmploymentHistory{})
	db.AutoMigrate(&Student{})
	db.AutoMigrate(&Teacher{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&EducationHistory{})
	db.AutoMigrate(&Teacher{})
}
