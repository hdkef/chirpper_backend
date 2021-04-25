package controller

import (
	"chirpper_backend/models"
	"chirpper_backend/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//AuthStruct to group auth controller
type Auth struct {
}

//Login is to login
func (a *Auth) Login(client *firestore.Client) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		_, err := verifyToken(res, req)

		if err == nil {
			http.Redirect(res, req, "/feed", 302)
			return
		}

		var loginForm struct {
			Username string
			Password string
		}

		err = json.NewDecoder(req.Body).Decode(&loginForm)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		//here implement finding hashedpass from database
		db := NewFindRepo(client)
		result, err := db.FindOneByUsername("users", loginForm.Username)
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

		bearerCookie := &http.Cookie{}

		var user models.User = models.User{
			ID:       result["ID"].(string),
			Username: result["Username"].(string),
			Email:    result["Email"].(string),
		}

		tokenString, err := createToken(&user)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		bearerCookie.Name = "bearer"
		bearerCookie.Value = tokenString
		// bearerCookies.Expires = time.Now().Add(5 * time.Minute)

		http.SetCookie(res, bearerCookie)

		err = json.NewEncoder(res).Encode(user)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
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
		db := NewInsertRepo(client)
		err = db.InsertOne("users", map[string]interface{}{
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
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		fmt.Println("tokenString", tokenString)
		return "", err
	}

	return tokenString, nil
}

//authenticate is to verify token
func verifyToken(res http.ResponseWriter, req *http.Request) (jwt.MapClaims, error) {

	var claimsModel = jwt.MapClaims{
		"Id":        "",
		"Subject":   "",
		"ExpiresAt": 0,
		"Issuer":    "",
		"Role":      "",
	}

	token := &http.Cookie{}

	storedCookie, _ := req.Cookie("bearer")
	if storedCookie == nil {
		return claimsModel, errors.New("BEARER COOKIE NOT FOUND")
	}

	token = storedCookie

	parsedToken, err := jwt.Parse(token.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return claimsModel, errors.New("Token was modified")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if ok && parsedToken.Valid {
		claimsModel = claims
		return claimsModel, nil
	}

	return claimsModel, errors.New("Token was either invalid or expired")
}
