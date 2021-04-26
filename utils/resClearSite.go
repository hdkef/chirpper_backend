package utils

import "net/http"

func ResClearSite(res *http.ResponseWriter) {
	(*res).Header().Set("Clear-Site-Data", "storage")
}
