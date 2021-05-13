package main

import (
	"chirpper_backend/controller"
	mydriver "chirpper_backend/driver"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

//ServerHTTP is to implement interface so that spaHandler can be use as Handler
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.staticPath, r.URL.Path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

//First init the env
func init() {
	_ = godotenv.Load()
}

func main() {

	client := mydriver.ConnectDB()
	defer client.Close()

	authHandler := controller.Auth{}
	endPointsHandler := controller.EndPoints{}
	router := mux.NewRouter()

	router.HandleFunc("/auth/login", authHandler.Login(client))
	router.HandleFunc("/auth/register", authHandler.Register(client))
	router.HandleFunc("/endpoints/checktoken", endPointsHandler.CheckToken(client))
	router.HandleFunc("/auth/sendemailver", authHandler.SendEmailVer(client))
	router.HandleFunc("/auth/verifyemailver", authHandler.VerifyEmailVer(client))
	router.HandleFunc("/endpoints/ws", endPointsHandler.EstablishWS(client))
	router.HandleFunc("/endpoints/profile", endPointsHandler.Profile(client))
	router.HandleFunc("/endpoints/search", endPointsHandler.Search(client))
	router.HandleFunc("/endpoints/comment", endPointsHandler.EstablishComment(client))
	router.HandleFunc("/endpoints/postwithimage", endPointsHandler.PostWithImage(client))
	router.HandleFunc("/endpoints/setting", endPointsHandler.Setting(client))

	spa := spaHandler{staticPath: os.Getenv("STATICPATH"), indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	PORT := os.Getenv("PORT")

	address := fmt.Sprintf(":%v", PORT)

	fmt.Println("serving and listening")

	err := http.ListenAndServe(address, router)
	if err != nil {
		panic(err)
	}
}
