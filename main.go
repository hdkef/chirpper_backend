package main

import (
	"chirpper_backend/controller"
	mydriver "chirpper_backend/driver"
	"chirpper_backend/utils"
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

func init() {
	_ = godotenv.Load()
}

func main() {

	client := mydriver.ConnectDB()
	defer client.Close()

	authHandler := controller.Auth{}
	endPointsHandler := controller.EndPoints{}
	router := mux.NewRouter()

	router.HandleFunc("/auth/login", utils.Cors(authHandler.Login(client)))
	router.HandleFunc("/auth/register", utils.Cors(authHandler.Register(client)))
	router.HandleFunc("/endpoints/feeds", utils.Cors(endPointsHandler.Feeds(client)))
	router.HandleFunc("/auth/sendemailver", utils.Cors(authHandler.SendEmailVer(client)))
	router.HandleFunc("/auth/verifyemailver", utils.Cors(authHandler.VerifyEmailVer(client)))

	// spa := spaHandler{staticPath: "dist/angular", indexPath: "index.html"}
	// router.PathPrefix("/").Handler(spa)

	PORT := os.Getenv("PORT")

	address := fmt.Sprintf(":%v", PORT)

	fmt.Println("serving and listening")

	err := http.ListenAndServe(address, router)
	if err != nil {
		panic(err)
	}
}
