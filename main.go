package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Response struct {
	gorm.Model

	Service  string
	Rating   int
	Feedback string
}

var db *gorm.DB
var err error

func main() {
	// load environment variables
	
	
	
	
	os.Setenv("HOST", "localhost")
	os.Setenv("DBPORT", "5433")
	os.Setenv("USER", "postgres")
	os.Setenv("NAME", "feedback")
	os.Setenv("PASSWORD", "postgres")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	// database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, password, dbPort)

	// open connection to db
	db, err = gorm.Open("postgres", dbURI)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to database")
	}

	// close connection to db when main func terminates
	defer db.Close()

	// make migration to the db if they have not already been created
	db.AutoMigrate(&Response{})

	router := handleRequests()

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getResponses(w http.ResponseWriter, r *http.Request) {
	var responses []Response
	db.Find(&responses)
	json.NewEncoder(w).Encode(responses)
}

func createResponse(w http.ResponseWriter, r *http.Request) {
	var response Response
	json.NewDecoder(r.Body).Decode(&response)

	createdResponse := db.Create(&response)
	err = createdResponse.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&response)
	}
}

func handleRequests() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/responses", getResponses).Methods("GET")
	router.HandleFunc("/create/response", createResponse).Methods("POST")

	return router
}