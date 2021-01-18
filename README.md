# Golang API Example by GaryJX

Live example deployed at [https://garyjx-golang-api-example.herokuapp.com/api/](https://garyjx-golang-api-example.herokuapp.com/api/)

## Local Development Commands

1. Ensure that you have Golang installed on your machine (This repo uses v1.15.6).
2. Ensure that you have PostgreSQL installed on your machine. Create 2 empty databases (one for playground and one for running tests).
3. Clone the repo: `git clone https://github.com/GaryJX/golang-api-example.git`
4. Change directory: `cd golang-api-example`
5. Configure environment variables for your postgres connection: `cp .env.sample .env`
6. Run tests: `go test -v`
7. Install `go-swagger` [here](https://goswagger.io/install.html) or run:

```
go get -u github.com/go-swagger/go-swagger/cmd/swagger
```

8. Re-generate Swagger API docs: `swagger generate spec -o ./api/swagger.json`
9. Start server:

```
go build
./golang-api-example.exe
```

10. View API at http://localhost:8080/api/

## Resources

- Swagger UI - cloned from https://github.com/swagger-api/swagger-ui/tree/master/dist
- Go Swagger Syntax & Documentation: https://goswagger.io/
