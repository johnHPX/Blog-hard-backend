package utils

import (
	"net/http"
)

func HeaderMethods(next http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, X-PINGOTHER, X-Auth-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		} else if r.Method == method {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func Authenticate(nextFunction http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := validateToken(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			response := CreateHttpErrorResponse(http.StatusUnauthorized, 01, err, "token is invalid")
			EncodeResponse(nil, w, response)
			return
		}
		nextFunction(w, r)
	}
}
