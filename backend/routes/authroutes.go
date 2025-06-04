package routes

import (
	

	controller "github.com/adityanath-madishetti/todo/backend/controllers"
	"github.com/gorilla/mux"
)




func AuthRoutes(r  *mux.Router){

	
	sr:=r.PathPrefix("/auth").Subrouter()

	sr.HandleFunc("/register",controller.SignUpController).Methods("POST")
	sr.HandleFunc("/login",controller.LoginController).Methods("POST")
	
}