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

// ! Swagger documentation @ https://goswagger.io/use/spec.html
func main() {
	a := App{}
	a.Initialize(
		"postgres",
		"postgres",
		"postgres",
	)
	a.Run(":8080")
}
