package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
)

var a App
var mock sqlmock.Sqlmock

func (a *App) InitializeMock() {
	var err error
	a.DB, mock, err = sqlmock.New()
	
	if err != nil {
		log.Fatalf("An error '%s' occurred while opening a mock DB connection.", err);
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func TestMain(m *testing.M) {
	a.InitializeMock();
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func TestGetProductsEmptyTable(t *testing.T) {
	clearTable()

	query := "^SELECT[\\s\\S]+FROM products"
	rows := sqlmock.NewRows([]string{"id", "name", "price"})
	mock.ExpectQuery(query).WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/api/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentProduct(t *testing.T) {
	clearTable()

	id := 1
	query := "^SELECT[\\s\\S]+FROM products WHERE[\\s\\S]+"
	rows := sqlmock.NewRows([]string{"id", "name", "price"})
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/api/product/" + fmt.Sprint(id), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got '%s'", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {
	clearTable()

	product := Product{
		ID: 1,
		Name: "test name",
		Price: 11.12,
	}
	query := "^INSERT INTO products[\\s\\S]+VALUES[\\s\\S]+"
	rows := sqlmock.NewRows([]string{"id", "name", "price"}).AddRow(product.ID, product.Name, product.Price)
	mock.ExpectQuery(query).WithArgs(product.Name, product.Price).WillReturnRows(rows)

	var jsonStr = []byte(fmt.Sprintf(`{"name":"%v", "price": %v}`, product.Name, product.Price))
	req, _ := http.NewRequest("POST", "/api/product/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != product.Name {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != product.Price {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	product := Product{
		ID: 1,
		Name: "Product 1",
		Price: 20.0,
	}

	query := "^SELECT[\\s\\S]+FROM products WHERE[\\s\\S]+"
	rows := sqlmock.NewRows([]string{"id", "name", "price"}).AddRow(product.ID, product.Name, product.Price)
	mock.ExpectQuery(query).WithArgs(product.ID).WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/api/product/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	product := Product{
		ID: 1,
		Name: "Product 1",
		Price: 20.0,
	}

	query := "^SELECT[\\s\\S]+FROM products WHERE[\\s\\S]+"
	rows := sqlmock.NewRows([]string{"id", "name", "price"}).AddRow(product.ID, product.Name, product.Price)
	mock.ExpectQuery(query).WithArgs(product.ID).WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/api/product/1", nil)
	response := executeRequest(req)
	var originalProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProduct)

	product = Product{
		ID: 1,
		Name: "New Product 1",
		Price: 20.1,
	}

	query = "^UPDATE products SET[\\s\\S]+WHERE[\\s\\S]+"
	mock.ExpectExec(query).WithArgs(product.Name, product.Price, product.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	var jsonStr = []byte(fmt.Sprintf(`{"name":"%v", "price": %v}`, product.Name, product.Price))

	req, _ = http.NewRequest("PUT", "/api/product/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalProduct["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], m["id"])
	}

	if m["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], m["name"], m["name"])
	}

	if m["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], m["price"], m["price"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	product := Product{
		ID: 1,
		Name: "Product 1",
		Price: 20.0,
	}

	query := "^SELECT[\\s\\S]+FROM products WHERE[\\s\\S]+"
	rows := sqlmock.NewRows([]string{"id", "name", "price"}).AddRow(product.ID, product.Name, product.Price)
	mock.ExpectQuery(query).WithArgs(product.ID).WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/api/product/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	query = "^DELETE FROM products WHERE[\\s\\S]+"
	mock.ExpectExec(query).WithArgs(product.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	req, _ = http.NewRequest("DELETE", "/api/product/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	query = "^SELECT[\\s\\S]+FROM products WHERE[\\s\\S]+"
	rows = sqlmock.NewRows([]string{"id", "name", "price"})
	mock.ExpectQuery(query).WithArgs(product.ID).WillReturnRows(rows)

	req, _ = http.NewRequest("GET", "/api/product/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// Empty products table
func clearTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

// Add count products into table
func addProducts(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO products(name, price) VALUES($1, $2)", "Product "+strconv.Itoa(i), (i+1.0)*10)
	}
}

// Execute test HTTP request
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

// Compare expected and received response code, error if different
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
