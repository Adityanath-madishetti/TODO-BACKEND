package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	db "github.com/adityanath-madishetti/todo/backend/DB"
	"github.com/adityanath-madishetti/todo/backend/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)


type User struct {
	ID bson.ObjectID 		`bson:"_id,omitempty" json:"id,omitempty"`
    UserID       string     `bson:"userId" json:"userId"`
    Password     string     `bson:"password" json:"-"` 
    Name         string     `bson:"name" json:"name"`
    CreationTime time.Time  `bson:"creationTime" json:"creationTime"`
    LastLogin    *time.Time `bson:"lastLogin,omitempty" json:"lastLogin,omitempty"` 
}



		// this methods just creates user with given deatils thats it!!



// ----------------------------------	VARS 	-----------------------------------

var ErrUserExists = errors.New("user already exists")

// --------------------------------------------------------------------------------


func CreateUser(newuser User) error {


		if newuser.ID.IsZero(){
			newuser.ID=bson.NewObjectID()
		}

	val, _ := utils.IsUserNameTaken(newuser.Name)

	if val {
		return errors.New("username already taken")
	}


		newuser.UserID=uuid.New().String()
		newuser.CreationTime = time.Now()
		newuser.LastLogin = nil

		res,err:=db.UserCollection.InsertOne(context.Background(),newuser)

		if(err!=nil){
			return fmt.Errorf("error from Createusermodel: %w ",err)
		}

		fmt.Println("the object is inserted in database with insertId as ",res.InsertedID)

	return nil;
}
















