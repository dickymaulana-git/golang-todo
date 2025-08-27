package model

import (
	"context"
	"log"
	"main/config"
	"time"
)

type Reminder struct {
	ID            int       `json:"id"`
	PlaceName     string    `json:"place_name"`
	PlaceCity     string    `json:"place_city"`
	EventDate     time.Time `json:"event_date"` // PostgreSQL DATE → Go time.Time
	ReminderRange int       `json:"reminder_range"`
	Price         int64     `json:"price"` // BIGINT → int64
	Status        string    `json:"status"`
}

func CreateReminder(todo Reminder) int64 {
	db := config.CreateConnection()
	defer db.Close()

	sqlStatement := `
		INSERT INTO reminders (place_name, place_city, event_date, reminder_range, price, status) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id`

	var id int64
	ctx := context.Background()

	err := db.QueryRowContext(ctx, sqlStatement,
		todo.PlaceName,
		todo.PlaceCity,
		todo.EventDate, // time.Time works fine with DATE in Postgres
		todo.ReminderRange,
		todo.Price,
		todo.Status,
	).Scan(&id)

	if err != nil {
		log.Fatal("❌ Error inserting reminder:", err)
		return 0
	}

	return id
}

func GetAllReminders() ([]Reminder, error) {
	db := config.CreateConnection()
	defer db.Close()

	sqlStatement := "SELECT * FROM reminders"
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatal("❌ Error fetching reminders:", err)
		return nil, err
	}
	defer rows.Close()

	var reminders []Reminder
	for rows.Next() {
		var reminder Reminder
		err := rows.Scan(&reminder.ID, &reminder.PlaceName, &reminder.PlaceCity, &reminder.EventDate, &reminder.ReminderRange, &reminder.Price, &reminder.Status)
		if err != nil {
			log.Fatal("❌ Error scanning reminder:", err)
			return nil, err
		}
		reminders = append(reminders, reminder)
	}

	return reminders, nil
}
