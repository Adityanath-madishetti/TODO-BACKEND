package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"golang.org/x/crypto/bcrypt"
	model "github.com/adityanath-madishetti/todo/backend/models"
	"github.com/adityanath-madishetti/todo/backend/utils"
	"github.com/golang-jwt/jwt/v5"
)

// authentocation controller

//signup and login of the user

//just creates a user , it dosenot sign in the user

// POST /auth/register

type SignInRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}


// the body should only contain  username  and password as json fields

func SignUpController(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")


    if r.Header.Get("Content-Type") != "application/json" {
        utils.SendJSONError(w,http.StatusUnsupportedMediaType, "Content-Type must be application/json")
        return
    }

	// fmt.Println(r.Header.Get("Content-Type"))



    var req SignInRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        // http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		utils.SendJSONError(w,http.StatusBadRequest,"Invalid JSON body")
        return
    }

    if req.Username == "" || req.Password == "" {
        // http.Error(w, "Username and password are required", http.StatusBadRequest)
		utils.SendJSONError(w,http.StatusBadRequest,"Username and password are required")
        return
    }

    pass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        // http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		utils.SendJSONError(w,http.StatusInternalServerError,"Failed to hash password")
        return
    }

    err = model.CreateUser(model.User{
        Password: string(pass),
        Name:     req.Username,
    })


    if errors.Is(err, model.ErrUserExists) {
        // http.Error(w, "Username already exists choose any other username", http.StatusConflict)		
		utils.SendJSONError(w,http.StatusConflict,"Username already exists choose any other username")
        return
    }
    if err != nil {
        // http.Error(w, "Failed to create user", http.StatusInternalServerError)
		utils.SendJSONError(w,http.StatusInternalServerError,"Failed to create user")
        return
    }

    newuser, err := model.GetUserFromUsername(req.Username)
    if err != nil {
		fmt.Println("error in creating user",err)
        // http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		utils.SendJSONError(w,http.StatusInternalServerError,"Failed to fetch user")
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{"message":"succesful","userinfo":newuser})
}


//this method sends teh token after succesful login and this 

// actual logic of comparing teh passwor dis pending


func LoginController(w http.ResponseWriter,r *http.Request){

	w.Header().Set("Content-Type", "application/json")

	//first bring teh data from body



    if r.Header.Get("Content-Type") != "application/json" {
    	utils.SendJSONError(w,http.StatusUnsupportedMediaType, "Content-Type must be application/json")
        return
    }



    var req SignInRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        // http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		utils.SendJSONError(w,http.StatusBadRequest,"Invalid JSON body")
        return
    }


	if req.Username == "" || req.Password == "" {
        // http.Error(w, "Username and password are required", http.StatusBadRequest)
		utils.SendJSONError(w,http.StatusBadRequest,"Username and password are required")
        return
    }

	if flag, err := utils.IsUserNameTaken(req.Username); err != nil {

	utils.SendJSONError(w, http.StatusInternalServerError, "error in checking for unique user")

	return

	} else if !flag {
		// Handle username not present: signup first 
		utils.SendJSONError(w, http.StatusUnauthorized, "Bad credentials")
		return
	}


		user,_:=model.GetUserFromUsername(req.Username)	


	//check authenticity


	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.SendJSONError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}


	//update lastlogintime

	err:=model.UpdateLastLoginTime(req.Username)

	if(err!=nil){
		utils.SendJSONError(w,http.StatusInternalServerError,"error in updating login time ")
		return
	}


	

	// create a token and send it here 

	claims:=jwt.MapClaims{
		"exp": time.Now().Add(time.Hour*7).Unix(),
		"username":req.Username,
		"userid":user.UserID,
	}


	tokenGenerated:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	tokenString,err:=tokenGenerated.SignedString([]byte("Aditya@5002"))

	if(err!=nil){
		//handel http error 
		utils.SendJSONError(w,http.StatusInternalServerError,"error in generating signed string")
		return
	}

	// send token  and oprtionally any message 


	json.NewEncoder(w).Encode(map[string]string{
		"token":tokenString,
		"message":"successful",
	})


}


