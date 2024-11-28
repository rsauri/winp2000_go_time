package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv" // Import the godotenv package
)

// Create new struct
type TimeLog struct {
	ID        int    `json:"id"`
	Timestamp string `json:"timestamp"`
}

var db *sql.DB // Global database connection

func saveTime(w http.ResponseWriter, r *http.Request) {
	// Load the Eastern Time Zone location
	loc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		log.Fatalf("Error loading time zonessss: %v", err)
	}

	currentTime := time.Now().In(loc).Format("2006-01-02 15:04:05")
	// currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Insert data into the database
	insertQuery := "INSERT INTO time_log (timestamp) VALUES (?)"
	_, err = db.Exec(insertQuery, currentTime)
	if err != nil {
		log.Fatalf("Error inserting data: %v", err)
	}

	// Show the data
	json.NewDecoder(r.Body).Decode(&currentTime)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentTime)
}

func retrieveTime(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, timestamp FROM time_log"

	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}
	defer rows.Close()

	// Iterate through the rows
	var timeLogs []TimeLog
	for rows.Next() {
		var timeLog TimeLog

		err := rows.Scan(&timeLog.ID, &timeLog.Timestamp)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}

		timeLogs = append(timeLogs, timeLog)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Fatalf("Row iteration error: %v", err)
	}

	// Respond with the list
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(timeLogs)
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get MySQL connection details from environment variables
	user := os.Getenv("DB_USER")
	if user == "" {
		log.Fatal("DB_USER environment variable not set")
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		log.Fatal("DB_PASSWORD environment variable not set")
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "mysql" // Default to localhost if not set
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306" // Default to 3306 if not set
	}

	database := os.Getenv("DB_NAME")
	if database == "" {
		log.Fatal("DB_NAME environment variable not set")
	}

	// Format the connection string (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)

	// Open a connection to the database
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Successfully connected to MySQL!")

	// Handlers
	http.HandleFunc("/current-time", saveTime)
	http.HandleFunc("/show-time", retrieveTime)

	// Start the server
	fmt.Println("Server is running on port 80...")
	log.Fatal(http.ListenAndServe(":80", nil))
}
