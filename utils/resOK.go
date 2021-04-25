package utils

import (
	"encoding/json"
	"net/http"
)

//ResOK is to
func ResOK(res http.ResponseWriter, msg string) {
	obj := struct {
		MESSAGE string
	}{MESSAGE: msg}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(obj)
}
