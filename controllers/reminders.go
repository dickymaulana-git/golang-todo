package controllers

import (
	"encoding/json"
	"main/model"
	"net/http"
)

type reminderResponse struct {
	Id      int64  `json:"id_todo,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message"`
}

func CreateReminder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var responseData model.Reminder
	err := json.NewDecoder(r.Body).Decode(&responseData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	insertedId := model.CreateReminder(responseData)

	res := reminderResponse{
		Id:      insertedId,
		Status:  "success",
		Message: "Reminder created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

func GetAllReminders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reminders, err := model.GetAllReminders()
	if err != nil {
		http.Error(w, "Error fetching reminders", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status  string           `json:"status"`
		Message string           `json:"message"`
		Data    []model.Reminder `json:"data"`
	}{
		Status:  "200",
		Message: "Reminders fetched successfully",
		Data:    reminders,
	}

	json.NewEncoder(w).Encode(response)
}
