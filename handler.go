package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["id"]

	// Попытка получения данных из кэша
	if cachedData, found := cacheVar.Get(orderID); found {
		log.Println("Cached data has been found")
		sendResponse(w, cachedData.(*Order))
		return
	}

	// Попытка получения данных из БД
	dbData, err := getOrderFromDB(orderID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
		return
	}

	// Сохранение данных в кэш
	cacheVar.SetDefault(orderID, dbData)
	sendResponse(w, dbData)
}

func sendTestMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Тестовая отправка сообщения
	message := `{"id": 20, "details": "test details"}`
	err := sendToNATS(message)
	if err != nil {
		http.Error(w, "Test message wasn't send", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Test message sent successfully"))
}
