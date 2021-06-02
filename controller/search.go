package controller

import (
	"chirpper_backend/utils"
	"encoding/json"
	"errors"
	"net/http"

	"cloud.google.com/go/firestore"
)

//Search someone
func (x *EndPoints) Search(client *firestore.Client) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		valid := verifyToken(req)
		if valid != true {
			utils.ResClearSite(&res)
			utils.ResError(res, http.StatusUnauthorized, errors.New("INVALID TOKEN"))
			return
		}

		var payload struct {
			Searchkey string
		}

		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		db := NewDBRepo(client)

		result, err := db.FindAllByField("users", "Username", payload.Searchkey)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		err = json.NewEncoder(res).Encode(&result)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}
	}
}
