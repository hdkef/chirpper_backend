package utils

import (
	"fmt"
	"net/http"
)

func ResError(res http.ResponseWriter, errCode int, err error) {
	res.WriteHeader(errCode)
	msg := fmt.Sprintf("%v", err.Error())
	res.Write([]byte(msg))
}
