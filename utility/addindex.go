package utility

import (
	"STUDENT_REGISTRATION/indexes"
	"STUDENT_REGISTRATION/model"
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Addindex() {
	col1, db := DB()
	studentDB := DB1()
	ctx := context.Background()
	var tmp = []model.SimpleIndexes{
		{
			CollectionName: col1.Name(),
			Indexes: []model.Simple{
			
				{
					Name: "CourseCollection1",
					Keys: []primitive.E{
						{Key: "courseId", Value: 1},
						{Key: "courseName", Value: 1},
						{Key: "year", Value: 1},
					},
					Unique: true, 
				},
			},
		},
		{
			CollectionName: studentDB.Name(),
			Indexes: []model.Simple{

				{
					Name: "StudentCollection",
					Keys: []primitive.E{
						{Key: "rollNo", Value: 1},
						{Key: "phoneNo", Value: 1},
						{Key: "email", Value: 1},
					},
					Unique: true,
				},
				{
					Name: "StudentCollection1",
					Keys: []primitive.E{
						{Key: "course.courseId", Value: 1},
						{Key: "course.courseName", Value: 1},
						{Key: "course.year", Value: 1},
					},
					Unique: true,
				},
			},
		},
	}

	err := indexes.CreateIndexes(db, tmp, ctx)
	if err != nil {
		fmt.Println("Indexes().CreateMany() ERROR:", err)
		os.Exit(1) // exit in case of error
	}

}
