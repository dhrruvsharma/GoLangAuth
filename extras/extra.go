package extras

import (
	"GoAuth/models"
	"fmt"
	u "GoAuth/utils"
)

func Extras(account *models.Account) map[string]interface{} {
	dbAccount := &models.Account{}
	err := models.GetDB().Table("accounts").Where("ID = ?",account.ID).First(dbAccount).Error

	if err != nil {
		fmt.Println(err)
		return u.Message(false,"Not Successfull")
	}
	fmt.Println(dbAccount)
	return u.Message(true,"Successfull")
}
