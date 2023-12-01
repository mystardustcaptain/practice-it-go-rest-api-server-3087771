package backend

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	Port   string
	DB     *sql.DB
	DBType string
	DBPath string
}

func (a *App) Initialize() {
	DB, err := sql.Open(a.DBType, a.DBPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	a.DB = DB
}

func (a *App) Run() {
	http.HandleFunc("/", helloWorld)
	fmt.Println("Server started on port ", a.Port)
	log.Fatal(http.ListenAndServe(a.Port, nil))
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World\n")
}

func (a *App) GetData() {
	rows, err := a.DB.Query("SELECT id, name, inventory, price FROM products")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var p Product

		rows.Scan(&p.id, &p.name, &p.inventory, &p.price)
		fmt.Println("Product: ", p.id, " ", p.name, " ", p.inventory, " ", p.price)
	}
}

type Product struct {
	id        int
	name      string
	inventory int
	price     int
}
