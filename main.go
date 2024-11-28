package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MySQL connection details
const (
	dbUser     = "root"        // Replace with your MySQL username
	dbPassword = "Thanooj@001" // Replace with your MySQL password
	dbName     = "toronto_thanooj_log"
)

// Database connection
var db *sql.DB

func main() {
	// Connect to MySQL database
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", dbUser, dbPassword, dbName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	defer db.Close()

	// Verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error verifying database connection: %v\n", err)
	}
	fmt.Println("Connected to the database successfully!")

	// Set up HTTP routes
	http.HandleFunc("/current-time", currentTimeHandler)

	// Start the server
	port := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// Handler for /current-time endpoint
func currentTimeHandler(w http.ResponseWriter, r *http.Request) {
	// Get the current Toronto time
	location, err := time.LoadLocation("America/Toronto")
	if err != nil {
		http.Error(w, "Error loading Toronto timezone", http.StatusInternalServerError)
		return
	}
	torontoTime := time.Now().In(location)
	formattedTime := torontoTime.Format("2006-01-02 15:04:05")
	// Log the time to the database
	_, err = db.Exec("INSERT INTO thanoojtable_log (timestamp) VALUES (?)", formattedTime)
	if err != nil {
		http.Error(w, "Error logging time to database", http.StatusInternalServerError)
		return
	}

	// Respond with the current Toronto time in JSON format
	response := map[string]string{
		"current_time": torontoTime.Format(time.RFC1123),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}