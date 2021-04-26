package utils

import "net/http"

func ResClearSite(res *http.ResponseWriter) {
	(*res).Header().Set("Access-Control-Expose-Headers", "CLEARBEARER")
	(*res).Header().Set("CLEARBEARER", "storage")
}
