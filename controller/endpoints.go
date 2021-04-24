package controller

import (
	"net/http"
)

type EndPoints struct {
}

func (x *EndPoints) Feeds() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
	}
}
