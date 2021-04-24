package controller

import (
	"chirpper_backend/models"
	"chirpper_backend/utils"
	"errors"
	"fmt"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

//AuthStruct to group auth controller
type Auth struct {
}

//Login is to login
func (a *Auth) Login() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		_, err := verifyToken(res, req)

		if err == nil {
			http.Redirect(res, req, "/test", 302)
			return
		}

		bearerCookie := &http.Cookie{}

		var user models.User = models.User{
			ID:   "1",
			Name: "hadekha",
			Role: "admin",
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

	}
}

//Register is to register
func (a *Auth) Register() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

	}
}

//createToken is to create token after login
func createToken(user *models.User) (string, error) {
	secret := os.Getenv("SECRET")

	var claims = jwt.MapClaims{
		"ID":   user.ID,
		"Name": user.Name,
		"Role": user.Role,
	}

	fmt.Println("secret", secret, "claims", claims)

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
