package main

import (
	"log"
	"net/http"

	// "sample/route"
	"sample/pkg/db"

	"github.com/labstack/echo/middleware"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type tasks struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	db := db.Connect()
	defer db.Close()

	initDb()

	e := echo.New()

	// 全てのリクエストで差し込みたいミドルウェア（ログとか）はここ
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/tasks/create", taskCreate)
	e.GET("/tasks/read", read)
	e.PUT("/tasks/update", taskUpdate)
	e.DELETE("/tasks/delete", taskDelete)

	e.Start("localhost:1323")
}

func taskCreate(c echo.Context) error {

	db := db.Connect()
	defer db.Close()

	name := c.QueryParam("name")

	db.Query(`INSERT INTO tasks (name) VALUES (?)`, name)
	return nil
}

func read(c echo.Context) error {

	db := db.Connect()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		log.Fatal(err)
	}

	tasksArray := make([]tasks, 0)
	for rows.Next() {
		var tasks tasks
		err := rows.Scan(&tasks.ID, &tasks.Name)
		if err != nil {
			log.Fatal(err)
		}
		tasksArray = append(tasksArray, tasks)
	}
	return c.JSON(http.StatusOK, tasksArray)
}

func taskUpdate(c echo.Context) error {

	db := db.Connect()
	defer db.Close()

	id := c.QueryParam("id")
	name := c.QueryParam("name")

	db.Query(`UPDATE tasks set name = ? WHERE id = ?`, name, id)
	return nil
}

func taskDelete(c echo.Context) error {

	db := db.Connect()
	defer db.Close()

	id := c.QueryParam("id")

	db.Query(`DELETE FROM tasks WHERE id = (?)`, id)
	return nil
}

// データベース初期化
func initDb() {
	db := db.Connect()
	defer db.Close()
	db.Query(`CREATE TABLE IF NOT EXISTS tasks (
				id SERIAL NOT NULL PRIMARY KEY,
				name VARCHAR(255)
			)`)
	db.Query(`TRUNCATE TABLE tasks`)
	db.Query(`INSERT INTO tasks (name) VALUES ('Task1')`)
	db.Query(`INSERT INTO tasks (name) VALUES ('Task2')`)
	db.Query(`INSERT INTO tasks (name) VALUES ('Task3')`)
}
