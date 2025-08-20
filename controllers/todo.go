package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"main/model"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type response struct {
	Id          int64  `json:"id_todo,omitempty"`
	Status      string `json:"status,omitempty"`
	Description string `json:"description"`
	Message     string `json:"message"`
}

type Response struct {
	Status      string      `json:"status"`
	Message     string      `json:"message"`
	Description string      `json:"description"`
	Data        interface{} `json:"todo"`
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var todo model.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Fatal("Error decoding request body:", err)
	}

	insertedId := model.CreateTodo(todo)

	res := response{
		Id:          insertedId,
		Status:      "success",
		Description: todo.Description,
		Message:     "Todo created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

func GetAllTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	todos, err := model.GetAllTodo()
	if err != nil {
		log.Fatal("Error fetching todos:", err)
	}

	response := Response{
		Status:      "200",
		Message:     "Todos fetched successfully",
		Description: "List of all todos",
		Data:        todos,
	}

	json.NewEncoder(w).Encode(response)
}

func GetDetailTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal("Error converting id to integer:", err)
	}

	todo, err := model.GetDetailTodo(int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// todo tidak ketemu
			w.WriteHeader(http.StatusNotFound)
			webResponse := response{
				Status:      "404",
				Message:     "Todo not found",
				Description: "No todo found with the given ID",
			}
			json.NewEncoder(w).Encode(webResponse)
			return
		}

		// error lain (misalnya DB down)
		w.WriteHeader(http.StatusInternalServerError)
		webResponse := response{
			Status:      "500",
			Message:     "Internal server error",
			Description: err.Error(),
		}
		json.NewEncoder(w).Encode(webResponse)
		return
	}

	res := Response{
		Status:      "200",
		Message:     "Todo fetched successfully",
		Description: "Details of the todo",
		Data:        todo,
	}

	json.NewEncoder(w).Encode(res)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var todo model.Todo
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedRows := model.UpdateTodo(int64(id), todo)
	if updatedRows == 0 {
		// kalau nggak ada row yang kena update
		webResponse := response{
			Status:      "404",
			Message:     "Todo not found or not updated",
			Description: "No todo found with the given ID or nothing to update",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(webResponse)
		return
	}

	res := response{
		Id:          int64(id),
		Status:      todo.Status,
		Description: todo.Description,
		Message:     "Todo updated successfully",
	}

	json.NewEncoder(w).Encode(res)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	deletedRows := model.DeleteTodo(int64(id))
	if deletedRows == 0 {
		webResponse := response{
			Status:      "404",
			Message:     "Todo not found",
			Description: "No todo found with the given ID",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(webResponse)
		return
	}

	res := response{
		Id:          int64(id),
		Status:      "deleted",
		Description: "Todo deleted successfully",
		Message:     "Todo deleted successfully",
	}

	json.NewEncoder(w).Encode(res)
}
