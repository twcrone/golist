package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

type Item struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Action string `json:"action"`
}

func repeatHandler(r int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var buffer bytes.Buffer
		for i := 0; i < r; i++ {
			buffer.WriteString("Hello Geeks!\n")
		}
		c.String(http.StatusOK, buffer.String())
	}
}

func listItems(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := db.Exec("CREATE TABLE IF NOT EXISTS items (id serial PRIMARY KEY, name varchar(50), action varchar(10));")
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("error creating database table: %q", err))
			return
		}
		rows, err := db.Query("SELECT id, name, action from items;")
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("error fetching table rows: %q", err))
			return
		}
		defer rows.Close()
		var items []Item
		for rows.Next() {
			var item Item
			err := rows.Scan(&item.Id, &item.Name, &item.Action)
			if err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("error scanning table row: %q", err))
				return
			}
			items = append(items, item)
		}
		c.IndentedJSON(http.StatusOK, items)
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	tStr := os.Getenv("REPEAT")
	repeat, err := strconv.Atoi(tStr)
	if err != nil {
		log.Printf("error converting $REPEAT to an int: %q - Using default\n", err)
		repeat = 5
	}
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("error opening database: %q", err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/old-index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/repeat", repeatHandler(repeat))

	router.GET("/", listItems(db))

	router.GET("/create", func(c *gin.Context) {
		name := c.Request.URL.Query().Get("name")
		c.String(http.StatusOK, "Name is "+name+"\n")
		if _, err := db.Exec("INSERT INTO items (name, action) VALUES ('" + name + "','');"); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating item: %q", err))
			return
		}
	})

	router.GET("/delete", func(c *gin.Context) {
		id := c.Request.URL.Query().Get("id")
		c.String(http.StatusOK, "Deleting player with ID = "+id+"\n")
		if _, err := db.Exec("DELETE FROM players WHERE id=" + id + ";"); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error deleting player: %q", err))
			return
		}
	})

	router.Run(":" + port)
}
