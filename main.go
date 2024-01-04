package main

import (
	"database/sql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"github.com/patrickmn/go-cache"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	db         *sql.DB
	cacheVar   *cache.Cache
	natsStream stan.Conn
)

func init() {
	var err error

	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error while loading .env file")
		return
	}
	log.Println("Env file has been loaded")

	// Подключаемся к БД
	db, err = sql.Open("postgres", os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("DB connection has been loaded successfully")

	// Создание таблицы
	err = createDBTables()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Table 'orders' has been created successfully or already exists")

	// Инициализируем кэш
	cacheVar = cache.New(5*time.Minute, 10*time.Minute)

	// Загружаем данные из БД в кэш
	loadCacheFromDB()

	// Подключаемся к NATS
	clusterID := "test-cluster"
	clientID := "test-id"

	natsStream, err = stan.Connect(clusterID, clientID, stan.NatsURL("nats://nats:4222"))
	if err != nil {
		log.Fatal(err)
	}

	// Подписка на канал "test-channel" и обработка полученных сообщений
	_, err = natsStream.Subscribe("test-channel", func(m *stan.Msg) {
		msgData := string(m.Data)
		log.Printf("Received a message: %s\n", msgData)

		// Сохраняем полученные данные в кэш
		dbData, _ := saveDataToDB(msgData)
		cacheVar.SetDefault(strconv.Itoa(int(dbData.ID)), dbData)
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/order/{id}", getOrderHandler).Methods(http.MethodGet)
	router.HandleFunc("/send-test-message", sendTestMessageHandler).Methods(http.MethodPost)

	log.Println("Server has been started at :8000")
	log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, router)))
}
