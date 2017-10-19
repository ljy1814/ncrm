package main

import (
	"net/http"
	"strconv"
)

func checkContext(request *http.Request) bool {
	if len(request.Header.Get("X-Zc-Nonce")) <= 0 {
		return false
	}
	return true
}

func getListDomainEnv(request *http.Request) int64 {
	env := request.Header.Get("X-Zc-Env")
	if len(env) <= 0 {
		return -1
	}
	ienv, _ := strconv.ParseInt(env, 9, 64)
	return ienv
}
