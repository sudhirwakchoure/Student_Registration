package controller

import (
	"STUDENT_REGISTRATION/model"
	"STUDENT_REGISTRATION/utility"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func Homepage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "welcome to the home page of student registration service",
	})
}

func Addcoures(c *gin.Context) {
	var cources, _ model.Courses
	collection, _ := utility.DB()
	var ctx context.Context
	err := c.ShouldBind(&cources)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		log.Print(err)
		return
	}
	_, err = collection.InsertOne(ctx, cources)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": "ALready Exists "})
		log.Print(err)
		return
	}
	c.JSON(http.StatusCreated, cources)
}

func FindAllCources(c *gin.Context) {

	var Allcources []model.Courses
	collection, _ := utility.DB()
	var ctx context.Context

	courseName := c.Query("courseName")
	courseId := c.Query("courseId")

	log.Printf("\ncourse name:%v\n", courseName)
	params := []primitive.M{}

	filter := primitive.M{}

	if courseName != "" {
		params = append(params, primitive.M{"courseName": courseName})
		filter = primitive.M{"courseName": courseName}
	}
	if courseId != "" {
		id, err := strconv.Atoi(courseId)
		if err != nil {
			err = errors.New("failed to convert course id to integer")
			c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
			log.Print(err)
			return
		}
		params = append(params, primitive.M{"courseId": id})
		filter = primitive.M{"courseId": id}
	}

	if len(params) > 1 {
		filter = primitive.M{"$and": params}
	}

	count, err := collection.CountDocuments(ctx, filter)
	log.Println("error find", err)
	if err != nil || count == 0 {
		err = errors.New("no  such a course found")
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		return
	}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		log.Print(err)
		return
	}
	for cur.Next(ctx) {
		var course model.Courses
		err := cur.Decode(&course)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
			log.Print(err)
			return
		}
		Allcources = append(Allcources, course)
	}

	c.JSON(http.StatusFound, Allcources)

}

func DeleteCourse(c *gin.Context) {
	courseId := c.Param("id")

	log.Println(courseId)
	var ctx context.Context

	collection, _ := utility.DB()

	id, err := strconv.Atoi(courseId)
	if err != nil {
		err = errors.New("failed to convert course id to integer")
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		log.Print(err)
		return
	}

	filter := primitive.M{"courseId": id}
	log.Println(filter)

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil || count == 0 {
		err = errors.New("no documents is already deleted")
		c.JSON(http.StatusGone, gin.H{"Response": err.Error()})
		return
	}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusGone, gin.H{"Response": err.Error()})
		log.Print(err)
		return

		//log.Fatal(err)
	}
	c.JSON(http.StatusGone, result)

}

func UpdateCourse(c *gin.Context) {
	collection, _ := utility.DB()
	var ctx context.Context

	courseId := c.Param("id")
	id, err := strconv.Atoi(courseId)
	if err != nil {
		err = errors.New("failed to convert course id to integer")
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		log.Print(err)
		return
	}

	var course model.Courses
	filter := primitive.M{"courseId": id}

	err = c.BindJSON(&course)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		log.Print(err)
		return
	}
	err = collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": course}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&course)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"Response": err.Error()})
		log.Print(err)
		fmt.Println(err)
		return

	}
	log.Println(course)

	c.JSON(http.StatusCreated, course)
}

func CreateStudent(c *gin.Context) {
	courses := []model.Courses{}
	var student, _ model.Students
	collection := utility.DB1()
	coll, _ := utility.DB()
	var ctx context.Context
	err := c.ShouldBind(&student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		log.Print(err)
		return
	}

	cur, err := coll.Find(ctx, primitive.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		log.Print(err)
		return
	}
	for cur.Next(ctx) {
		var course model.Courses
		err := cur.Decode(&course)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
			log.Print(err)
			return
		}
		courses = append(courses, course)
	}
	for _, courcess := range student.Course {
		if !StringInSlice(courcess, courses) {
			log.Printf("Course Not found:%+v", courcess.CourseName)
			c.JSON(http.StatusNotFound, "please check courses list")
			return
		}
	}

	_, err = collection.InsertOne(ctx, student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": "ALready Exists"})
		log.Print(err)
		return
	}

	c.JSON(http.StatusCreated, student)

}

