package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"github.com/thampaponn/go-fellowie/controller"
	_ "github.com/thampaponn/go-fellowie/models"
	"github.com/thampaponn/go-fellowie/repository"
)

func main() {
	_, err := repository.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}
	defer repository.CloseDB()

	// Initialize Echo framework
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/users", controller.CreateUser)
	e.GET("/users", controller.GetUsers)
	e.PATCH("/users/:id", controller.UpdateUser)
	e.DELETE("/users/:id", controller.DeleteUser)

	e.Logger.Fatal(e.Start(":8080"))
}
