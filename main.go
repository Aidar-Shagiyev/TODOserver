package main

import (
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type TodoItem struct {
	gorm.Model
	Deadline string
	Action   string
}

func main() {
	db, err := gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	db.AutoMigrate(&TodoItem{})

	router := gin.Default()
	router.GET("/", GetIndex(db, router))
	router.POST("/delete", DeleteTodos(db))
	router.POST("/add", AddTodos(db))
	router.Run(":8888")
}

func GetIndex(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	router.LoadHTMLGlob("templates/*")

	return func(c *gin.Context) {
		var todoitems []TodoItem
		db.Order("deadline").Find(&todoitems)
		c.HTML(http.StatusOK, "index.html", gin.H{"todolist": todoitems})
	}
}

func DeleteTodos(db *gorm.DB) gin.HandlerFunc {
	type DeleteForm struct {
		DeleteList []int `form:"deletelist[]"`
	}

	return func(c *gin.Context) {
		var deleteform DeleteForm
		c.ShouldBind(&deleteform)
		db.Delete(&TodoItem{}, deleteform.DeleteList)
		c.Redirect(http.StatusFound, "/")
	}
}

func AddTodos(db *gorm.DB) gin.HandlerFunc {
	type AddForm struct {
		Deadline string `form:"deadline"`
		Action   string `form:"action"`
	}

	return func(c *gin.Context) {
		var addform AddForm
		c.ShouldBind(&addform)
		_, err := time.Parse("2006-01-02", addform.Deadline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format."})
			return
		}
		db.Create(&TodoItem{Deadline: addform.Deadline, Action: addform.Action})
		c.Redirect(http.StatusFound, "/")
	}
}
