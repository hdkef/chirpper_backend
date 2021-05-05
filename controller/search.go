package controller

import (
	"chirpper_backend/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
)

//Search someone
func (x *EndPoints) Search(client *firestore.Client) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		fmt.Println("Search")

		var payload struct {
			Searchkey string
		}

		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		fmt.Println("Searchkey", payload.Searchkey)

		db := NewDBRepo(client)

		result, err := db.FindAllByField("users", "Username", payload.Searchkey)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		fmt.Println(result)

		err = json.NewEncoder(res).Encode(&result)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}
	}
}
