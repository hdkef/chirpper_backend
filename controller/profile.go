package controller

import (
	"chirpper_backend/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
)

//Profile give info about someone's profile
func (x *EndPoints) Profile(client *firestore.Client) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		fmt.Println("Profile")

		var payload struct {
			ID string
		}

		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		fmt.Println("ID", payload.ID)

		db := NewDBRepo(client)

		result, err := db.FindOneByID("users", payload.ID)
		if err != nil {
			fmt.Println("db error", err)
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		payloadToBeSent := struct {
			Username       string
			Email          string
			AvatarURL      string
			Desc           string
			FollowerCount  int
			FollowingCount int
			Feed           []map[string]interface{}
		}{
			Username:       result["Username"].(string),
			Email:          result["Username"].(string),
			AvatarURL:      result["AvatarURL"].(string),
			Desc:           result["Desc"].(string),
			FollowerCount:  len(result["followers"].([]interface{})),
			FollowingCount: len(result["followings"].([]interface{})),
		}

		result2, err := db.FindAllSubColByIDField("users", payload.ID, "ID", payload.ID, "feed")
		if err != nil && err.Error() != "NO RESULT" {
			fmt.Println("db error", err)
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		} else if err != nil && err.Error() == "NO RESULT" {
			payloadToBeSent.Feed = nil
		} else {
			payloadToBeSent.Feed = result2
		}

		fmt.Println(payloadToBeSent)

		err = json.NewEncoder(res).Encode(&payloadToBeSent)
		if err != nil {
			fmt.Println("db error", err)
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}
	}
}
