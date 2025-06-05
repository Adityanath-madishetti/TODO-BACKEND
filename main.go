package main

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/adityanath-madishetti/todo/backend/DB"
	"github.com/adityanath-madishetti/todo/backend/routes"
	"github.com/gorilla/mux"
)



func main(){
		fmt.Println("Welcome to MongoDB-based API")
		db.MakeDbConnection()

		r:=mux.NewRouter()
		sr2:=r.PathPrefix("/internals").Subrouter()
		sr:=r.PathPrefix("/api").Subrouter()
		routes.AuthRoutes(sr)
		routes.TaskRoutes(sr)
		routes.Userroutes(sr)
		routes.Internaluserroutes(sr2)
		
		

    log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)


}