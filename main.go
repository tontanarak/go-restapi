package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

	render.HandleFunc("/spasessions", SpaSessionList).Methods(http.MethodGet)        //show all spa session ****without jwt
	render.HandleFunc("/spasessions", SpaSessionCreate).Methods(http.MethodPost)     //add spa session
	render.HandleFunc("/spasessions/{id}", SpaSessionDel).Methods(http.MethodDelete) //delete spa session
	render.HandleFunc("/spasessions/{id}", SpaSessionBook).Methods(http.MethodPatch) //book spa session

	http.ListenAndServe(":1234", render)

}

func SpaSessionList(w http.ResponseWriter, r *http.Request) {

	allSessions, _ := database.Query("SELECT * FROM spasession;")
	defer allSessions.Close()
	sessions := []SpaSession{}
	var session SpaSession
	for allSessions.Next() {
		if err := allSessions.Scan(&session.ID, &session.CUSTOMER, &session.TIME); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(err)
		}
		sessions = append(sessions, session)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&sessions)
}

func SpaSessionCreate(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("Missing Authorization Header")
		return
	}
	result, jwterr := jwtHandler(r)
	if jwterr != nil {
		fmt.Println(jwterr)
	}
	admin := result.(jwt.MapClaims)["admin"].(bool)
	if admin {
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
		var res SpaSession
		for Sessions.Next() {
			if err := Sessions.Scan(&res.ID, &res.CUSTOMER, &res.TIME); err != nil {
				fmt.Println(err)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	} else {
		w.WriteHeader(http.StatusForbidden)
		fmt.Println("Only admin can create session")
	}

}

func SpaSessionDel(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Missing Authorization Header"))
		return
	}
	result, jwterr := jwtHandler(r)
	if jwterr != nil {
		fmt.Println(jwterr)
	}
	admin := result.(jwt.MapClaims)["admin"].(bool)
	if admin {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		delsession := fmt.Sprintf("DELETE FROM spasession WHERE id=%d", id)
		_, err := database.Exec(delsession)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
		fmt.Println("Only admin can delete session")
	}

}

func SpaSessionBook(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Missing Authorization Header"))
		return
	}
	result, jwterr := jwtHandler(r)
	if jwterr != nil {
		fmt.Println(jwterr)
	}
	name := result.(jwt.MapClaims)["name"].(string)

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	booksession := fmt.Sprintf("UPDATE spasession SET customer = '%s' WHERE id=%d ", name, id)
	_, err := database.Exec(booksession)
	if err != nil {
		fmt.Println(err)
	}
	query := "SELECT * FROM spasession WHERE id= " + strconv.Itoa(int(id))
	Sessions, _ := database.Query(query)
	var res SpaSession
	for Sessions.Next() {
		if err := Sessions.Scan(&res.ID, &res.CUSTOMER, &res.TIME); err != nil {
			fmt.Println(err)
		}
	}
	if res.ID == 0 {
		fmt.Println("Incorrect Session ID")
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}

func jwtHandler(r *http.Request) (jwt.Claims, error) {
	body := strings.Split(r.Header["Authorization"][0], " ")
	tokenString := body[1]
	signingKey := []byte("thekeyiskeykeykey")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, err
}
