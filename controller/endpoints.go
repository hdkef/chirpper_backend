package controller

import (
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"
)

type EndPoints struct {
}

func (x *EndPoints) Feed(client *firestore.Client) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		_, err := verifyToken(res, req)
		if err != nil {
			return
		}
		payload := struct {
			MESSAGE string
		}{MESSAGE: "feed served"}

		err = json.NewEncoder(res).Encode(&payload)
		if err != nil {
			return
		}
	}
}
