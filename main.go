package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

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

	migrationsDir := "./migrations"
	err = applyMigrations(db, migrationsDir)
	checkErr(err)

	stmt, err := db.Prepare("INSERT INTO userinfo(username, created) values(?,?)")
	checkErr(err)

	res, err := stmt.Exec("marcel0ll", time.Now())
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	log.Println("id", id)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("GET /static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("GET /", func(res http.ResponseWriter, req *http.Request) {
		HTML(res, req, components.Index(count))
	})

	http.HandleFunc("POST /add", func(res http.ResponseWriter, req *http.Request) {
		count = count + 1
		HTML(res, req, components.Index(count))
	})

	port := envPortOr("8080")
	log.Println("Listenting at", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func envPortOr(port string) string {
	// If `PORT` variable in environment exists, return it
	if envPort := os.Getenv("PORT"); envPort != "" {
		return ":" + envPort
	}
	// Otherwise, return the value of `port` variable from function argument
	return ":" + port
}

func applyMigrations(db *sql.DB, migrationsDir string) error {
	log.Println("Running migrations...")
	migrationFiles, err := filepath.Glob(filepath.Join(migrationsDir, "*.up.sql"))
	if err != nil {
		return err
	}

	sort.Strings(migrationFiles)

	for _, migrationFile := range migrationFiles {
		log.Println("Migrating", migrationFile)
		content, err := os.ReadFile(migrationFile)
		if err != nil {
			return err
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return err
		}
	}

	return nil
}
