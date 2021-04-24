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

	mydriver.ConnectDB()

	authHandler := controller.Auth{}
	router := mux.NewRouter()

	router.HandleFunc("/login", authHandler.Login())
	router.HandleFunc("/register", authHandler.Register())

	PORT := os.Getenv("PORT")

	address := fmt.Sprintf(":%v", PORT)

	fmt.Println("serving and listening")

	err := http.ListenAndServe(address, router)
	if err != nil {
		panic(err)
	}
}
