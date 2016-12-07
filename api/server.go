package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	PORT string = ":3000"
)

type Todo struct {
	ID        uint64 `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"isCompleted"`
}

type TodoJSON struct {
	Todo Todo
}

type Todos struct {
	Todos []Todo
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	api := e.Group("/api")
	v1 := api.Group("/v1")

	db, err := gorm.Open("postgres", "host=postgres user=testuser dbname=test sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to db")
	}
	db.AutoMigrate(&Todo{})

	v1.GET("/todos", func(c echo.Context) error {
		todos := make([]Todo, 0)
		db.Find(&todos)
		return c.JSON(http.StatusOK, Todos{Todos: todos})
	})

	v1.GET("/todos/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		var todo Todo
		if db.First(&todo, id).RecordNotFound() {
			return c.String(http.StatusNotFound, "No Todo Found")
		}
		return c.JSON(http.StatusOK, todo)
	})

	v1.POST("/todos", func(c echo.Context) error {
		todoJson := new(TodoJSON)
		if err := c.Bind(todoJson); err != nil {
			return err
		}
		todo := todoJson.Todo
		db.Create(&todo)
		todoJson.Todo = todo
		return c.JSON(http.StatusOK, todoJson)
	})

	v1.PUT("/todos/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		todoJson := new(TodoJSON)
		if err := c.Bind(todoJson); err != nil {
			return err
		}
		var todo Todo
		if db.First(&todo, id).RecordNotFound() {
			return c.String(http.StatusNotFound, "Bad")
		}
		todo.ID = uint64(id)
		todo.Title = todoJson.Todo.Title
		todo.Completed = todoJson.Todo.Completed
		db.Save(&todo)
		todoJson.Todo = todo
		return c.JSON(http.StatusOK, todoJson)
	})

	v1.DELETE("/todos/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		var todo Todo
		if db.First(&todo, id).RecordNotFound() {
			return c.String(http.StatusNotFound, "No Todo Found")
		}
		db.Delete(&todo)
		todoJson := TodoJSON{Todo: todo}
		return c.JSON(http.StatusOK, todoJson)
	})

	e.Logger.Fatal(e.Start(PORT))
}
