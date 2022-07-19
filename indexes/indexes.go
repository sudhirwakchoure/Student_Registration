package indexes

import (
	"STUDENT_REGISTRATION/model"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexes(db *mongo.Database, tmp []model.SimpleIndexes, ctx context.Context) error {

	for _, val := range tmp {
		coll := db.Collection(val.CollectionName)
		for _, v1 := range val.Indexes {
			var mod []mongo.IndexModel
			Unique := v1.Unique
			mod = []mongo.IndexModel{
				{
					Keys: v1.Keys,
					Options: &options.IndexOptions{
						Unique: &Unique,
					},
				},
			}
			opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
			_, err := coll.Indexes().CreateMany(ctx, mod, opts)

			// Check for the options errors
			if err != nil {
				return err
			} else {
				fmt.Println("CreateMany() option:", opts)
			}
		}
	}
	return nil
}
