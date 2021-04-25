package utils

import "net/http"

func cors(res *http.ResponseWriter) {
	(*res).Header().Set("Access-Control-Allow-Origin", "*")
	(*res).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*res).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Auth, Username, ID")
}

//Cors is to handle cors
func Cors(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		cors(&res)
		if req.Method == "OPTIONS" {
			res.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(res, req)
	}
}
