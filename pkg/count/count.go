package count

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	connStr = "user=pdericson password=pdericson dbname=pdericson sslmode=disable"
)

type Count struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	First string `json:"first"`
	Last  string `json:"last"`
}

var c chan Count

func worker() {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("count: worker: %s\n", err.Error())
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

	for {
		count := <-c

		stmt, err := db.Prepare(`insert into count(date, name) values(now(), $1)`)
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

//swagger:parameters PostHandler
type PostHandlerParams struct {
	// in: body
	// required: true
	Count Count `json:"count"`
}

// swagger:route POST /count count PostHandler
//
// Increment the count.
//
// Responses:
//   204:
//   400:
func PostHandler(w http.ResponseWriter, r *http.Request) {
	count := Count{}

	if c == nil {
		c = make(chan Count)

		go worker()
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

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
}

//swagger:parameters GetHandler
type GetHandlerParams struct {
	// The name of the thing being counted.
	//
	// in: path
	// required: true
	Name string `json:"name"`
}

// swagger:route GET /count/{name} count GetHandler
//
// Get the count.
//
// Responses:
//   200: Count
func GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	count := Count{}

	if c == nil {
		c = make(chan Count)

		go worker()
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer db.Close()

	sqlStmt := `
        create table if not exists count(
            date timestamp without time zone not null,
            name text                        not null
        );
        `

	_, err = db.Exec(sqlStmt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	count.Name = vars["name"]

	// count
	stmt, err := db.Prepare(`select count(*) from count where name=$1`)
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

	// first
	stmt, err = db.Prepare(`select date from count where name=$1 order by date asc limit 1`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(count.Name).Scan(&count.First)
	if err != nil {
		if err != sql.ErrNoRows {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	// last
	stmt, err = db.Prepare(`select date from count where name=$1 order by date desc limit 1`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(count.Name).Scan(&count.Last)
	if err != nil {
		if err != sql.ErrNoRows {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(count)
	return
}
