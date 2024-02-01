package app

import (
	"context"
	u "GoAuth/utils"
	"net/http"
	"os"
	"strings"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"GoAuth/models"
)

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
		notAuth := []string{"/api/user/new","/api/user/login","/api/user/new/verifyOTP"}	//These path do not require auth
		requestPath := r.URL.Path

		for _,value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w,r) //Serve the request if no Authentication needed
				return
			}
		}
		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			response = u.Message(false,"Missing Auth Token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-type","application/json")
			u.Respond(w,response)
			return
		}
		splitted := strings.Split(tokenHeader," ") //The token comes in the format `Bearer {token-body},check`

		if len(splitted) != 2 {
			fmt.Println(len(splitted))
			response = u.Message(false,"Invalid Token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("content-type","application/json")
			u.Respond(w,response)
			return
		}
		tokenPart := splitted[1]
		tk := &models.Token{}

		token,err := jwt.ParseWithClaims(tokenPart,tk,func(token *jwt.Token) (interface{},error){
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil {
			response = u.Message(false,"Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-type","application/json")
			u.Respond(w,response)
			return
		}
		if !token.Valid {
			response = u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}
		//Everything Went Well
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //Proceed in the middleware chain
	})
}