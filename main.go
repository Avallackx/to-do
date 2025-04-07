package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Task struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var db *gorm.DB

func initDB() {
	dsn := "host=postgres user=postgres password=admin dbname=db port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	db.AutoMigrate(&Task{})
}

func getTasks(c *gin.Context) {
	var tasks []Task
	db.Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func createTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&task)
	c.JSON(http.StatusCreated, task)
}

func updateTask(c *gin.Context) {
	var task Task
	id := c.Param("id")

	db.First(&task, id)
	if task.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&task)
	c.JSON(http.StatusOK, task)
}

func deleteTask(c *gin.Context) {
	var task Task
	id := c.Param("id")

	db.First(&task, id)
	if task.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	db.Delete(&task)
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

func main() {
	initDB()
	r := gin.Default()

	r.GET("/tasks", getTasks)
	r.POST("/tasks", createTask)
	r.PUT("/tasks/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)

	r.Run(":8080")
}
