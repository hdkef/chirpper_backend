package controller

import (
	"net/http"

	"cloud.google.com/go/firestore"
)

type EndPoints struct {
}

func (x *EndPoints) Feeds(client *firestore.Client) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

	}
}