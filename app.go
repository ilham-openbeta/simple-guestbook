package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// App berisi mux Router dan sql DB
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initialize untuk inisialisasi aplikasi
func (a *App) Initialize(host, port, database, username, password string) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	var err error
	a.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()

	fmt.Printf("Initialize Successful\n")
}

// Run untuk menjalankan aplikasi
func (a *App) Run(port string) {
	fmt.Printf("Server Running on Port %s\n", port)
	port = fmt.Sprintf(":%s", port)
	log.Fatal(http.ListenAndServe(port, a.Router))
}

func (a *App) initializeRoutes() {
	api := a.Router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/message", a.createMessage).Methods("POST")
	api.HandleFunc("/messages", a.getMessages).Methods("GET")
	a.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	// a.Router.HandleFunc("/", a.home)
}

// func (a *App) home(w http.ResponseWriter, r *http.Request) {
// 	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Hello World!"})
// }

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) createMessage(w http.ResponseWriter, r *http.Request) {
    var m message
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&m); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    if err := m.createMessage(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

	// respondWithJSON(w, http.StatusCreated, m)
	respondWithJSON(w, http.StatusCreated, map[string]string{"status": "message submitted"})
}

func (a *App) getMessages(w http.ResponseWriter, r *http.Request) {
	var count, start int
	query := r.URL.Query()

	qcount := query.Get("count")
	qstart := query.Get("start")

	if  qcount == "" ||  qstart == "" {
		start = 0
		count = 0
	} else {
		start, _ = strconv.Atoi(qstart)
		count, _ = strconv.Atoi(qcount)
	}

    messages, err := getMessages(a.DB, start, count)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, messages)
}

