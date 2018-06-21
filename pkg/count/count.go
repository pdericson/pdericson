package count

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Count struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

var c chan Count

func worker() {
	db, err := sql.Open("sqlite3", "./count.db")
	if err != nil {
		log.Fatalf("count: worker: %s\n", err.Error())
	}
	defer db.Close()

	for {
		count := <-c

		stmt, err := db.Prepare(`insert into count(date, name) values(date('now'), ?)`)
		if err != nil {
			log.Fatalf("count: worker: %s\n", err.Error())
		}
		defer stmt.Close()

		_, err = stmt.Exec(count.Name)
		if err != nil {
			log.Fatalf("count: worker: %s\n", err.Error())
		}
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	count := Count{}

	if c == nil {
		c = make(chan Count)

		go worker()
	}

	db, err := sql.Open("sqlite3", "./count.db")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer db.Close()

	sqlStmt := `
        create table if not exists count(
            date text not null,
            name text not null
        );
        `

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("count: worker: %s\n", err.Error())
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "GET" {
		count.Name = vars["name"]

		stmt, err := db.Prepare(`select count(*) from count where name=?`)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer stmt.Close()

		err = stmt.QueryRow(count.Name).Scan(&count.Count)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(count)
		return
	}
	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&count)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if count.Name == "" {
			http.Error(w, `count.Name == ""`, 400)
			return
		}

		c <- count

		w.WriteHeader(204)
		return
	}
}
