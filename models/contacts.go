package models

import (
	u "GoAuth/utils"
)

type Contacts struct {
	Number		string		`json:"number"`
	AccountID 	uint 		`json:"account_id"`
	Name 		string 		`json:"name"`
}

func (account *Account) AddContact(contact *Contacts) (map[string]interface{}) {
	userID := account.ID
	existingContact := Contacts{}
	err := GetDB().Where("account_id = ? AND number = ?",userID,contact.Number).First(&existingContact).Error
	if err != nil {
		return u.Message(false,"Error checking the contact")
	}
	newContact := Contacts{
		Number: contact.Number,
		AccountID: userID,
		Name: contact.Name,
	}
	err = GetDB().Create(&newContact).Error
	if err != nil {
		return u.Message(false,"Error while checking the request")
	}
	resp := u.Message(true,"Contact added successfully")
	resp["contact"] = newContact
	return resp
}