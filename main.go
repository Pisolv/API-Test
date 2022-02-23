package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Person struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Country  string `json:"country"`
	Language string `json:"language"`
	Contact  string `json:"contact"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "AyenC123^"
	dbname   = "myaccount"
)

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	rows, err := db.Query("SELECT * FROM person")
	if err != nil {
		log.Fatal(err)
	}

	var people []Person

	for rows.Next() {
		var person Person
		rows.Scan(&person.ID, &person.Name, &person.Email, &person.Country, &person.Language, &person.Contact)
		people = append(people, person)
	}

	peopleBytes, _ := json.MarshalIndent(people, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

	defer rows.Close()
	defer db.Close()
}
func Getdatabyemail(w http.ResponseWriter, r *http.Request) {

	db := OpenConnection()

	vars := mux.Vars(r)
	mail := vars["email"]

	fmt.Println(`id ==== := `, mail)

	rows, err := db.Query("SELECT name, email, country ,language , contact FROM email=$", mail)
	if err != nil {
		log.Fatal(err)
	}

	var people []Person

	for rows.Next() {
		var person Person
		rows.Scan(&person.ID, &person.Name, &person.Email, &person.Country, &person.Language, &person.Contact)
		people = append(people, person)
	}

	peopleBytes, _ := json.MarshalIndent(people, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)
	defer rows.Close()
	defer db.Close()

}

func POSTHandler(w http.ResponseWriter, r *http.Request) {

	db := OpenConnection()

	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}
	sqlStatement := `INSERT INTO public.person(id, name, email, country, language, contact) VALUES ($1,$2,$3,$4,$5,$6);`
	_, err = db.Exec(sqlStatement, p.ID, p.Name, p.Email, p.Country, p.Language, p.Contact)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()

}

func PUTHandler(w http.ResponseWriter, r *http.Request) {

	db := OpenConnection()
	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}
	sqlStatement := `UPDATE person SET  name=$2 , country=$3, language=$4, contact=$5 WHERE email= $1;`
	_, err = db.Exec(sqlStatement, p.Email, "NewName", "NewCountry", "NewLanguage", "NewContact")

	if err != nil {
		panic(err)

	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()
	defer db.Close()

}

func DELETEHandler(w http.ResponseWriter, r *http.Request) {

	db := OpenConnection()

	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}
	sqlStatement := `DELETE FROM person WHERE email=$1;`
	_, err = db.Exec(sqlStatement, p.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()

}

func main() {
	http.HandleFunc("/", GETHandler)
	http.HandleFunc("/insert", POSTHandler)
	http.HandleFunc("/getbyemail/{email}", Getdatabyemail)
	http.HandleFunc("/update", PUTHandler)
	http.HandleFunc("/delete", DELETEHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
