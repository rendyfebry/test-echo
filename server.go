package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	users = []user{}
	seq   = 1
)

func homeHandler(c echo.Context) error {
	data := map[string]string{
		"msg": "Welcome",
	}
	return c.JSON(http.StatusOK, data)
}

func createUser(c echo.Context) error {
	u := new(user)
	u.ID = seq

	if err := c.Bind(u); err != nil {
		return err
	}

	users = append(users, *u)
	seq++
	return c.JSON(http.StatusCreated, u)
}

func getUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	for i := range users {
		if users[i].ID == id {
			return c.JSON(http.StatusOK, users[i])
		}
	}

	return c.JSON(http.StatusNotFound, nil)
}

func getAllUser(c echo.Context) error {
	return c.JSON(http.StatusOK, users)
}

func updateUser(c echo.Context) error {
	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}

	id, _ := strconv.Atoi(c.Param("id"))

	for i := range users {
		if users[i].ID == id {
			users[i].Name = u.Name
			return c.JSON(http.StatusOK, users[i])
		}
	}

	return c.JSON(http.StatusNotFound, nil)
}

func deleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var newUsers []user
	for i := range users {
		if users[i].ID != id {
			newUsers = append(newUsers, users[i])
		}
	}
	users = newUsers

	return c.JSON(http.StatusOK, users)
}

func loginHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "jon" && password == "shhh!" {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "Jon Snow"
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return echo.ErrUnauthorized
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", homeHandler)
	e.POST("/login", loginHandler)

	e.POST("/users", createUser)
	e.GET("/users", getAllUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", restricted)

	e.Logger.Fatal(e.Start(":9000"))
}
