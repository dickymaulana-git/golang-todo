package model

import (
	"context"
	"encoding/json"
	"log"
	helper "main/Helper"
	"main/config"
	"net/http"

	"github.com/lib/pq"
)

type User struct {
	Id        int    `json:"id,omitempty"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Country   string `json:"country,omitempty"`
}

type SignUpRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Country   string `json:"country"`
	Password  string `json:"password"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	db := config.CreateConnection()
	defer db.Close()

	var user SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": "Invalid request body",
		})
		return
	}

	// ✅ Field validations per field
	fieldErrors := make(map[string]string)
	if user.Email == "" {
		fieldErrors["email"] = "Email is required"
	}
	if user.Password == "" {
		fieldErrors["password"] = "Password is required"
	} else if len(user.Password) < 6 {
		fieldErrors["password"] = "Password must be at least 6 characters"
	}
	if user.FirstName == "" {
		fieldErrors["first_name"] = "First name is required"
	}
	if user.LastName == "" {
		fieldErrors["last_name"] = "Last name is required"
	}
	if user.Country == "" {
		fieldErrors["country"] = "Country is required"
	}

	if len(fieldErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":       http.StatusBadRequest,
			"message":      "Validation errors",
			"field_errors": fieldErrors,
		})
		return
	}

	// ✅ Hash password
	hashedPassword, err := helper.GenerateHashPassword(user.Password)
	if err != nil {
		log.Println("Error generating password hash:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": "Internal Server Error",
		})
		return
	}

	// ✅ Insert into DB (return user id)
	sqlStatement := `INSERT INTO users(email, password, first_name, last_name, country)
	                 VALUES ($1, $2, $3, $4, $5) RETURNING id`

	ctx := context.Background()
	var userID int
	err = db.QueryRowContext(ctx, sqlStatement,
		user.Email, hashedPassword, user.FirstName, user.LastName, user.Country,
	).Scan(&userID)

	if err != nil {
		// Check for duplicate email (unique violation)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  http.StatusConflict,
				"message": "Email already exists",
			})
			return
		}

		log.Println("Error inserting new user:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": "Could not create user",
		})
		return
	}

	// ✅ Success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  http.StatusCreated,
		"message": "User registered successfully",
		"data": map[string]interface{}{
			"id":         userID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"country":    user.Country,
		},
	})
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	db := config.CreateConnection()
	defer db.Close()

	var RequestUser Authentication

	err := json.NewDecoder(r.Body).Decode(&RequestUser)
	if err != nil {
		helper.WriteError(w, "Invalid request body")
		return
	}

	// ✅ Validate input fields first
	fieldErrors := make(map[string]string)
	if RequestUser.Email == "" {
		fieldErrors["email"] = "Email is required"
	}
	if RequestUser.Password == "" {
		fieldErrors["password"] = "Password is required"
	}
	if len(fieldErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"errors":  fieldErrors,
		})
		return
	}

	var user User
	err = db.QueryRow("SELECT email, password FROM users WHERE email=$1", RequestUser.Email).
		Scan(&user.Email, &user.Password)
	if err != nil {
		helper.WriteError(w, "User not found")
		return
	}

	check := helper.CheckPasswordHash(RequestUser.Password, user.Password)
	if !check {
		helper.WriteError(w, "Invalid password")
		return
	}

	validToken, err := helper.GenerateJWT(user.Email)
	if err != nil {
		helper.WriteError(w, "Error generating token")
		return
	}

	// ✅ return token + user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user": map[string]interface{}{
			"email": user.Email,
		},
		"token": validToken,
	})
}

func GetUserByEmail(email string) (User, error) {
	db := config.CreateConnection()
	defer db.Close()

	var user User
	err := db.QueryRow("SELECT id, email, first_name, last_name, country FROM users WHERE email=$1", email).
		Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Country)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	// Since JWT is stateless, sign-out can be handled on the client side by deleting the token.
	// Optionally, you can implement token blacklisting on the server side if needed.

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  http.StatusOK,
		"message": "User signed out successfully",
	})
}
