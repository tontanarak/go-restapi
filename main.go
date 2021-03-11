package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type SpaSession struct {
	ID       int    `json:"id"`
	CUSTOMER string `json: "customer"`
	TIME     string `json: "time"`
}

var database, dataerr = sql.Open("sqlite3", "./spa.db")

func main() {
	render := mux.NewRouter()
	//databased create table
	if dataerr != nil {
		fmt.Println("db file %s: %s", "spa.db", dataerr)
	}
	createtable, err := database.Prepare("CREATE TABLE IF NOT EXISTS spasession (id INTEGER PRIMARY KEY, customer TEXT, time TEXT)")
	if err != nil {
		fmt.Println(err)
	}
	createtable.Exec()
	defer database.Close()

	//render.HandleFunc("/spasessions", SpaSessionList).Methods(http.MethodGet)        //show all spa session ****without jwt
	render.HandleFunc("/spasessions", SpaSessionCreate).Methods(http.MethodPost)     //add spa session
	render.HandleFunc("/spasessions/{id}", SpaSessionDel).Methods(http.MethodDelete) //delete spa session
	render.HandleFunc("/spasessions/{id}", SpaSessionBook).Methods(http.MethodPatch) //book spa session
	http.ListenAndServe(":1234", render)

}

func SpaSessionList(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	fmt.Println(t)

	allSessions, _ := database.Query("SELECT * FROM spasession;")
	defer allSessions.Close()
	sessions := []SpaSession{}
	var session SpaSession
	for allSessions.Next() {
		if err := allSessions.Scan(&session.ID, &session.CUSTOMER, &session.TIME); err != nil {
			fmt.Println(err)
		}
		sessions = append(sessions, session)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&sessions)
}

func SpaSessionCreate(w http.ResponseWriter, r *http.Request) {
	var s SpaSession
	_ = json.NewDecoder(r.Body).Decode(&s)
	addsession := fmt.Sprintf("INSERT INTO spasession (customer,time) VALUES ('%s','%s')", "Avaliable", s.TIME)
	result, _ := database.Exec(addsession)
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}
	query := "SELECT * FROM spasession WHERE id = " + strconv.Itoa(int(id))
	Sessions, err := database.Query(query)
	if err != nil {
		fmt.Println(w, err)
	}
	for Sessions.Next() {
		if err := Sessions.Scan(&s.ID, &s.CUSTOMER, &s.TIME); err != nil {
			fmt.Println(err)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GenToken(s.ID, s.CUSTOMER, s.TIME))
}

func SpaSessionDel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	delsession := fmt.Sprintf("DELETE FROM spasession WHERE id=%d", id)
	_, err := database.Exec(delsession)
	if err != nil {
		fmt.Println(w, err)
	}
}

func SpaSessionBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	var s SpaSession
	_ = json.NewDecoder(r.Body).Decode(&s)

	booksession := fmt.Sprintf("UPDATE spasession SET customer = '%s' WHERE id=%d", s.CUSTOMER, id)
	_, err := database.Exec(booksession)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	query := "SELECT * FROM spasession WHERE id= " + strconv.Itoa(int(id))
	Sessions, _ := database.Query(query)
	for Sessions.Next() {
		if err := Sessions.Scan(&s.ID, &s.CUSTOMER, &s.TIME); err != nil {
			fmt.Println(w, err)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GenToken(s.ID, s.CUSTOMER, s.TIME))
}

type customClaims struct {
	ID       int    `json:"id"`
	CUSTOMER string `json: "customer"`
	TIME     string `json: "time"`
	jwt.StandardClaims
}

func GenToken(id int, customer string, time string) string {
	claims := &customClaims{
		id,
		customer,
		time,
		jwt.StandardClaims{},
	}
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("tonssecretkeyisnothing"))

	return token
}
