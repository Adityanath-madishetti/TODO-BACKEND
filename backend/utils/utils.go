package utils

import (
	"context"
	"fmt"

	db "github.com/adityanath-madishetti/todo/backend/DB"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)


var allowedFilters = map[string]bool{
	"completed":   true,
	"category": true,
	"priority": true,
	"userId":   true,
	"taskId":true,
}



func CleanFilter(raw bson.M) bson.M {
	clean := bson.M{}
	for key, value := range raw {
		if allowedFilters[key] {
			clean[key] = value

		} else {
			fmt.Printf("Warning: Ignored invalid filter key: %s\n", key)
		}
	}
	return clean
}



func EnsureUserNameUniqueIndex() error {
    indexModel := mongo.IndexModel{
        Keys: bson.M{"name": 1}, // create index on 'name' field
        Options: options.Index().SetUnique(true),
    }

    _, err := db.UserCollection.Indexes().CreateOne(context.Background(), indexModel)
    if err != nil {
        return fmt.Errorf("failed to create unique index on name: %w", err)
    }

    return nil
}


func IsUserNameTaken(name string) (bool, error) {
    count, err := db.UserCollection.CountDocuments(context.Background(), bson.M{"name": name})
    if err != nil {
        return false, err
    }
    return count > 0, nil
}
