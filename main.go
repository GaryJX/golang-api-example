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
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	Router *mux.Router
	DB     *mongo.Database
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
	a.Router.HandleFunc("/api/product/{id}", a.getProduct).Methods("GET")

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
	a.Router.HandleFunc("/api/product/{id}", a.updateProduct).Methods("PUT")

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
	a.Router.HandleFunc("/api/product/{id}", a.deleteProduct).Methods("DELETE")

	// Serve Swagger API Docs
	a.Router.PathPrefix("/api").Handler(http.StripPrefix("/api/", http.FileServer(http.Dir("./api"))))
}

// Initialize sets up the database connection and router
func (a *App) Initialize(connectionURI string) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	a.DB = client.Database("example-database")
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(port string) {
	fmt.Println("Running server on port" + port)
	log.Fatal(http.ListenAndServe(port, a.Router))
}

func main() {
	port := os.Getenv("PORT")

	// Default port in development
	if port == "" {
		port = "8080"
	}

	a := App{}
	a.Initialize("mongodb://localhost:27017")
	a.Run(":" + port)
}
