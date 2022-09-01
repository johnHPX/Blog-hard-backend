package authn

import (
	"net/http"

	"github.com/johnHPX/blog-hard-backend/internal/appl/service"
	"github.com/johnHPX/blog-hard-backend/internal/infra/utils/responseAPI"
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

		tokenFunc := service.NewAccessService()
		err := tokenFunc.ValidateAToken(r)

		if err != nil && err.Error() == "Token is expired" {

			userID, err := tokenFunc.ExtractInvalideToken(r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				response := responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 01, err, "ti")
				responseAPI.EncodeResponse(nil, w, response)
				return
			}
			atoken, err := tokenFunc.GenerateNewToken(userID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				response := responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 02, err, "ti")
				responseAPI.EncodeResponse(nil, w, response)
				return
			}
			r.Header.Set("Authorization", atoken)
		}

		nextFunction(w, r)
	}
}
