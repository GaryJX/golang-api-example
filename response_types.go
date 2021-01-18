package main

// Status OK
// swagger:response okResponse
type OkResponse struct {
	// in:body
	Body struct {
	   // HTTP status code 200 - OK
	   Code int `json:"code"`
	}
 }

// Error Not Found
// swagger:response notFoundResponse
type NotFoundResponse struct {
	// in:body
	Body struct {
	   // HTTP status code 404 - Not Found
	   Code int `json:"code"`
	}
 }

// Error Bad Request
// swagger:response badRequestResponse
type badRequestResponse struct {
	// in:body
	Body struct {
	   // HTTP status code 400 -  Bad Request
	   Code int `json:"code"`
	}
 }