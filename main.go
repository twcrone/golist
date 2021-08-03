package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
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
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("error opening database: %q", err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.GET("/items", listItems(db))

	router.POST("/items", func(c *gin.Context) {
		var newItem Item

		if err := c.BindJSON(&newItem); err != nil {
			return
		}

		if newItem.Name != "" {
			_, err := db.Exec("INSERT INTO items (name, action) VALUES ('" + newItem.Name + "','');")
			if err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error creating item: %q", err))
			}
			c.IndentedJSON(http.StatusCreated, newItem)
		} else if newItem.Id > 0 && newItem.Action != "" {
			_, err := db.Exec("UPDATE items set action = ? where id = ?", newItem.Action, newItem.Id)
			if err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error updating item: %q", err))
			}
			c.IndentedJSON(http.StatusOK, newItem)
		}
	})

	router.DELETE("/items", func(c *gin.Context) {
		c.String(http.StatusOK, "Deleting actioned items\n")
		if _, err := db.Exec("DELETE FROM items WHERE action NOT NULL;"); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error deleting actioned items: %q", err))
			return
		}
	})

	router.Run(":" + port)
}
