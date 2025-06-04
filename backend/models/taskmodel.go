package model

// https://cloud.mongodb.com/v2/6839a019d5bf1b10435ecd6c#/explorer

import (
	"context"
	"fmt"

	"time"

	db "github.com/adityanath-madishetti/todo/backend/DB"
	"github.com/adityanath-madishetti/todo/backend/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)



type Task struct {
    Id             bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` 
    UserId         string             `json:"userId" bson:"userId"`
	TaskId  		string 			  `json:"taskId" bson:"taskId"`
    Category       string             `json:"category" bson:"category"`          
    Title          string             `json:"title" bson:"title"`
    Completed         bool               `json:"completed" bson:"completed"`
    CreationTime   time.Time          `json:"creationTime" bson:"creationTime"`
    CompletionTime *time.Time         `json:"endTime,omitempty" bson:"endTime,omitempty"`
	Priority 	int 				   `json:"priority" bson:"priority"`
}





// user id is given to you in task struct automatically so jsut handel times and id remaining is hoped to be done by controller
func CreateTask(newtask Task ) error{

	if newtask.Id.IsZero(){
		newtask.Id=bson.NewObjectID()
	}


	newtask.TaskId = uuid.New().String()
	newtask.CreationTime = time.Now()
	newtask.CompletionTime=nil
	newtask.Completed=false


	res,err:=db.Taskcollection.InsertOne(context.Background(),newtask)

	if err!=nil{
		return fmt.Errorf("error from AddTask %w",err)
	}

	fmt.Println("the task is added with insertid ",res.InsertedID)

	return nil
}




func Toggle(taskID string) error {
	ctx := context.Background()

	// 1. Find the task by UUID
	var foundTask Task
	filter := bson.M{"taskId": taskID}

	err := db.Taskcollection.FindOne(ctx, filter).Decode(&foundTask)
	if err != nil {
		return fmt.Errorf("error finding task: %w", err)
	}

	// 2. Prepare update
	var update bson.M
	if foundTask.Completed {
		// Mark as incomplete
		update = bson.M{
			"$set": bson.M{
				"completed":  false,
				"endTime": nil,
			},
		}
	} else {
		// Mark as complete
		now := time.Now()
		update = bson.M{
			"$set": bson.M{
				"completed":  true,
				"endTime": &now,
			},
		}
	}

	// 3. Apply update
	_, err = db.Taskcollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating task: %w", err)
	}

	return nil
}



func RemoveTask(taskID string) error{

	
	res,err:=db.Taskcollection.DeleteOne(context.Background(),bson.M{"taskId": taskID})	

	if(err!=nil){
		return fmt.Errorf("error from RemoveTask %w",err)
	}

	fmt.Println("the resultant delete count is ",res.DeletedCount)


	return nil
}


func ChangeTitle(taskID string, newtitle  string) error{

	update := bson.M{
		"$set":bson.M{
			"title":newtitle,
		},
	}

	filter:=bson.M{
		"taskId":taskID,
	}

	_, err :=db.Taskcollection.UpdateOne(context.Background(),filter,update)
	
		if(err!=nil){
			return fmt.Errorf("error in updating from changeTitle %w",err)
		}



		return nil

}




func ChangeCategory(taskID string, newcategory  string) error{

	update := bson.M{
		"$set":bson.M{
			"title":newcategory,
		},
	}

	filter:=bson.M{
		"taskId":taskID,
	}

	_, err :=db.Taskcollection.UpdateOne(context.Background(),filter,update)
	
		if(err!=nil){
			return fmt.Errorf("error in updating from changeTitle %w",err)

		}

		return nil

}



func ChangePriority(taskID string, newpriority int) error{
	update:= bson.M{
		"set":bson.M{
			"priority":newpriority,
		},
	}

	filter:=bson.M{
		"taskId":taskID,
	}

	_, err :=db.Taskcollection.UpdateOne(context.Background(),filter,update)

	if(err!=nil){
			return fmt.Errorf("error in updating from changeTitle %w",err)
	}

		return nil
}




func GetTaskById(taskID string) (Task,error) {

	var foundTask Task


	err:=db.Taskcollection.FindOne(context.Background(),bson.M{"taskId":taskID}).Decode(&foundTask)

	if(err!=nil){
		return foundTask,fmt.Errorf("error from GetTaskById %w",err)
	}
	return foundTask,nil


}


func GetAllTasksforUser(userID string)([]Task,error){
	

	curr,err:=db.Taskcollection.Find(context.Background(),bson.M{"userId":userID})

	var tasks []Task

	if(err!=nil){
		return tasks, fmt.Errorf("error from GetTasksByID  , %w",err)
	}

	defer curr.Close(context.Background())


	for curr.Next(context.Background()){
		var res Task

		if err:=curr.Decode(&res);err!=nil{
			return tasks, fmt.Errorf("error from GetTasksById , %w",err)
		}

		tasks = append(tasks, res)
	}

	return tasks,nil
}


func GetTasksByCategoryForUser( category string,userID string)([]Task,error){

	curr,err:=db.Taskcollection.Find(context.Background(),bson.M{"userId":userID,"category":category})

	var tasks []Task

	if(err!=nil){
		return tasks, fmt.Errorf("error from GetTasksByCategoryForUser  , %w",err)
	}

	defer curr.Close(context.Background())


	for curr.Next(context.Background()){
		var res Task

		if err:=curr.Decode(&res);err!=nil{
			return tasks, fmt.Errorf("error from GetTasksByCategoryForUser , %w",err)
		}

		tasks = append(tasks, res)
	}

	return tasks,nil



}


func GeneralFilter(filter bson.M)([]Task,error){

		var tasks []Task
		filter = utils.CleanFilter(filter)


	curr,err:=db.Taskcollection.Find(context.Background(),filter)
	

	if(err!=nil){
		return tasks, fmt.Errorf("error from generalFilter  , %w",err)
	}

	defer curr.Close(context.Background())	

	for curr.Next(context.Background()){
		var res Task

		if err:=curr.Decode(&res);err!=nil{
			return tasks, fmt.Errorf("error from generalFilter  cursor , %w",err)
		}

		tasks = append(tasks, res)
	}


	if err := curr.Err(); err != nil {
		return nil, fmt.Errorf("GeneralFilter: cursor error: %w", err)
	}

	return tasks,nil
}