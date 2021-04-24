package main

import (
	"chirpper_backend/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func helloworld() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		utils.ResOK(res, "okkkkkkk")
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", helloworld())
	err := http.ListenAndServe(":4040", router)
	if err != nil {
		log.Fatal("something wrong")
	}
}