func StringInSlice(str model.Courses, list []model.Courses) bool {
	for _, b := range list {
		if str == b {
			return true
		}
	}
	return false
}

func Getstudent(c *gin.Context) {

	var Allstudent []model.Students
	collection := utility.DB1()
	var ctx context.Context

	firstName := c.Query("firstName")

	course := c.Query("course")
	rollNo := c.Query("rollNo")

	log.Println("@@@@@@@@@@@@@@@@@@@@@", rollNo)
	params := []primitive.M{}

	filter := primitive.M{}

	if firstName != "" {
		params = append(params, primitive.M{"firstName": firstName})
		filter = primitive.M{"firstName": firstName}
	}
	if rollNo != "" {
		no, err := strconv.Atoi(rollNo)
		if err != nil {
			log.Println("Error parsing", err)
			c.JSON(http.StatusInternalServerError, " Error parsing please provide valid if no ")
			return
		}
		params = append(params, primitive.M{"rollNo": no})
		filter = primitive.M{"rollNo": no}
	}
	if course != "" {
		params = append(params, primitive.M{"course.courseName": course})
		filter = primitive.M{"course.courseName": course}
	}

	if len(params) > 1 {
		filter = primitive.M{"$and": params}
	}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		log.Print(err)
		return
	}
	for cur.Next(ctx) {
		var student model.Students
		err := cur.Decode(&student)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
			log.Print(err)
			return
		}

		Allstudent = append(Allstudent, student)
	}

	c.JSON(http.StatusFound, Allstudent)
}

func DeleteStudent(c *gin.Context) {
	No := c.Param("rollNo")

	rollNo, err := strconv.Atoi(No)
	if err != nil {
		log.Println("Error parsing", err)
		c.JSON(http.StatusInternalServerError, " Error parsing please provide valid  no ")
		return
	}
	log.Println(rollNo)
	var ctx context.Context

	collection := utility.DB1()

	filter := primitive.M{"rollNo": rollNo}
	log.Println(filter)

	count, err := collection.CountDocuments(ctx, filter)
	log.Println("error find", err)
	if err != nil || count == 0 {
		err = errors.New("no  such a student found")
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		return
	}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusGone, gin.H{"Response": err.Error()})
		log.Print(err)
		return

		//log.Fatal(err)
	}
	c.JSON(http.StatusGone, result)

}

func UpdateStudent(c *gin.Context) {
	collection := utility.DB1()
	var ctx context.Context
	coll, _ := utility.DB()
	courses := []model.Courses{}

	No := c.Param("rollNo")
	rollNo, err := strconv.Atoi(No)
	if err != nil {
		log.Println("Error parsing", err)
		c.JSON(http.StatusInternalServerError, " Error parsing please provide valid  no ")
		return
	}

	var student model.Students
	filter := primitive.M{"rollNo": rollNo}

	err = c.BindJSON(&student)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		log.Print(err)
		return
	}

	cur, err := coll.Find(ctx, primitive.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		log.Print(err)
		return
	}
	for cur.Next(ctx) {
		var course model.Courses
		err := cur.Decode(&course)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
			log.Print(err)
			return
		}
		courses = append(courses, course)
	}
	for _, courcess := range student.Course {
		if !StringInSlice(courcess, courses) {
			log.Printf("Course Not found:%+v", courcess.CourseName)
			c.JSON(http.StatusNotFound, "please check courses list")
			return
		}
	}

	err = collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": student}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&student)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"Response": err.Error()})
		log.Print(err)
		fmt.Println(err)
		return

	}
	log.Println(student)

	c.JSON(http.StatusCreated, student)
}
