package controller

import (
	"chirpper_backend/utils"
	"encoding/json"
	"errors"
	"net/http"

	"cloud.google.com/go/firestore"
)

//Profile give info about someone's profile
func (x *EndPoints) Profile(client *firestore.Client) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		valid := verifyToken(req)
		if valid != true {
			utils.ResClearSite(&res)
			utils.ResError(res, http.StatusUnauthorized, errors.New("INVALID TOKEN"))
			return
		}

		var payload struct {
			ID string
		}

		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		db := NewDBRepo(client)

		result, err := db.FindOneByID("users", payload.ID)
		if err != nil {
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
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		} else if err != nil && err.Error() == "NO RESULT" {
			payloadToBeSent.Feed = []map[string]interface{}{}
		} else {
			payloadToBeSent.Feed = result2
		}

		err = json.NewEncoder(res).Encode(&payloadToBeSent)
		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}
	}
}
