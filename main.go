// Golang Example API
//
// Example Description
//
// Terms Of Service:
//
//     Schemes: http, https
//     Host: localhost:8080
//     Version: 1.0.0
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

// App contains Router and Database connection
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) initializeDBTables() {
	// Initialize products table if it does not already exist
	var tableCreationQuery = `CREATE TABLE IF NOT EXISTS products
	(
		id SERIAL,
		name TEXT NOT NULL,
		price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
		CONSTRAINT products_pkey PRIMARY KEY (id)
	)`

	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func (a *App) initializeRoutes() {
	baseURL := "/api"

	fs := http.FileServer(http.Dir("./swaggerui"))
	a.Router.PathPrefix(baseURL).Handler(http.StripPrefix("/swaggerui/", fs))

	// TODO: Move this to products.go

	// swagger:operation GET /products products getProducts
	// ---
	// summary: Get all products.
	// description: Returns all available products.
	// responses:
	//  200:
	//    $ref: "#/responses/productsResponse"
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}

// Initialize sets up the database connection and router
func (a *App) Initialize(user, password, dbname string) {
	// Postgres connection string for Heroku server
	connectionString := os.Getenv("DATABASE_URL")

	// Postgres connection string for local development
	if connectionString == "" {
		connectionString = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	}

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.initializeDBTables()

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run starts up the server
func (a *App) Run(port string) {
	log.Fatal(http.ListenAndServe(port, a.Router))
}

func main() {
	port := os.Getenv("PORT")

	// Default port in development
	if port == "" {
		port = "8080"
	}

	a := App{}
	a.Initialize(
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	a.Run(":" + port)
}
