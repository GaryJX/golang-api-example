# Golang API Example by GaryJX

Live example deployed at [https://garyjx-golang-api-example.herokuapp.com/api/](https://garyjx-golang-api-example.herokuapp.com/api/)

## Local Development Commands

1. Clone the repo: `git clone https://github.com/GaryJX/golang-api-example.git`
2. Configure environment variables: `cp .env.sample .env`
3. Run tests: `go test`
4. Re-generate Swagger API docs: `swagger generate spec -o ./api/swagger.json`

## Running Server Locally

1. Ensure that you have PostgreSQL installed on your machine. Create 2 empty databases (one for playground and one for running tests) and configure the environment variables correctly for your postgres connection.
2. For Mac & Linux users:

```
go run ./main.go
```

3. For Windows users:

```
go build
.\golang-api-example.exe
```

## Resources

- Swagger UI - cloned from https://github.com/swagger-api/swagger-ui/tree/master/dist
- Go Swagger Syntax & Documentation: https://goswagger.io/
