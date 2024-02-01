package main

import (
	"GoAuth/app"
	"GoAuth/controllers"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.Use(app.JwtAuthentication)
	// defer models.GetDB().Close()
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	fmt.Println("port", port)
	
	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/user/new/verifyOTP",controllers.ConfirmOTP).Methods("POST")
	router.HandleFunc("/api/testing",controllers.ExtraHandler).Methods("GET")
	router.HandleFunc("/api/contacts/add",controllers.NewContact).Methods("POST")

	err := http.ListenAndServe(":8080",router)
	if err != nil {
		fmt.Print(err)
	}
}
