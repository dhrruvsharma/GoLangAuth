package controllers

import (
	"GoAuth/models"
	u "GoAuth/utils"
	"encoding/json"
	"fmt"
	"net/http"
	extras "GoAuth/extras"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}

	//decode the request body into struct and failed if any error occur
	err := json.NewDecoder(r.Body).Decode(account) 
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := account.Create() //Create account
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}

	//decode the request body into struct and failed if any error occur
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(account.Email, account.Password)
	u.Respond(w, resp)
}

var ConfirmOTP = func(w http.ResponseWriter,r *http.Request) {
	account := &models.Account{}
	var otp models.OTP
	err := json.NewDecoder(r.Body).Decode(&otp)
	if err != nil {
		u.Respond(w, u.Message(false,"Invalid Request"))
		return
	}
	account.Email = otp.Email
	resp := account.VerifyOtp(otp.Otp)
	u.Respond(w,resp)
}

var ExtraHandler = func(w http.ResponseWriter,r *http.Request){
	userValue := r.Context().Value("user")
	account := &models.Account{}
	if userValue != nil {
		userID, ok := userValue.(uint)
		if ok {
			account.ID = userID
			fmt.Println("User ID :",userID)
			extras.Extras(account)
			u.Respond(w,u.Message(true,"ID fetched successfully"))
			return
		}
	}
	http.Error(w,"Unauthorized",http.StatusUnauthorized)
	return
}

var NewContact = func(w http.ResponseWriter,r *http.Request) {
	uservalue := r.Context().Value("user")
	contact := &models.Contacts{}
	json.NewDecoder(r.Body).Decode(contact)
	account := &models.Account{}
	if uservalue != nil {
		userID,ok := uservalue.(uint)
		if ok {
			account.ID = userID
			resp := account.AddContact(contact)
			u.Respond(w,resp)
			return
		}
	}
}