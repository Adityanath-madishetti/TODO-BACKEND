package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/adityanath-madishetti/todo/backend/utils"
	"github.com/golang-jwt/jwt/v5"
)


type contextKey string

const(
	ContextKeyUserID   contextKey = "userid"
    ContextKeyUsername contextKey = "username"
)



func AuthenticationMiddleware(next http.Handler) http.Handler{

		return http.HandlerFunc(
			(
				func(w http.ResponseWriter, r *http.Request){


				//here i need to take the token from the header and now i need to check wether it is valid

				//if valid then extract teh payload and give it back

				authHeader:= strings.Split(r.Header.Get("Authorization"),"Bearer ") // it is now basically slice of strings

				
					if(len(authHeader)!=2){
						//handel error

						utils.SendJSONError(w,http.StatusUnauthorized,"Malformed Token")
						return
					}

				//handel validation now
				tokensent:=authHeader[1]

				actualToken,err:=	jwt.Parse(tokensent,func(token *jwt.Token)(interface{},error){
																		if token.Method != jwt.SigningMethodHS256 {
    																	return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
																		}
																	return []byte("Aditya@5002"),nil
				})

				if err != nil || !actualToken.Valid {
					utils.SendJSONError(w, http.StatusUnauthorized, "Invalid token: "+err.Error())
					return
				}


				//claims is just an interface
				claims,ok:=actualToken.Claims.(jwt.MapClaims)

				if !ok {
					utils.SendJSONError(w, http.StatusUnauthorized, "Invalid token claims")
					return
				}



				userid, ok1 := claims["userid"].(string)
				username, ok2 := claims["username"].(string)

				if !ok1 || !ok2 {
					utils.SendJSONError(w, http.StatusUnauthorized, "Invalid token claims data")
					return
				}

				// Add user info to context for next handlers
				ctx := context.WithValue(r.Context(), ContextKeyUserID, userid)
				ctx = context.WithValue(ctx, ContextKeyUsername, username)

				next.ServeHTTP(w, r.WithContext(ctx))
	}),
)
}


// func loggingMiddleware(next http.Handler) http.Handler{
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 	log.Println(r.Method, r.URL.Path)
// 	next.ServeHTTP(w, r)
// 	})
// }