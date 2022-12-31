package main

import (
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"learn-golang-restful-api/app"
	"learn-golang-restful-api/controller"
	"learn-golang-restful-api/helper"
	"learn-golang-restful-api/middleware"
	"learn-golang-restful-api/repository"
	"learn-golang-restful-api/service"
	"net/http"
)

func main() {

	db := app.NewDB()
	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)
	router := app.NewRouter(categoryController)

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: middleware.NewAuthMiddleware(router),
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
