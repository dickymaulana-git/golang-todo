package middleware

import (
	"context"
	"encoding/json"
	helper "main/Helper"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

var (
	secretKey string = "secretkeyjwt"
)

type contextKey string

// ✅ Exported supaya bisa diakses dari controller
const UserEmailKey contextKey = "email"

func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errResponse := helper.SetError(helper.Error{}, "Missing token in header")
			json.NewEncoder(w).Encode(errResponse)
			return
		}

		// Expect "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // "Bearer " not found
			errResponse := helper.SetError(helper.Error{}, "Invalid token format")
			json.NewEncoder(w).Encode(errResponse)
			return
		}

		var mySigningKey = []byte(secretKey)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
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

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			errResponse := helper.SetError(helper.Error{}, "Unauthorized access")
			json.NewEncoder(w).Encode(errResponse)
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			errResponse := helper.SetError(helper.Error{}, "Email not found in token")
			json.NewEncoder(w).Encode(errResponse)
			return
		}

		ctx := context.WithValue(r.Context(), UserEmailKey, email)
		handler.ServeHTTP(w, r.WithContext(ctx))
	}
}
