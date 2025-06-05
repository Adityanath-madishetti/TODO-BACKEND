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
	"go.mongodb.org/mongo-driver/v2/mongo"
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
var ErrNouserExists=errors.New("no new user with this detaisl exist")

// --------------------------------------------------------------------------------


func CreateUser(newuser User) error {


		if newuser.ID.IsZero(){
			newuser.ID=bson.NewObjectID()
		}

	val, _ := utils.IsUserNameTaken(newuser.Name)

	if val {
		return fmt.Errorf("CreateUser: %w", ErrUserExists)
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




// not needed much cause its checked directly afetr fetching user
func CheckUserPresence(hashedPassword string, username string) (bool, error) {
	count, err := db.UserCollection.CountDocuments(context.Background(), bson.M{
		"name":     username,
		"password": hashedPassword,
	})

	if(errors.Is(err,mongo.ErrNoDocuments)){
		return false, fmt.Errorf("CheckUserPresence %w",ErrNouserExists)
	}
	if err != nil {
		return false, fmt.Errorf("CheckUserPresence: %w", err)
	}
	return count == 1, nil
}

func GetUserObject(userID string) (User, error) {
	var user User

	err := db.UserCollection.FindOne(context.Background(), bson.M{"userId": userID}).Decode(&user)
	if err != nil {
		return user, fmt.Errorf("GetUserObject: %w", err)
	}

	return user, nil
}


func GetUserFromUsername(userName string) (User,error){
	var user User

	err := db.UserCollection.FindOne(context.Background(), bson.M{"name": userName}).Decode(&user)
	if err != nil {
		return user, fmt.Errorf("GetUserObject: %w", err)
	}

	return user, nil
}




func GetAllUsers() ([]User,error){

	var users []User

	curr,err:=db.UserCollection.Find(context.Background(),bson.M{})

	if err!=nil{
		return users,fmt.Errorf("error from Get all users: %w",err)
	}

	defer curr.Close(context.Background())

	for curr.Next(context.Background()){
			var user User

		if err:=curr.Decode(&user);err!=nil{
			return users, fmt.Errorf("error from Get all users: %w",err)
		}

		users=append(users, user)
	}

	return users,nil


}



func GetUserFromUserId(userid string) (User,error){
	var user User

	err := db.UserCollection.FindOne(context.Background(), bson.M{"userId": userid}).Decode(&user)
	if err != nil {
		return user, fmt.Errorf("GetUserObject: %w", err)
	}

	return user, nil
}



func UpdateLastLoginTime(username string)error{



	now := time.Now()
	res,err:=db.UserCollection.UpdateOne(context.Background(),bson.M{"name":username},bson.M{
		"$set":bson.M{
			"lastLogin":&now,
		},
	})

	if(err!=nil){
		return fmt.Errorf("error from UpdateLastLogin, %w",err)
	}

	fmt.Println("info regarding update (Mc, U)c is : ",res.MatchedCount,res.ModifiedCount)

	return nil

}




func UpdatePassword(userId string, newPasswordHash string) error {
	res, err := db.UserCollection.UpdateOne(
		context.Background(),
		bson.M{"userId": userId},
		bson.M{"$set": bson.M{"password": newPasswordHash}},
	)
	if err != nil {
		return fmt.Errorf("error from UpdatePassword model: %w", err)
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("no user found with userId %s", userId)
	}

	if res.ModifiedCount == 0 {
		fmt.Println("Password was already the same; nothing updated")
	} else {
		fmt.Println("Password updated successfully")
	}

	return nil
}





