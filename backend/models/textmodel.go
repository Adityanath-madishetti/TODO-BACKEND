package model

import (
	"context"
	"errors"
	"fmt"

	db "github.com/adityanath-madishetti/todo/backend/DB"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)


type TaskDescription struct{
	
	Id     bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` 
	UserId         string             `json:"userId" bson:"userId"`
	TaskId  		string 			  `json:"taskId" bson:"taskId"`
	Text   string `bson:"text" json:"text"`   
}



// to add description


func CreateDescription(newcollection TaskDescription) error{

	if newcollection.Id.IsZero(){
		newcollection.Id=bson.NewObjectID()
	}

	if newcollection.TaskId==""{
		return errors.New("TaskId is empty ")
	}

	if newcollection.UserId==""{
		return errors.New("UserID is empty")
	}

	res,err:=db.DescriptionCollection.InsertOne(context.Background(),newcollection)

	if err!=nil{
		return err
	}

	fmt.Println("the insertID of  result : ",res.InsertedID)

	return nil


}


func UpdateTaskDescription(taskId string,text string ) error{

	if taskId==""{
		return errors.New("TaskId is empty ")
	}

	res,err:=db.DescriptionCollection.UpdateOne(context.Background(),bson.M{
		"taskId":taskId,
	},bson.M{
		"$set":bson.M{
			"text":text,
		},
	})

	if err!=nil{
		return err;
	}

	fmt.Println("the update result : ",res.MatchedCount," ",res.ModifiedCount)

	return nil
}

func RemoveTaskDescription(taskId string) error{
	if taskId==""{
		return errors.New("TaskId is empty ")
	}

	dc,err:=db.DescriptionCollection.DeleteOne(context.Background(),bson.M{
		"taskId":taskId,
	})

	if err!=nil{
		return err
	}

	fmt.Println("the delete result : ",dc.DeletedCount)
	return nil
}



func GetTextDescription(taskId string) (string,error){
	if taskId==""{
		return "",errors.New("TaskId is empty ")
	}

	var found_taskdec TaskDescription

	sr:=db.DescriptionCollection.FindOne(context.Background(),bson.M{"taskId":taskId})

	if errors.Is(sr.Err(),mongo.ErrNoDocuments){
		return "", sr.Err()
	}

	if sr.Err()!=nil{
		return "",sr.Err()
	}


	if err:=sr.Decode(&found_taskdec);err!=nil{
		return "",err
	}



	return found_taskdec.Text,nil
}




