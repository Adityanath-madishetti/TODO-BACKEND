package controller

import (
	"encoding/json"
	"net/http"

	model "github.com/adityanath-madishetti/todo/backend/models"
	"github.com/adityanath-madishetti/todo/backend/utils"
	"github.com/gorilla/mux"
)



//put based on id
func UpdateDescriptionController(w http.ResponseWriter,r *http.Request) {

	w.Header().Set("Content-Type","application/json")

	if r.Header.Get("Content-Type") != "application/json" {
        // http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		utils.SendJSONError(w,http.StatusUnsupportedMediaType, "Content-Type must be application/json")
        return
    }

	variables:=mux.Vars(r)

	taskId:=variables["id"]

	var reqbody map[string]string

	if err:=json.NewDecoder(r.Body).Decode(&reqbody);err!=nil{
		utils.SendJSONError(w,http.StatusBadRequest,"invalid json sent , text should be string")
		return
	}


	if taskId==""{
		utils.SendJSONError(w,http.StatusBadRequest,"taskId is empty")
		return
	}

	val ,exist:= reqbody["text"]

	if !exist{
		utils.SendJSONError(w,http.StatusBadRequest,"text parametere is missing in the body")
	}

	err:=model.UpdateTaskDescription(taskId,val)

	if err!=nil{
		utils.SendJSONError(w,http.StatusInternalServerError,"error in updating :"+err.Error())
		return
	}



	json.NewEncoder(w).Encode(map[string]string{"message":"succesfully updated"})
	


}


func GetTaskDescriptionController(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")

	// basically its get request bassed on id

	variables:=mux.Vars(r)
	taskId:=variables["id"]

	if taskId==""{
		utils.SendJSONError(w,http.StatusBadRequest,"taskId is empty")
		return
	}


	text,err:=model.GetTextDescription(taskId)

	if err!=nil{
		utils.SendJSONError(w,http.StatusInternalServerError,"error from textdesc.controller :"+err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"text":text,"message":"succesfully sent the text"})


}