package model

import "go.mongodb.org/mongo-driver/bson/primitive"

//Courses ...
type Courses struct {
	CourseId   int    `bson:"courseId" json:"courseId"`
	CourseName string `bson:"courseName" json:"courseName"`
	Year       int64  `bson:"year" json:"year"`
}

type Students struct {
	FirstName string    `bson:"firstName" json:"firstName"`
	LastName  string    `bson:"lastName" json:"lastName"`
	Age       int       `bson:"age" json:"age"`
	PhoneNo   int64     `bson:"phoneNo" json:"phoneNo"`
	Email     string    `bson:"email" json:"email"`
	RollNo    int64     `bson:"rollNo" json:"rollNo"`
	Course    []Courses `bson:"course" json:"course"`
}

type TestStudents struct {
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
	Age       int    `bson:"age" json:"age"`
	RollNo    int64  `bson:"rollNo" json:"rollNo"`
	//Course    []Courses `bson:"course" json:"course"`
}
type Simple struct {
	Name   string
	Keys   []primitive.E
	Unique bool
}

type SimpleIndexes struct {
	CollectionName string
	Indexes        []Simple
}
