package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func homeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Home!")
}

func userHandler(c echo.Context) error {
	return c.String(http.StatusOK, "User!")
}

func main() {
	e := echo.New()

	loggerConfig := middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}

	e.Use(middleware.LoggerWithConfig(loggerConfig))

	e.GET("/", homeHandler)
	e.GET("/users", userHandler)

	e.Logger.Fatal(e.Start(":9000"))
}
