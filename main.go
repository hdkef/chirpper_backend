package main

import (
	"chirpper_backend/controller"
	mydriver "chirpper_backend/driver"
	"chirpper_backend/utils"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {

	client := mydriver.ConnectDB()
	defer client.Close()

	authHandler := controller.Auth{}
	endPointsHandler := controller.EndPoints{}
	router := mux.NewRouter()

	router.HandleFunc("/login", utils.Cors(authHandler.Login(client)))
	router.HandleFunc("/register", utils.Cors(authHandler.Register(client)))
	router.HandleFunc("/feeds", utils.Cors(endPointsHandler.Feeds(client)))
	router.HandleFunc("/sendemailver", utils.Cors(authHandler.SendEmailVer(client)))
	router.HandleFunc("/verifyemailver", utils.Cors(authHandler.VerifyEmailVer(client)))

	PORT := os.Getenv("PORT")

	address := fmt.Sprintf(":%v", PORT)

	fmt.Println("serving and listening")

	err := http.ListenAndServe(address, router)
	if err != nil {
		panic(err)
	}
}
