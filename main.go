package main

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/adityanath-madishetti/todo/backend/DB"
	controller "github.com/adityanath-madishetti/todo/backend/controllers"
	"github.com/adityanath-madishetti/todo/backend/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func init() {
 if err := godotenv.Load(); err != nil {
    log.Println("⚠️  No .env file found. Assuming env variables are set by host.")
}

}

func main(){
		fmt.Println("Welcome to MongoDB-based API")
		db.MakeDbConnection()




		r:=mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(controller.NotFoundHandler)
		sr2:=r.PathPrefix("/internals").Subrouter()
		sr:=r.PathPrefix("/api").Subrouter()
		routes.AuthRoutes(sr)
		routes.TaskRoutes(sr)
		routes.Userroutes(sr)
		routes.Internaluserroutes(sr2)
		

				handler := cors.New(cors.Options{
    AllowedOrigins: []string{"http://localhost:3000"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
}).Handler(r)
		

    log.Println("Server running on :8080")
	http.ListenAndServe(":8080", handler)


}