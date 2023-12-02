package backend_test

import (
	"os"
	"testing"

	"bytes"
	"encoding/json"
	"example.com/backend"
	"log"
	"net/http"
	"net/http/httptest"
)

var a backend.App

const tableProductCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
	id INT NOT NULL PRIMARY KEY AUTOINCREMENT,
	productCode VARCHAR(25) NOT NULL,
	name VARCHAR(256) NOT NULL,
	inventory INT NOT NULL,
	price INT NOT NULL,
	status VARCHAR(64) NOT NULL
)`

const tableOrderCreationQuery = `CREATE TABLE IF NOT EXISTS orders 
(
	id INT NOT NULL PRIMARY KEY AUTOINCREMENT,
	customerName NVARCHAR(50) NOT NULL,
	total INT NOT NULL,
	status NVACHAR(25) NOT NULL
)`

const tableOrderItemsCreationQuery = `CREATE TABLE IF NOT EXISTS order_items 
(
	order_id INT NOT NULL,
	product_id NOT NULL,
	quantity INT NOT NULL,
	FOREIGN KEY (order_id) REFERENCES orders(id),
	FOREIGN KEY (product_id) REFERENCES products(id)
	PRIMARY KEY (order_id, product_id)
)`

func TestMain(m *testing.M) {
	a = backend.App{Port: ":9003", DBType: "sqlite3", DBPath: "./test.sqlite"}
	a.Initialize()
	ensureTableExists()
	code := m.Run()

	clearProductTable()
	clearOrderTable()
	clearOrderItemsTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableProductCreationQuery); err != nil {
		log.Fatal(err)
	}

	if _, err := a.DB.Exec(tableOrderCreationQuery); err != nil {
		log.Fatal(err)
	}

	if _, err := a.DB.Exec(tableOrderItemsCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearProductTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'products'")
}

func insertProductDummy(productCode string, name string, inventory int, price int, status string) {
	a.DB.Exec("INSERT INTO products(productCode, name, inventory, price, status) VALUES(?, ?, ?, ?, ?)", productCode, name, inventory, price, status)
}

func clearOrderTable() {
	a.DB.Exec("DELETE FROM orders")
	a.DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'orders'")
}

func clearOrderItemsTable() {
	a.DB.Exec("DELETE FROM order_items")
}

func insertOrderDummy(customerName string, total int, status string) {
	a.DB.Exec("INSERT INTO orders(customerName, total, status) VALUES(?, ?, ?)", customerName, total, status)
}

func insertOrderItemsDummy(order_id int, product_id int, quantity int) {
	a.DB.Exec("INSERT INTO order_items(order_id, product_id, quantity) VALUES(?, ?, ?)", order_id, product_id, quantity)
}

func TestGetNonExistentProduct(t *testing.T) {
	clearProductTable()

	req, _ := http.NewRequest("GET", "/product/101", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusInternalServerError, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "sql: no rows in result set" {
		t.Errorf("Expected the 'error' key of the response to be set to 'sql:no rows in result set'. Got '%s'", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {
	clearProductTable()

	payload := []byte(`{"productCode":"TEST12345", "name":"ProductTest", "inventory":1, "price":1, "status":"testing"}`)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	verifyProduct(t, m, "TEST12345", "ProductTest", 1, 1, "testing", 1)
}

func TestGetProduct(t *testing.T) {
	clearProductTable()

	insertProductDummy("TEST12345", "ProductTest", 1, 1, "testing")

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	verifyProduct(t, m, "TEST12345", "ProductTest", 1, 1, "testing", 1)
}

func verifyProduct(t *testing.T, m map[string]interface{}, productCode string, name string, inventory int, price int, status string, id int) {
	if m["productCode"] != productCode {
		t.Errorf("Expected productCode to be '%s'. Got '%v'", productCode, m["productCode"])
	}
	if m["name"] != name {
		t.Errorf("Expected name to be '%s'. Got '%v'", name, m["name"])
	}
	if m["inventory"] != float64(inventory) {
		t.Errorf("Expected inventory to be '%v'. Got '%v'", inventory, m["inventory"])
	}
	if m["price"] != float64(price) {
		t.Errorf("Expected price to be '%v'. Got '%v'", price, m["price"])
	}
	if m["status"] != status {
		t.Errorf("Expected status to be '%s'. Got '%v'", status, m["status"])
	}
	if m["id"] != float64(id) {
		t.Errorf("Expected id to be '%v'. Got '%v'", id, m["id"])
	}
}

func TestGetOrder(t *testing.T) {
	clearOrderTable()

	insertOrderDummy("TestCustomer", 132, "testStatus")

	req, _ := http.NewRequest("GET", "/order/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	verifyOrder(t, m, "TestCustomer", 132, "testStatus", 1)
}

func TestCreateOrder(t *testing.T) {
	clearOrderTable()

	payload := []byte(`{"customerName":"TEST12345", "total":32, "status":"testing"}`)

	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	verifyOrder(t, m, "TEST12345", 32, "testing", 1)
}

func TestGetNonExistentOrder(t *testing.T) {
	clearOrderTable()

	req, _ := http.NewRequest("GET", "/order/101", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusInternalServerError, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "sql: no rows in result set" {
		t.Errorf("Expected the 'error' key of the response to be set to 'sql:no rows in result set'. Got '%s'", m["error"])
	}
}

func TestCreateOrderItem(t *testing.T) {
	clearOrderItemsTable()

	payload := []byte(`[{"order_id":132, "product_id":32, "quantity":54}]`)

	req, _ := http.NewRequest("POST", "/orderitems", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m []map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	verifyOrderItem(t, m[0], 132, 32, 54)
}

func verifyOrder(t *testing.T, m map[string]interface{}, customerName string, total int, status string, id int) {
	if m["customerName"] != customerName {
		t.Errorf("Expected customerName to be '%s'. Got '%v'", customerName, m["customerName"])
	}
	if m["total"] != float64(total) {
		t.Errorf("Expected total to be '%v'. Got '%v'", total, m["total"])
	}
	if m["status"] != status {
		t.Errorf("Expected status to be '%s'. Got '%v'", status, m["status"])
	}
	if m["id"] != float64(id) {
		t.Errorf("Expected id to be '%v'. Got '%v'", id, m["id"])
	}
}

func verifyOrderItem(t *testing.T, m map[string]interface{}, order_id int, product_id int, quantity int) {

	if m["order_id"] != float64(order_id) {
		t.Errorf("Expected order_id to be '%v'. Got '%v'", order_id, m["order_id"])
	}
	if m["product_id"] != float64(product_id) {
		t.Errorf("Expected product_id to be '%v'. Got '%v'", product_id, m["product_id"])
	}
	if m["quantity"] != float64(quantity) {
		t.Errorf("Expected quantity to be '%v'. Got '%v'", quantity, m["quantity"])
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
