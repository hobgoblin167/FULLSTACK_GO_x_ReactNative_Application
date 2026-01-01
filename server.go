package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// База данных
var conn *pgx.Conn

// JSON
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CORS
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
}


func loginHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		json.NewEncoder(w).Encode(AuthResponse{
			Success: false,
			Message: "Неверный формат данных",
		})
		return
	}

	var passwordFromDB string
	err = conn.QueryRow(
		context.Background(),
		"SELECT password FROM userBase WHERE name = $1",
		req.Username,
	).Scan(&passwordFromDB)

	if err != nil {
		json.NewEncoder(w).Encode(AuthResponse{
			Success: false,
			Message: "Пользователь не найден",
		})
		return
	}

	if passwordFromDB != req.Password {
		json.NewEncoder(w).Encode(AuthResponse{
			Success: false,
			Message: "Неверный пароль",
		})
		return
	}

	json.NewEncoder(w).Encode(AuthResponse{
		Success: true,
		Message: "Успешный вход",
	})
}


func registerHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		json.NewEncoder(w).Encode(AuthResponse{
			Success: false,
			Message: "Неверный формат данных",
		})
		return
	}

	var exists string
	err = conn.QueryRow(
		context.Background(),
		"SELECT password FROM userBase WHERE name = $1",
		req.Username,
	).Scan(&exists)

	if err == nil {
		json.NewEncoder(w).Encode(AuthResponse{
			Success: false,
			Message: "Пользователь уже существует",
		})
		return
	}

	_, err = conn.Exec(
		context.Background(),
		"INSERT INTO userBase (name, password) VALUES ($1, $2)",
		req.Username,
		req.Password,
	)

	if err != nil {
		json.NewEncoder(w).Encode(AuthResponse{
			Success: false,
			Message: "Ошибка создания пользователя",
		})
		return
	}

	json.NewEncoder(w).Encode(AuthResponse{
		Success: true,
		Message: "Регистрация успешна",
	})
}

func main() {
	// Загружаем env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	dbURL := os.Getenv("DATABASE_URL")

	conn, err = pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v\n", err)
	}
	defer conn.Close(context.Background())

	fmt.Println("Подключение к БД успешно")

	// Создаём таблицу
	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS userBase (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы: %v\n", err)
	}

	fmt.Println("Таблица userBase готова")

	// Роутинг
	mux := http.NewServeMux()
	mux.HandleFunc("/api/login", loginHandler)
	mux.HandleFunc("/api/register", registerHandler)

	fmt.Println("API сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
