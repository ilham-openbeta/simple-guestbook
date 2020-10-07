package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	godotenv.Load()
	a.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"))

	ensureTableExist()
	clearTable()
	os.Exit(m.Run())
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS messages (
    id int(5) NOT NULL,
    name varchar(50) NOT NULL,
    message varchar(160) NOT NULL,
    contact varchar(50) NOT NULL,
    created_at datetime NOT NULL
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

func ensureTableExist() {
	a.DB.Exec(tableCreationQuery)
}

func clearTable() {
	a.DB.Exec("DELETE FROM messages")
	a.DB.Exec("ALTER TABLE messages AUTO_INCREMENT = 1")
}

// func TestHelloWorld(t *testing.T) {
// 	req, _ := http.NewRequest("GET", "/", nil)
// 	response := executeRequest(req)

// 	checkResponseCode(t, http.StatusOK, response.Code)

// 	var m map[string]interface{}
// 	json.Unmarshal(response.Body.Bytes(), &m)

// 	if m["message"] != "Hello World!" {
// 		t.Errorf("Expected ID to be 'Hello World!'. Got '%v'", m["message"])
// 	}
// }

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestCreateMesssage(t *testing.T) {
	clearTable()
	name := "John Doe"
	message := "Lorem Ipsum Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce tristique augue eu blandit posuere."
	contact := "test@test.com"

	var jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "message":"%s", "contact":"%s"}`, name, message, contact))
	req, _ := http.NewRequest("POST", "/api/v1/message", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)
    
    if m["status"] != "message submitted" {
		t.Errorf("Expected status to be 'message submitted'. Got '%v'", m["status"])
	}

	// if m["id"] != 1.0 {
	// 	t.Errorf("Expected ID to be '1'. Got '%v'", m["id"])
	// }

	// if m["name"] != name {
	// 	t.Errorf("Expected name to be '%s'. Got '%v'", name, m["name"])
	// }

	// if m["message"] != message {
	// 	t.Errorf("Expected message to be '%s'. Got '%v'", message, m["message"])
	// }

	// if m["contact"] != contact {
	// 	t.Errorf("Expected contact to be '%s'. Got '%v'", contact, m["contact"])
	// }
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/api/v1/messages", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetAllMessages(t *testing.T) {
	clearTable()
	addMessages(5)

	req, _ := http.NewRequest("GET", "/api/v1/messages", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m []interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if len(m) < 1 {
		t.Errorf("Expected non empty array. Got 0")
	}
}

func addMessages(count int) {
	var msg []message
	var m message

	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {

		m.Name = fmt.Sprintf("Test %d", i)
		m.Message = fmt.Sprintf("Lorem ipsum dolor sit amet number %d", i)
		m.Contact = fmt.Sprintf("0812345678%d", i)
		msg = append(msg, m)
	}
	a.DB.Create(&msg)
}

func TestGetAllMessagesInRange(t *testing.T) {
	clearTable()
	addMessages(10)

	req, _ := http.NewRequest("GET", "/api/v1/messages?start=5&count=5", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m []interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if len(m) != 5 {
		t.Errorf("Expected 5 data")
	}
}
