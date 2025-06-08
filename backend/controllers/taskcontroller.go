package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adityanath-madishetti/todo/backend/middleware"
	model "github.com/adityanath-madishetti/todo/backend/models"
	"github.com/adityanath-madishetti/todo/backend/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// update sany thing

/*

expected body is

{
	"types": ["category","title","toggle"]
	updates:{
		"catogery":"work"
		"title":"project"
		"toggle":true
	}
}


*/

//PUT



type UpdateRequest struct {
    Types   []string               `json:"types"`
    Updates map[string]interface{} `json:"updates"`
	TaskId  string	`json:"taskId"`
}



// depends on taskId so userId dosent matter
func UpdateTaskController(w http.ResponseWriter, r *http.Request){



	w.Header().Set("Content-Type", "application/json")


    if r.Header.Get("Content-Type") != "application/json" {
        // http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		utils.SendJSONError(w,http.StatusUnsupportedMediaType, "Content-Type must be application/json")
        return
    }

	// expects the  body json
	var req UpdateRequest
	if err:=json.NewDecoder(r.Body).Decode(&req); err!=nil{
		utils.SendJSONError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}


	if req.TaskId == "" {
		utils.SendJSONError(w, http.StatusBadRequest, "TaskId is required")
		return
	}


	// validating input

	for _, v := range req.Types {
	switch v {
			case "category":
				if _, ok := req.Updates["category"].(string); !ok {
					utils.SendJSONError(w, http.StatusUnprocessableEntity, "Category must be a string")
					return
				}
			case "title":
				if _, ok := req.Updates["title"].(string); !ok {
					utils.SendJSONError(w, http.StatusUnprocessableEntity, "Title must be a string")
					return
				}
			case "priority":
				p, ok := req.Updates["priority"].(float64) // JSON numbers are float64
				if !ok || int(p) < 0 || int(p) > 2 {
					utils.SendJSONError(w, http.StatusUnprocessableEntity, "Priority must be an integer between 0 and 2")
					return
				}
			case "toggle":
				if _, ok := req.Updates["toggle"].(bool); !ok {
					utils.SendJSONError(w, http.StatusUnprocessableEntity, "Toggle must be a boolean")
					return
				}


				
			default:
				utils.SendJSONError(w, http.StatusBadRequest, "Unknown update type: "+v)
				return
		}
	}





	for _,v:=range req.Types{

		switch v {

			case "category":
				// check wether the type of value is string r not

				// if string then call model to chaneg and handel the returnd error appropriate;y
				newCategory,_:= req.Updates["category"].(string)

			
				if err:=model.ChangeCategory(req.TaskId,newCategory);err!=nil{
					utils.SendJSONError(w,http.StatusInternalServerError,"from Update controller : "+err.Error())
					return
				}


			case "title":
				newtitle,_:= req.Updates["title"].(string)

				if err:=model.ChangeTitle(req.TaskId,newtitle);err!=nil{
					utils.SendJSONError(w,http.StatusInternalServerError,"from Update controller : "+err.Error())
					return
				}

			case "priority":
				newPriority,_:=req.Updates["priority"].(float64)

				

				if err:=model.ChangePriority(req.TaskId,int(newPriority));err!=nil{
					utils.SendJSONError(w,http.StatusInternalServerError,"from update controller : "+err.Error())
					return
				}


			case "toggle":

				if err:=model.Toggle(req.TaskId);err!=nil{
					utils.SendJSONError(w,http.StatusInternalServerError,"from update controller : "+err.Error())
					return 
				}

			default:
				utils.SendJSONError(w, http.StatusBadRequest, "Unknown update type: " + v)
				return

			}
	}

	
	json.NewEncoder(w).Encode(map[string]interface{}{"message":"succesfuly updated all the given fields"})


}




//dosent use userId
func GetTaskFromId(w http.ResponseWriter , r *http.Request){



	w.Header().Set("Content-Type", "application/json")

	fmt.Println("endpoint got hit from id task controller")
	mapc:=mux.Vars(r)
	taskId:=mapc["id"]
	

	if(taskId==""){
		utils.SendJSONError(w,http.StatusBadRequest,"taskId should be string")
		return
	}

	var task model.Task

	task,err:=model.GetTaskById(taskId);
	if err!=nil{
		utils.SendJSONError(w,http.StatusInternalServerError,"error from GetTask controller : "+err.Error())
		return
	}
	

	json.NewEncoder(w).Encode(map[string]interface{}{"task":task,"message":"succesful"})
	
}



// in this function u  get userId from token

func GetTasksForUser(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type", "application/json")


	
	
	userId, ok:=r.Context().Value(middleware.ContextKeyUserID).(string)
	
	if(!ok){
		utils.SendJSONError(w,http.StatusUnprocessableEntity,"UserId should be string so problem with token")
		return
	}

	if(userId==""){
				utils.SendJSONError(w,http.StatusBadRequest,"userId should be string")
				return
	}

	var tasks []model.Task

	tasks,err:=model.GetAllTasksforUser(userId);
	if err!=nil{
		utils.SendJSONError(w,http.StatusInternalServerError,"error from GetTaskForUser controller : "+err.Error())
		return
	}
	

	json.NewEncoder(w).Encode(map[string]interface{}{"tasks":tasks,"message":"succesful"})

}





