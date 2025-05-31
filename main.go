package main

import (
	"fmt"

	db "github.com/adityanath-madishetti/todo/backend/DB"
)



func main(){
		fmt.Println("Welcome to MongoDB-based API")
		db.MakeDbConnection()

		fmt.Println("Starting server on :4000...")


}