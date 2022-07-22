package controller

import (
	"STUDENT_REGISTRATION/model"
	"STUDENT_REGISTRATION/utility"
	"context"
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

	FindAllCources(c)

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
		params = append(params, primitive.M{"courseId": courseId})
		filter = primitive.M{"courseId": courseId}
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
		var course model.Courses
		err := cur.Decode(&course)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
			log.Print(err)
			return
		}

		// collection.Find().sort({cocourse.courseId:-1}).limit(1)

		Allcources = append(Allcources, course)
	}

	c.JSON(http.StatusFound, Allcources)

}

func DeleteCourse(c *gin.Context) {
	id := c.Param("id")
	log.Println(id)
	var ctx context.Context

	collection, _ := utility.DB()

	filter := primitive.M{"courseId": id}
	log.Println(filter)

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

	id := c.Param("id")

	var course model.Courses
	filter := primitive.M{"courseId": id}

	var update map[string]string

	err := c.BindJSON(&update)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		log.Print(err)
		return
	}
	err = collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&course)

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
	No := c.Query("rollNo")
	rollNo, err := strconv.Atoi(No)
	if err != nil {
		log.Println("Error parsing", err)
		c.JSON(http.StatusInternalServerError, " Error parsing please provide valid if no ")
		return
	}

	log.Println("@@@@@@@@@@@@@@@@@@@@@", rollNo)
	params := []primitive.M{}

	filter := primitive.M{}

	if firstName != "" {
		params = append(params, primitive.M{"firstName": firstName})
		filter = primitive.M{"firstName": firstName}
	}
	if rollNo != 0 {

		params = append(params, primitive.M{"rollNo": rollNo})
		filter = primitive.M{"rollNo": rollNo}
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

// func Aggregate(c *gin.Context) {
// 	course := c.Query("course")

// 	collection := utility.DB1()

// 	pipeline := []primitive.M{}

// 	match := primitive.M{
// 		"$match": primitive.M{
// 			"course.courseName": course,
// 		},
// 	}
// 	// unwind := primitive.M{
// 	// 	"$unwind": primitive.M{
// 	// 		"path": "$course",
// 	// 	},
// 	// }

// 	pipeline = append(pipeline, match)

// 	printPipeline(pipeline)

// 	result, err := collection.Aggregate(c, pipeline)
// 	if err != nil {
// 		log.Println("Error aggregating pipeline: ", err)
// 		return

// 	}

// 	allresult := []model.TestStudents{}
// 	for result.Next(c) {
// 		student := model.TestStudents{}
// 		err = result.Decode(&student)

// 		// log.Printf("result :%+v", testgroup)
// 		if err != nil {
// 			log.Println("Error decoding aggregation result: ", err)
// 			return

// 		}
// 		log.Printf("\nresult test group:%v\n", student)
// 		allresult = append(allresult, student)

// 	}
// 	c.JSON(http.StatusOK, allresult)

// }
// func printPipeline(pipeline []primitive.M) {
// 	jpipeline, _ := json.Marshal(pipeline)
// 	log.Printf("\npipeline: %v\n", string(jpipeline))
// }
