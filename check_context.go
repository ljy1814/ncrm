package main

import (
	"net/http"
)

func checkContext(request *http.Request) bool {
	if len(request.Header.Get("X-Zc-Nonce")) <= 0 {
		return false
	}
	return true
}
