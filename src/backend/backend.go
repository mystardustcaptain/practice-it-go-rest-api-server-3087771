package backend

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	Port   string
	DB     *sql.DB
	DBType string
	DBPath string
	Router *mux.Router
}

func (a *App) Initialize() {
	db, err := sql.Open(a.DBType, a.DBPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	a.DB = db

	a.Router = mux.NewRouter()
	a.initializeRouters()
}

func (a *App) initializeRouters() {
	a.Router.HandleFunc("/", helloWorld_GET).Methods("GET")
	a.Router.HandleFunc("/", helloWorld_POST).Methods("POST")

	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product/{id}", a.getProduct).Methods("GET")
}

func (a *App) RunServer() {
	fmt.Println("Server started on port ", a.Port)
	log.Fatal(http.ListenAndServe(a.Port, a.Router))
}

func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := GetProducts(a.DB)
	if err != nil {
		// log.Fatal(err.Error())
		fmt.Printf("getProducts error: %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// for _, p := range products {
	// 	fmt.Println("Product: ", p.ID, " ", p.Name, " ", p.Inventory, " ", p.Price)
	// }
	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	// check the request
	vars := mux.Vars(r)
	id := vars["id"]

	var p product
	p.ID, _ = strconv.Atoi(id)
	err := p.GetProduct(a.DB)
	if err != nil {
		// log.Fatal(err.Error())
		fmt.Printf("getProduct error: %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func helloWorld_GET(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World\n")
}

func helloWorld_POST(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World - POST\n")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
