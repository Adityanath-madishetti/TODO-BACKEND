package routes

import (
	controller "github.com/adityanath-madishetti/todo/backend/controllers"
	"github.com/adityanath-madishetti/todo/backend/middleware"
	"github.com/gorilla/mux"
)

func Userroutes(r *mux.Router) {
	
	sr:=r.PathPrefix("/user").Subrouter()
	
	sr.Use(middleware.AuthenticationMiddleware)
	sr.HandleFunc("/passwordchange",controller.ChangePassword).Methods("POST")
	

}