package controllers

import (
	"encoding/json"
	"main/middleware"
	"main/model"
	"net/http"
)

type UserResponse struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Country   string `json:"country"`
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	emailCtx := r.Context().Value(middleware.UserEmailKey)
	if emailCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	email := emailCtx.(string)

	user, err := model.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Map to response struct
	resp := UserResponse{
		Id:        user.Id,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Country:   user.Country,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
