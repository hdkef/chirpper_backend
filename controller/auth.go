package controller

import (
	"chirpper_backend/models"
	"chirpper_backend/utils"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const expiresaAtInt int64 = 604800

//AuthStruct to group auth controller
type Auth struct {
}

//Login is to login
func (a *Auth) Login(client *firestore.Client) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		var loginForm struct {
			Username string
			Password string
		}

		err := json.NewDecoder(req.Body).Decode(&loginForm)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		//here implement finding hashedpass from database
		db := NewDBRepo(client)
		result, err := db.FindOneByField("users", "Username", loginForm.Username)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}
		if result == nil {
			utils.ResError(res, http.StatusUnauthorized, errors.New("Username not found"))
			return
		}

		var hashedPass string = result["Password"].(string)

		err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(loginForm.Password))
		if err != nil {
			utils.ResError(res, http.StatusUnauthorized, err)
			return
		}

		var user models.User = models.User{
			ID:        result["ID"].(string),
			Username:  result["Username"].(string),
			Email:     result["Email"].(string),
			AvatarURL: result["AvatarURL"].(string),
		}

		tokenString, err := createToken(&user)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		user.Token = tokenString

		err = json.NewEncoder(res).Encode(user)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}
	}
}

//sendEmailVer is sent email verification to user
func (a *Auth) SendEmailVer(client *firestore.Client) http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		var payload struct {
			Email string
		}

		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		//implement random digit generation
		code := rand.Int31()
		var codePayload map[string]interface{} = map[string]interface{}{
			"Email": payload.Email,
			"Code":  code,
		}

		db := NewDBRepo(client)
		_, err = db.InsertOne("emailver", codePayload)

		emailMe := os.Getenv("EMAIL")
		pswdMe := os.Getenv("PSWD")
		host := "smtp.gmail.com"
		emailTo := []string{payload.Email}
		port := os.Getenv("SMTPPORT")
		addr := fmt.Sprintf("%s:%s", host, port)
		auth := smtp.PlainAuth("", emailMe, pswdMe, host)
		msgString := "Subject: Email Verification\n\n" + "Hi, this is your email verification code : " + fmt.Sprint(code)
		msgBytes := []byte(msgString)

		err = smtp.SendMail(addr, auth, emailMe, emailTo, msgBytes)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		utils.ResOK(res, "email verification sent")
	}
}

//VerifyEmailVer is to verify the digit code given to user's email
func (a *Auth) VerifyEmailVer(client *firestore.Client) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var payload struct {
			Email string
			Code  string
		}

		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		db := NewDBRepo(client)
		result, err := db.FindOneByField("emailver", "Email", payload.Email)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		verified := strconv.FormatInt(result["Code"].(int64), 10) == payload.Code

		if verified == true {
			err = db.DeleteByID("emailver", result["ID"].(string))
			if err != nil {
				utils.ResError(res, http.StatusInternalServerError, err)
				return
			}
			res.Write([]byte(fmt.Sprint(verified)))
		} else {
			res.Write([]byte(fmt.Sprint(verified)))
		}
	}
}

//Register is to register
func (a *Auth) Register(client *firestore.Client) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		var registerForm struct {
			Username string
			Email    string
			Password string
		}

		err := json.NewDecoder(req.Body).Decode(&registerForm)

		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		hashedPassByte, err := bcrypt.GenerateFromPassword([]byte(registerForm.Password), 5)

		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		//here implement inserting registerForm to database
		db := NewDBRepo(client)
		_, err = db.InsertOne("users", map[string]interface{}{
			"Username": registerForm.Username,
			"Password": string(hashedPassByte),
			"Email":    registerForm.Email,
		})
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		utils.ResOK(res, "register success")
	}
}

//createToken is to create token after login
func createToken(user *models.User) (string, error) {
	secret := os.Getenv("SECRET")

	var claims = jwt.MapClaims{
		"ID":       user.ID,
		"Username": user.Username,
		"Email":    user.Email,
		"exp":      time.Now().Unix() + expiresaAtInt,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//VerifyToken is to verify token in header
func verifyToken(req *http.Request) bool {

	fmt.Println("VerifyToken")

	token := req.Header.Get("BEARER")

	fmt.Println("Token", token)

	return verifyTokenString(token)
}

//VerifyTokenString will verify token from parameter string
func verifyTokenString(token string) bool {

	fmt.Println("VerifyTokenString")

	if token == "" {
		return false
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return false
	}

	if parsedToken.Valid == true {
		return true
	} else {
		return false
	}
}
