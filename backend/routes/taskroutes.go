package routes

import (
	controller "github.com/adityanath-madishetti/todo/backend/controllers"
	"github.com/adityanath-madishetti/todo/backend/middleware"
	"github.com/gorilla/mux"
)



func TaskRoutes(r *mux.Router) {

	
	
	sr:=r.PathPrefix("/tasks").Subrouter() 
	sr.Use(middleware.AuthenticationMiddleware)
	sr.HandleFunc("/",controller.GetTasksForUser).Methods("GET") // get all tasks 
	sr.HandleFunc("/",controller.AddTaskcontroller).Methods("POST") // addtask
	sr.HandleFunc("/",controller.UpdateTaskController) .Methods("PUT")			// update the task
	sr.HandleFunc("/{id}",controller.RemoveController).Methods("DELETE") //removetask
	sr.HandleFunc("/{id}",controller.GetTaskFromId).Methods("GET") // get task by id
	sr.HandleFunc("/category/{category}",controller.GetTasksByCategory).Methods("GET") //get task by category


}