package main

import (
	"chirpper_backend/controller"
	mydriver "chirpper_backend/driver"
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

	router.HandleFunc("/login", authHandler.Login(client))
	router.HandleFunc("/register", authHandler.Register(client))
	router.HandleFunc("/feeds", endPointsHandler.Feeds(client))

	PORT := os.Getenv("PORT")

	address := fmt.Sprintf(":%v", PORT)

	fmt.Println("serving and listening")

	err := http.ListenAndServe(address, router)
	if err != nil {
		panic(err)
	}
}