// this function uses userId from token
func GetTasksByCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the category from the URL path
	vars := mux.Vars(r)
	category, ok := vars["category"]
	if !ok || category == "" {
		utils.SendJSONError(w, http.StatusBadRequest, "Missing or empty category in URL")
		return
	}

	// Get userId from context (from middleware)
	userId, ok := r.Context().Value(middleware.ContextKeyUserID).(string)
	if !ok || userId == "" {
		utils.SendJSONError(w, http.StatusBadRequest, "Invalid or missing user ID")
		return
	}

	// Fetch tasks for user by category
	tasks, err := model.GetTasksByCategoryForUser(category, userId)
	if err != nil {
		utils.SendJSONError(w, http.StatusInternalServerError, "Error from GetTasksByCategory: "+err.Error())
		return
	}

	// Return tasks
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tasks":   tasks,
		"message": "successful",
	})
}




//add a task and remove a task controllers

func AddTaskcontroller(w http.ResponseWriter,r * http.Request){

	// just adds a task


	w.Header().Set("Content-Type", "application/json")


    if r.Header.Get("Content-Type") != "application/json" {
        // http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		utils.SendJSONError(w,http.StatusUnsupportedMediaType, "Content-Type must be application/json")
        return
    }

	userId,ok:=r.Context().Value(middleware.ContextKeyUserID).(string)


	if(!ok){
		utils.SendJSONError(w,http.StatusUnprocessableEntity,"UserId should be string so problem with token")
		return
	}

	if(userId==""){
		utils.SendJSONError(w,http.StatusBadRequest,"userId should be string so problem with token")
				return
	}

	var newtask model.Task

	//first decide how  request is sent

	type taskrequest struct{
		Category string `json:"category"`
		Title string 	`json:"title"`
		Priority int `json:"priority"`
	}

	var reqbody taskrequest
	if err:=json.NewDecoder(r.Body).Decode(&reqbody);err!=nil{
		utils.SendJSONError(w,http.StatusBadRequest,"Invalid json object given")
		return
	}

	newtask.Category=reqbody.Category
	newtask.Title=reqbody.Title
	newtask.Priority=reqbody.Priority
	newtask.UserId=userId
	newtask.Completed=false

	var taskIdgenerated string
	taskIdgenerated,err:=model.CreateTask(newtask);
	if err!=nil{
		utils.SendJSONError(w,http.StatusInternalServerError,"error from AddTaskcontroller: "+err.Error())
		return
	}

	var newcollection model.TaskDescription
	newcollection.TaskId=taskIdgenerated
	newcollection.Text=""
	newcollection.UserId=userId

	if err:=model.CreateDescription(newcollection);err!=nil{
		utils.SendJSONError(w,http.StatusInternalServerError,"error from createTask: "+err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"message":"succesfully done"})

}


func RemoveController(w http.ResponseWriter,r * http.Request){
	w.Header().Set("Content-Type", "application/json")


    

	vars:=mux.Vars(r)
	taskid:=vars["id"]


	if taskid == "" {
    utils.SendJSONError(w, http.StatusBadRequest, "Missing task ID in URL")
    return
	}

	if err:=model.RemoveTask(taskid);err!=nil{
		utils.SendJSONError(w,http.StatusInternalServerError,"error from RemoveTaskcontroller: "+err.Error())
		return
	}

	if err:=model.RemoveTaskDescription(taskid);err!=nil{
		utils.SendJSONError(w,http.StatusInternalServerError,"error from RemoveTaskController: "+err.Error())
		return

	}

	json.NewEncoder(w).Encode(map[string]interface{}{"message":"succesfully done"})


}


func GeneralFiltercontroller(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters (use r.URL.Query().Get("key"))



	// fmt.Println("from the general Filter")


	title := r.URL.Query().Get("title")
	priority := r.URL.Query().Get("priority")
	category := r.URL.Query().Get("category")
	completed:=r.URL.Query().Get("status")

	// fmt.Println("priority: ",priority)

	//add user id

	userId, ok := r.Context().Value(middleware.ContextKeyUserID).(string)
	if !ok || userId == "" {
		utils.SendJSONError(w, http.StatusBadRequest, "Invalid or missing user ID")
		return
	}

	




	// Create a dynamic BSON filter
	filter := bson.M{}

	filter["userId"]=userId
	if title != "" {
		filter["title"] = title
	}
	if priority != "" {
		valFloat, err := strconv.ParseFloat(priority, 64)
		if err != nil {
			utils.SendJSONError(w, http.StatusBadRequest, "priority must be a number (e.g., 1, 2, 3): "+err.Error())
			return
		}

		val := int(valFloat)
		if valFloat != float64(val) {
			utils.SendJSONError(w, http.StatusBadRequest, "priority must be a whole number without decimal (e.g., 1, 2, 3)")
			return
		}

		if val < 1 || val > 3 {
			utils.SendJSONError(w, http.StatusBadRequest, "priority must be between 1 and 3")
			return
		}

		// fmt.Println("converted priority val is", val)
		filter["priority"] = val
	}
	if category != "" {
		filter["category"] = category
	}

	if completed!=""{

		if completed=="true"{
			filter["completed"]=true
		}else{
			filter["completed"]=false
		}
	}

	tasks, err := model.GeneralFilter(filter)
	// fmt.Println("tasks being returned are ",tasks)
	if err != nil {
		utils.SendJSONError(w,  http.StatusInternalServerError,"Failed to fetch tasks")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message":"succesful","tasks":tasks})
}