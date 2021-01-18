// Golang Example API
//
// Example Description
//
// Terms Of Service:
//
//     BasePath: /api
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
	// swagger:operation GET /products products getProducts
	// ---
	// summary: Get all products.
	// description: Returns all available products.
	// responses:
	//  200:
	//    $ref: "#/responses/productsResponse"
	a.Router.HandleFunc("/api/products", a.getProducts).Methods("GET")

	// swagger:operation GET /product/{id} products getProduct
	// ---
	// summary: Get a product by id.
	// description: Return a product provided by the id.
	// parameters:
	// - name: id
	//   in: path
	//   description: Product ID
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/productResponse"
	//   "404":
	//     "$ref": "#/responses/notFoundResponse"
	a.Router.HandleFunc("/api/product/{id:[0-9]+}", a.getProduct).Methods("GET")

	// swagger:operation POST /product/ products createProduct
	// ---
	// summary: Creates a new product.
	// description: Create a new product.
	// parameters:
	// - name: product
	//   description: Product
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/Product"
	// responses:
	//   "200":
	//     "$ref": "#/responses/productResponse"
	//   "400":
	//     "$ref": "#/responses/badRequestResponse"
	a.Router.HandleFunc("/api/product/", a.createProduct).Methods("POST")

	// swagger:operation PUT /product/{id} products updateProduct
	// ---
	// summary: Update a product by ID.
	// description: Update a product by ID.
	// parameters:
	// - name: id
	//   in: path
	//   description: Product ID
	//   type: string
	//   required: true
	// - name: product
	//   description: Product
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/Product"
	// responses:
	//   "200":
	//     "$ref": "#/responses/productResponse"
	//   "400":
	//     "$ref": "#/responses/badRequestResponse"
	//   "404":
	//     "$ref": "#/responses/notFoundResponse"
	a.Router.HandleFunc("/api/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")

	// swagger:operation DELETE /product/{id} products deleteProduct
	// ---
	// summary: Delete a product by ID.
	// description: Delete a product by ID.
	// parameters:
	// - name: id
	//   in: path
	//   description: Product ID
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/okResponse"
	a.Router.HandleFunc("/api/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
	
	// Serve Swagger API Docs
	a.Router.PathPrefix("/api").Handler(http.StripPrefix("/api/", http.FileServer(http.Dir("./api"))))
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
