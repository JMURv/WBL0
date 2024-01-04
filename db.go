package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

func createDBTables() error {
	var err error

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS orders (
        id SERIAL PRIMARY KEY,
        details TEXT
    )
`)

	if err != nil {
		return err
	}

	return nil
}

func saveDataToDB(data string) (*Order, error) {
	var orderData Order
	if err := json.Unmarshal([]byte(data), &orderData); err != nil {
		log.Println("Error decoding JSON:", err)
		return &orderData, err
	}

	_, err := db.Exec("INSERT INTO orders (id, details) VALUES ($1, $2)", orderData.ID, orderData.Details)
	if err != nil {
		log.Println("Error saving data to DB:", err)
		return &orderData, err
	}
	return &orderData, nil
}

func getOrderFromDB(orderID string) (*Order, error) {
	row := db.QueryRow("SELECT id, details FROM orders WHERE id = $1", orderID)
	order := &Order{}
	err := row.Scan(&order.ID, &order.Details)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order with ID %s not found", orderID)
		}
		return nil, err
	}
	return order, nil
}

func loadCacheFromDB() {
	rows, err := db.Query("SELECT id, details FROM orders")
	if err != nil {
		log.Printf("Error querying orders from the database: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		if err = rows.Scan(&order.ID, &order.Details); err != nil {
			log.Printf("Error scanning row from the database: %v", err)
			continue
		}

		cacheVar.SetDefault(strconv.Itoa(int(order.ID)), &order)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating through rows: %v", err)
	}

	log.Println("Cache has been loaded from the database")
}
