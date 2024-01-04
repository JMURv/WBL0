package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func sendToNATS(message string) error {
	err := natsStream.Publish("test-channel", []byte(message))
	if err != nil {
		log.Println("Error publishing to NATS:", err)
		return err
	}
	return nil
}

func sendResponse(w http.ResponseWriter, order *Order) {
	jsonData, err := json.Marshal(order)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while encoding JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonData); err != nil {
		http.Error(w, fmt.Sprintf("Error while writing response: %v", err), http.StatusInternalServerError)
	}
}
