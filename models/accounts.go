package models

import (
	"fmt"
	"os"
	"strings"

	u "GoAuth/utils"

	"time"

	otpGen "GoAuth/otp"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//JWT claims struct

type Token struct {
	UserId uint
	jwt.StandardClaims
}

//Struct to represent a user account

type Account struct {
	gorm.Model
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Token     string    `json:"token" sql:"-"`
	Verified  bool      `json:"verified"`
	OTP       string    `json:"-"`
	OtpExpiry time.Time `json:"-"`
	Contacts []Contacts `gorm:"foreignKey:AccountID" json:"contacts"`
}

//Struct to represent an OTP

type OTP struct {
	Otp   string `json:"otp"`
	Email string `json:"email"`
}

// Validate Incoming User
func (account *Account) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email Address Is required"), false
	}
	if len(account.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}
	if account.Username == "" {
		return u.Message(false, "Username is required"), false
	}
	//Email must be unique
	temp := &Account{}

	//Check for errors and duplicate emails
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(&temp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection Error, Please Retry"), false
	}
	if temp.Email != "" && temp.Verified == true {
		return u.Message(false, "User already exists"), false
	}
	return u.Message(false, "Requirement Passed"), true
}

func (account *Account) Create() map[string]interface{} {
	if resp, ok := account.Validate(); !ok {
		return resp
	}

	existingUser := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(&existingUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection Error.Please retry")
	}

	otp, err := otpGen.GenerateOtp()

	if err != nil {
		return u.Message(false, "Failed to generate otp")
	}

	err = otpGen.SendOTP([]string{account.Email}, otp)

	if err != nil {
		fmt.Print(err)
		return u.Message(false, "Error while Sending OTP. Please retry.")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	if existingUser.Email != "" && !existingUser.Verified {
		account.OTP = otp
		existingUser.OTP = account.OTP
		existingUser.OtpExpiry = time.Now().UTC().Add(5 * time.Minute)
		existingUser.Verified = false
		existingUser.Username = account.Username	
		existingUser.Password = account.Password
		existingUser.Token = account.Token
		GetDB().Save(existingUser)

		existingUser.Password = ""

		response := u.Message(true, "OTP sent successfully and Information updated")
		response["account"] = existingUser
		return response
	}

	account.OTP = otp
	account.OtpExpiry = time.Now().UTC().Add(5 * time.Minute)
	account.Verified = false

	GetDB().Create(account)
	if account.ID < 0 {
		return u.Message(false, "Failed to create accont,connection error")
	}
	//Create new Jwt token for the newly registered account


	account.Password = "" //Delete Password

	response := u.Message(true, "OTP Sent Successfully")
	response["account"] = account
	return response
}

func Login(email, password string) map[string]interface{} {
	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "User not found")
		}
		return u.Message(false, "Connection error, Please retry")
	}
	if account.Verified == false {
		return u.Message(false, "Please verify your account first")
	}
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid Login Credentials. Please try again")
	}
	//Worked
	account.Password = ""

	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	return resp
}

func (account *Account) VerifyOtp(userOTP string) map[string]interface{} {

	//Check if the OTP is still valid
	dbAccount := &Account{}
	fmt.Printf("Email : %s\n",account.Email)
	err  := GetDB().Table("accounts").Where("email = ?",account.Email).First(&dbAccount).Error
	fmt.Printf(dbAccount.OTP)
	if err != nil {
		u.Message(false,"Error while fetching the account")
	}

	if dbAccount == nil {
		return u.Message(false,"Account not found in the database")
	}

	// account = dbAccount

	if time.Now().UTC().Before(dbAccount.OtpExpiry) && userOTP == dbAccount.OTP {
		dbAccount.Verified = true
		//Clear the OTP after successful validation
		dbAccount.OTP = ""
		dbAccount.OtpExpiry = time.Time{}
		resp := u.Message(true, "Otp Verified Successfully")
		resp["account"] = dbAccount
		err := GetDB().Save(dbAccount)
		if err != nil {
			return u.Message(false,"Error saving changes to database. Please try again later.")
		}
		return resp
	}
	return u.Message(false, "Otp Verification Failed")
}

func GetUser(u uint) *Account {
	acc := &Account{}
	GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" {
		return nil
	}
	acc.Password = ""
	return acc
}
