package controller

import (
	"encoding/json"
	"net/http"

	"github.com/adityanath-madishetti/todo/backend/middleware"
	model "github.com/adityanath-madishetti/todo/backend/models"
	"github.com/adityanath-madishetti/todo/backend/utils"
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
func UpdateTaskController(w http.ResponseWriter, r *http.Response){



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


    if r.Header.Get("Content-Type") != "application/json" {
        // http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		utils.SendJSONError(w,http.StatusUnsupportedMediaType, "Content-Type must be application/json")
        return
    }

	var body map[string]interface{}
	if err:=json.NewDecoder(r.Body).Decode(&body);err!=nil{
		utils.SendJSONError(w,http.StatusBadRequest,"Invalid Json: "+err.Error())
		return
	}

	taskId, ok :=body["taskId"].(string)

	if(!ok){
		utils.SendJSONError(w,http.StatusUnprocessableEntity,"taskId should be string")
		return
	}

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


    if r.Header.Get("Content-Type") != "application/json" {
        // http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		utils.SendJSONError(w,http.StatusUnsupportedMediaType, "Content-Type must be application/json")
        return
    }



	
	
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
func GetTasksByCategory(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type", "application/json")


    if r.Header.Get("Content-Type") != "application/json" {
        // http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		utils.SendJSONError(w,http.StatusUnsupportedMediaType, "Content-Type must be application/json")
        return
    }


	var body map[string]interface{}

	if err:=json.NewDecoder(r.Body).Decode(&body);err!=nil{
		utils.SendJSONError(w,http.StatusBadRequest,"Invalid json;: "+err.Error())
		return
	}


	category, ok :=body["category"].(string)
	if(!ok){
		utils.SendJSONError(w,http.StatusUnprocessableEntity,"category should be string")
		return
	}


	if(category==""){
		utils.SendJSONError(w,http.StatusBadRequest,"category should be string")
		return
	}
	userId, ok :=r.Context().Value(middleware.ContextKeyUserID).(string)

	if(!ok){
		utils.SendJSONError(w,http.StatusUnprocessableEntity,"userId should be string")
		return
	}



	if(userId==""){
		utils.SendJSONError(w,http.StatusBadRequest,"userId should be string,some problem with string")
		return
	}





	var tasks []model.Task
	tasks,err:=model.GetTasksByCategoryForUser(category,userId)

	if err!=nil{
			utils.SendJSONError(w,http.StatusInternalServerError,"error fro, GetTaskByCategory : "+err.Error())
			return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"tasks":tasks,"message":"succesful"})

}



//add a task and remove a task controllers

