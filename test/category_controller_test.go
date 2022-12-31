package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mitchellh/mapstructure"
	"io"
	"learn-golang-restful-api/app"
	"learn-golang-restful-api/controller"
	"learn-golang-restful-api/helper"
	"learn-golang-restful-api/middleware"
	"learn-golang-restful-api/model/domain"
	"learn-golang-restful-api/model/web"
	"learn-golang-restful-api/repository"
	"learn-golang-restful-api/service"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/belajar_golang_restful_api_test")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func setupRouter(db *sql.DB) http.Handler {
	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)
	router := app.NewRouter(categoryController)

	return middleware.NewAuthMiddleware(router)
}

func truncateCategory(db *sql.DB) {
	_, _ = db.Exec("TRUNCATE category")
}

func TestCreateCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name": "Gadget"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	var categoryResponse web.CategoryResponse
	_ = json.Unmarshal(body, &responseBody)
	_ = mapstructure.Decode(responseBody.Data, &categoryResponse)

	assert.Equal(t, http.StatusOK, responseBody.Code)
	assert.Equal(t, "OK", responseBody.Status)
	assert.Equal(t, "Gadget", categoryResponse.Name)
}

func TestCreateCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name": ""}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3001/api/categories", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	var categoryResponse web.CategoryResponse
	_ = json.Unmarshal(body, &responseBody)
	_ = mapstructure.Decode(responseBody.Data, &categoryResponse)

	assert.Equal(t, http.StatusBadRequest, responseBody.Code)
	assert.Equal(t, "BAD REQUEST", responseBody.Status)
}

func TestUpdateCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	_ = tx.Commit()

	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name": "Fashion"}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	var categoryResponse web.CategoryResponse
	_ = json.Unmarshal(body, &responseBody)
	_ = mapstructure.Decode(responseBody.Data, &categoryResponse)

	assert.Equal(t, http.StatusOK, responseBody.Code)
	assert.Equal(t, "OK", responseBody.Status)
	assert.Equal(t, category.Id, categoryResponse.Id)
	assert.Equal(t, "Fashion", categoryResponse.Name)
}

func TestUpdateCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	_ = tx.Commit()

	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name": ""}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	var categoryResponse web.CategoryResponse
	_ = json.Unmarshal(body, &responseBody)
	_ = mapstructure.Decode(responseBody.Data, &categoryResponse)

	assert.Equal(t, http.StatusBadRequest, responseBody.Code)
	assert.Equal(t, "BAD REQUEST", responseBody.Status)
}

func TestGetCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	_ = tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	var categoryResponse web.CategoryResponse
	_ = json.Unmarshal(body, &responseBody)
	_ = mapstructure.Decode(responseBody.Data, &categoryResponse)

	assert.Equal(t, http.StatusOK, responseBody.Code)
	assert.Equal(t, "OK", responseBody.Status)
	assert.Equal(t, category.Id, categoryResponse.Id)
	assert.Equal(t, category.Name, categoryResponse.Name)
}

func TestGetCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories/404", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	var categoryResponse web.CategoryResponse
	_ = json.Unmarshal(body, &responseBody)
	_ = mapstructure.Decode(responseBody.Data, &categoryResponse)

	assert.Equal(t, http.StatusNotFound, responseBody.Code)
	assert.Equal(t, "NOT FOUND", responseBody.Status)
}

func TestDeleteCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	_ = tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	var categoryResponse web.CategoryResponse
	_ = json.Unmarshal(body, &responseBody)
	_ = mapstructure.Decode(responseBody.Data, &categoryResponse)

	assert.Equal(t, http.StatusOK, responseBody.Code)
	assert.Equal(t, "OK", responseBody.Status)
}

func TestDeleteCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/404", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	var categoryResponse web.CategoryResponse
	_ = json.Unmarshal(body, &responseBody)
	_ = mapstructure.Decode(responseBody.Data, &categoryResponse)

	assert.Equal(t, http.StatusNotFound, responseBody.Code)
	assert.Equal(t, "NOT FOUND", responseBody.Status)
}

func TestListCategoriesSuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category0 := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	category1 := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Laptop",
	})
	_ = tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	var categoriesResponse []web.CategoryResponse
	_ = json.Unmarshal(body, &responseBody)
	_ = mapstructure.Decode(responseBody.Data, &categoriesResponse)
	categoryResponse0 := categoriesResponse[0]
	categoryResponse1 := categoriesResponse[1]

	assert.Equal(t, http.StatusOK, responseBody.Code)
	assert.Equal(t, "OK", responseBody.Status)
	assert.Equal(t, category0.Id, categoryResponse0.Id)
	assert.Equal(t, category0.Name, categoryResponse0.Name)
	assert.Equal(t, category1.Id, categoryResponse1.Id)
	assert.Equal(t, category1.Name, categoryResponse1.Name)
}

func TestUnauthorized(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories", nil)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	_ = json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusUnauthorized, responseBody.Code)
	assert.Equal(t, "UNAUTHORIZED", responseBody.Status)
}
