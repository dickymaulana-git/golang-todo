package middleware

import (
	"encoding/json"
	helper "main/Helper"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

var (
	secretKey string = "secretkeyjwt"
)

func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Header["Token"] == nil {
			errResponse := helper.SetError(helper.Error{}, "Missing token in header")
			json.NewEncoder(w).Encode(errResponse)
			return
		}

		var mySigningKey = []byte(secretKey)
		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("Unexpected signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return mySigningKey, nil
		})
		if err != nil {
			errResponse := helper.SetError(helper.Error{}, "Token Expired: "+err.Error())
			json.NewEncoder(w).Encode(errResponse)
			return
		}

		_, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			errResponse := helper.SetError(helper.Error{}, "Unauthorized access")
			json.NewEncoder(w).Encode(errResponse)
			return
		}
		handler.ServeHTTP(w, r)
	}
}
