package routes

import (
	controller "github.com/adityanath-madishetti/todo/backend/controllers"
	"github.com/gorilla/mux"
)

func Internaluserroutes(r *mux.Router){
	sr:=r.PathPrefix("/users").Subrouter()
	sr.HandleFunc("/",controller.GetAllUsersController)
}