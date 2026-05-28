package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"product-api/config"
	"product-api/models"
	"product-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5433")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")

	// Connect to default DB to create test database
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		host, port, user, password,
	)
	adminDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to postgres: " + err.Error())
	}

	adminDB.Exec("DROP DATABASE IF EXISTS products_db_test")
	adminDB.Exec("CREATE DATABASE products_db_test")

	// Connect to test database
	testDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=products_db_test sslmode=disable",
		host, port, user, password,
	)
	testDB, err := gorm.Open(postgres.Open(testDSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}

	if err := testDB.AutoMigrate(&models.Product{}); err != nil {
		panic("failed to migrate test database: " + err.Error())
	}

	config.DB = testDB
	router = routes.SetupRouter()

	code := m.Run()

	// Cleanup: drop test database
	sqlDB, _ := testDB.DB()
	sqlDB.Close()
	adminDB.Exec("DROP DATABASE IF EXISTS products_db_test")

	os.Exit(code)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func cleanProducts(t *testing.T) {
	t.Helper()
	config.DB.Exec("DELETE FROM products")
	config.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

func TestCreateProduct(t *testing.T) {
	cleanProducts(t)

	body := map[string]interface{}{
		"name":        "Test Product",
		"description": "A test product",
		"price":       29.99,
		"quantity":    10,
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Test Product", data["name"])
	assert.Equal(t, "A test product", data["description"])
	assert.Equal(t, 29.99, data["price"])
	assert.Equal(t, float64(10), data["quantity"])
}

func TestCreateProductValidationError(t *testing.T) {
	cleanProducts(t)

	// Missing required name
	body := map[string]interface{}{
		"price":    10.0,
		"quantity": 5,
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateProductNegativePrice(t *testing.T) {
	cleanProducts(t)

	body := map[string]interface{}{
		"name":  "Bad Product",
		"price": -5.0,
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetProducts(t *testing.T) {
	cleanProducts(t)

	// Create two products
	config.DB.Create(&models.Product{Name: "Product 1", Price: 10.0, Quantity: 5})
	config.DB.Create(&models.Product{Name: "Product 2", Price: 20.0, Quantity: 3})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/products", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})
	assert.Len(t, data, 2)
}

func TestGetProductsEmpty(t *testing.T) {
	cleanProducts(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/products", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})
	assert.Len(t, data, 0)
}

func TestGetProduct(t *testing.T) {
	cleanProducts(t)

	product := models.Product{Name: "Single Product", Price: 15.5, Quantity: 7}
	config.DB.Create(&product)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/products/%d", product.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Single Product", data["name"])
	assert.Equal(t, 15.5, data["price"])
}

func TestGetProductNotFound(t *testing.T) {
	cleanProducts(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/products/9999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateProduct(t *testing.T) {
	cleanProducts(t)

	product := models.Product{Name: "Old Name", Price: 10.0, Quantity: 5}
	config.DB.Create(&product)

	body := map[string]interface{}{
		"name":  "New Name",
		"price": 25.0,
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/products/%d", product.ID), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "New Name", data["name"])
	assert.Equal(t, 25.0, data["price"])
}

func TestUpdateProductNotFound(t *testing.T) {
	cleanProducts(t)

	body := map[string]interface{}{"name": "Updated"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/products/9999", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteProduct(t *testing.T) {
	cleanProducts(t)

	product := models.Product{Name: "Delete Me", Price: 5.0, Quantity: 1}
	config.DB.Create(&product)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/products/%d", product.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify soft-deleted: should not be found via normal query
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/products/%d", product.ID), nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusNotFound, w2.Code)

	// Verify record still exists with Unscoped
	var deleted models.Product
	config.DB.Unscoped().First(&deleted, product.ID)
	assert.NotNil(t, deleted.DeletedAt)
}

func TestDeleteProductNotFound(t *testing.T) {
	cleanProducts(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/products/9999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// --- Concurrency Tests ---

func TestConcurrentCreateProducts(t *testing.T) {
	cleanProducts(t)

	const numGoroutines = 50
	var wg sync.WaitGroup
	results := make([]int, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			body := map[string]interface{}{
				"name":     fmt.Sprintf("Concurrent Product %d", idx),
				"price":    float64(idx) + 1.0,
				"quantity": idx,
			}
			jsonBody, _ := json.Marshal(body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			results[idx] = w.Code
		}(i)
	}

	wg.Wait()

	for i, code := range results {
		assert.Equal(t, http.StatusCreated, code, "goroutine %d failed with status %d", i, code)
	}

	// Verify all products were persisted
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/products", nil)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})
	assert.Len(t, data, numGoroutines)
}

func TestConcurrentReadSameProduct(t *testing.T) {
	cleanProducts(t)

	product := models.Product{Name: "Read Target", Price: 42.0, Quantity: 100}
	config.DB.Create(&product)

	const numGoroutines = 100
	var wg sync.WaitGroup
	results := make([]int, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/products/%d", product.ID), nil)
			router.ServeHTTP(w, req)
			results[idx] = w.Code
		}(i)
	}

	wg.Wait()

	for i, code := range results {
		assert.Equal(t, http.StatusOK, code, "read goroutine %d got status %d", i, code)
	}
}

func TestConcurrentUpdateSameProduct(t *testing.T) {
	cleanProducts(t)

	product := models.Product{Name: "Update Target", Price: 10.0, Quantity: 0}
	config.DB.Create(&product)

	const numGoroutines = 30
	var wg sync.WaitGroup
	results := make([]int, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			body := map[string]interface{}{
				"name":     fmt.Sprintf("Updated by %d", idx),
				"price":    float64(idx) + 1.0,
				"quantity": idx,
			}
			jsonBody, _ := json.Marshal(body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/products/%d", product.ID), bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			results[idx] = w.Code
		}(i)
	}

	wg.Wait()

	for i, code := range results {
		assert.Equal(t, http.StatusOK, code, "update goroutine %d got status %d", i, code)
	}

	// Verify product still exists and has a valid state
	var final models.Product
	config.DB.First(&final, product.ID)
	assert.NotEmpty(t, final.Name)
	assert.GreaterOrEqual(t, final.Price, 1.0)
}

func TestConcurrentDeleteSameProduct(t *testing.T) {
	cleanProducts(t)

	product := models.Product{Name: "Delete Race", Price: 5.0, Quantity: 1}
	config.DB.Create(&product)

	const numGoroutines = 20
	var wg sync.WaitGroup
	results := make([]int, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/products/%d", product.ID), nil)
			router.ServeHTTP(w, req)
			results[idx] = w.Code
		}(i)
	}

	wg.Wait()

	successCount := 0
	notFoundCount := 0
	for _, code := range results {
		switch code {
		case http.StatusOK:
			successCount++
		case http.StatusNotFound:
			notFoundCount++
		default:
			t.Errorf("unexpected status code: %d", code)
		}
	}

	// At least one goroutine must have succeeded in deleting
	assert.GreaterOrEqual(t, successCount, 1, "at least one delete should succeed")
	assert.Equal(t, numGoroutines, successCount+notFoundCount, "all responses should be 200 or 404")

	// Verify product is soft-deleted
	var deleted models.Product
	err := config.DB.Unscoped().First(&deleted, product.ID).Error
	assert.NoError(t, err)
	assert.NotNil(t, deleted.DeletedAt)
}

func TestConcurrentMixedOperations(t *testing.T) {
	cleanProducts(t)

	// Seed some products for reads/updates/deletes
	for i := 0; i < 10; i++ {
		config.DB.Create(&models.Product{
			Name:     fmt.Sprintf("Seed Product %d", i),
			Price:    float64(i+1) * 5.0,
			Quantity: i * 2,
		})
	}

	const numGoroutines = 80
	var wg sync.WaitGroup
	errors := make(chan string, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			w := httptest.NewRecorder()

			switch idx % 4 {
			case 0: // Create
				body := map[string]interface{}{
					"name":     fmt.Sprintf("Mixed Create %d", idx),
					"price":    float64(idx) + 1.0,
					"quantity": idx,
				}
				jsonBody, _ := json.Marshal(body)
				req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(w, req)
				if w.Code != http.StatusCreated {
					errors <- fmt.Sprintf("create goroutine %d: expected 201, got %d", idx, w.Code)
				}

			case 1: // Read all
				req, _ := http.NewRequest("GET", "/api/v1/products", nil)
				router.ServeHTTP(w, req)
				if w.Code != http.StatusOK {
					errors <- fmt.Sprintf("list goroutine %d: expected 200, got %d", idx, w.Code)
				}

			case 2: // Read single (may be 200 or 404 depending on timing)
				productID := (idx % 10) + 1
				req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/products/%d", productID), nil)
				router.ServeHTTP(w, req)
				if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
					errors <- fmt.Sprintf("get goroutine %d: unexpected status %d", idx, w.Code)
				}

			case 3: // Update (may be 200 or 404 depending on timing)
				productID := (idx % 10) + 1
				body := map[string]interface{}{
					"name": fmt.Sprintf("Mixed Update %d", idx),
				}
				jsonBody, _ := json.Marshal(body)
				req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/products/%d", productID), bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(w, req)
				if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
					errors <- fmt.Sprintf("update goroutine %d: unexpected status %d", idx, w.Code)
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for errMsg := range errors {
		t.Error(errMsg)
	}
}
