package model

import (
	"context"
	"database/sql"
	"log"
	"main/config"
)

type Todo struct {
	Id          int64  `json:"id_todo"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

func CreateTodo(todo Todo) int64 {
	db := config.CreateConnection()
	defer db.Close()

	sqlStatement := "INSERT INTO todo (status, description) VALUES ($1, $2)"
	var id int64
	ctx := context.Background()

	_, err := db.ExecContext(ctx, sqlStatement, todo.Status, todo.Description)
	if err != nil {
		log.Fatal("Error inserting todo:", err)
		return 0
	}

	return id
}

func GetAllTodo() ([]Todo, error) {
	db := config.CreateConnection()
	defer db.Close()

	sqlStatement := "SELECT * FROM todo"
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatal("Error fetching todos:", err)
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.Id, &todo.Status, &todo.Description)
		if err != nil {
			log.Fatal("Error scanning todo:", err)
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func GetDetailTodo(id int64) (Todo, error) {
	db := config.CreateConnection()
	defer db.Close()

	sqlStatement := "SELECT id, status, description FROM todo WHERE id=$1"
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		return Todo{}, err
	}
	defer rows.Close()

	todo := Todo{}
	if rows.Next() {
		err := rows.Scan(&todo.Id, &todo.Status, &todo.Description)
		if err != nil {
			return Todo{}, err
		}
		return todo, nil
	}

	// kalau tidak ada row, balikin error
	return Todo{}, sql.ErrNoRows
}

func UpdateTodo(id int64, todo Todo) int64 {
	db := config.CreateConnection()
	defer db.Close()
	sqlStatement := `UPDATE todo SET status=$2, description=$3 WHERE id=$1`
	res, err := db.Exec(sqlStatement, id, todo.Status, todo.Description)
	if err != nil {
		log.Fatal("Error updating todo:", err)
		return 0
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal("Error getting rows affected:", err)
		return 0
	}
	if rowsAffected == 0 {
		log.Println("No rows updated")
		return 0
	}
	return rowsAffected
}

func DeleteTodo(id int64) int64 {
	db := config.CreateConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM todo WHERE id=$1`
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatal("Error deleting todo:", err)
		return 0
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal("Error getting rows affected:", err)
		return 0
	}
	if rowsAffected == 0 {
		log.Println("No rows deleted")
		return 0
	}
	return rowsAffected
}
