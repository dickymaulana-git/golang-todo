package routers

import (
	"main/controllers"
	"main/middleware"
	"main/model"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	// AUTH API
	router.HandleFunc("/auth/signup", model.SignUp).Methods("POST")
	router.HandleFunc("/auth/signin", model.SignIn).Methods("POST")
	router.HandleFunc("/auth/signout", middleware.IsAuthorized(model.SignOut)).Methods("POST", "OPTIONS")
	router.HandleFunc("/auth/getuser", middleware.IsAuthorized(controllers.GetUserHandler)).Methods("GET", "OPTIONS")

	// TODO API
	router.HandleFunc("/api/todo", middleware.IsAuthorized(controllers.CreateTodo)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/todo", middleware.IsAuthorized(controllers.GetAllTodo)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/todo/detail/{id}", middleware.IsAuthorized(controllers.GetDetailTodo)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/todo/detail/{id}", middleware.IsAuthorized(controllers.UpdateTodo)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/todo/detail/{id}", middleware.IsAuthorized(controllers.DeleteTodo)).Methods("DELETE", "OPTIONS")

	//REMINDER API
	router.HandleFunc("/reminder/create", middleware.IsAuthorized(controllers.CreateReminder)).Methods("POST", "OPTIONS")
	router.HandleFunc("/reminder", middleware.IsAuthorized(controllers.GetAllReminders)).Methods("GET", "OPTIONS")

	return router
}
