package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/marcel0ll/go-template/components"
	_ "github.com/mattn/go-sqlite3"
)

func HTML(res http.ResponseWriter, req *http.Request, comp templ.Component) error {
	res.Header().Set("Content-Type", "text/html")
	return comp.Render(req.Context(), res)
}

var count = 0

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func main() {
	db, err := sql.Open("sqlite3", "./db/sqlite.db")
	checkErr(err)

	stmt, err := db.Prepare("INSERT INTO userinfo(username, created) values(?,?)")
	checkErr(err)

	res, err := stmt.Exec("marcel0ll", "2012-12-09")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	println("id", id)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("GET /static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("GET /", func(res http.ResponseWriter, req *http.Request) {
		HTML(res, req, components.Index(count))
	})

	http.HandleFunc("POST /add", func(res http.ResponseWriter, req *http.Request) {
		count = count + 1
		HTML(res, req, components.Index(count))
	})

	println("Listenting at 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
