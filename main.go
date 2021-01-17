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
	"fmt"
	"os"
)

// ! Swagger documentation @ https://goswagger.io/use/spec.html
func main() {
	port := os.Getenv("PORT")

	fmt.Println("HELLO WORLD")
	if port == "" {
		port = "8080"
	}

	a := App{}
	a.Initialize(
		"postgres",
		"postgres",
		"postgres",
	)
	a.Run(":" + port)
}
