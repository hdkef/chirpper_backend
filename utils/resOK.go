package utils

import (
	"encoding/json"
	"net/http"
)

//ResOK is to
func ResOK(res http.ResponseWriter, msg string) {
	obj := struct {
		Message string
	}{Message: msg}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(obj)
}
