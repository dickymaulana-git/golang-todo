package routers

import (
	"main/controllers"
	"main/middleware"
	"main/model"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/signup", model.SignUp).Methods("POST")
	router.HandleFunc("/api/signin", model.SignIn).Methods("POST")

	router.HandleFunc("/api/todo", middleware.IsAuthorized(controllers.CreateTodo)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/todo", middleware.IsAuthorized(controllers.GetAllTodo)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/todo/detail/{id}", middleware.IsAuthorized(controllers.GetDetailTodo)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/todo/detail/{id}", middleware.IsAuthorized(controllers.UpdateTodo)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/todo/detail/{id}", middleware.IsAuthorized(controllers.DeleteTodo)).Methods("DELETE", "OPTIONS")

	return router
}
