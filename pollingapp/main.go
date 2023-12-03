// main.go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

type Poll struct {
	gorm.Model
	Question     string
	LikeCount    int
	DislikeCount int
}

func main() {
	r := gin.Default()

	// Initialize the database
	initDB()

	// Routes
	r.GET("/polls", getPolls)
	r.GET("/polls/:id", getPoll)
	r.POST("/polls", createPoll)
	r.PUT("/polls/:id/like", likePoll)
	r.PUT("/polls/:id/dislike", dislikePoll)

	// Run the server
	r.Run(":8080")
}

func initDB() {
	var err error
	db, err = gorm.Open("sqlite3", "polls.db")
	if err != nil {
		panic("Failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&Poll{})
}

func getPolls(c *gin.Context) {
	var polls []Poll
	db.Find(&polls)
	c.JSON(200, polls)
}

func getPoll(c *gin.Context) {
	var poll Poll
	id := c.Param("id")

	if err := db.First(&poll, id).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, poll)
}

func createPoll(c *gin.Context) {
	var input struct {
		Question string `json:"question" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	poll := Poll{
		Question: input.Question,
	}

	db.Create(&poll)
	c.JSON(201, poll)
}

func likePoll(c *gin.Context) {
	var poll Poll
	id := c.Param("id")

	if err := db.First(&poll, id).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	poll.LikeCount++
	db.Save(&poll)

	c.JSON(200, poll)
}

func dislikePoll(c *gin.Context) {
	var poll Poll
	id := c.Param("id")

	if err := db.First(&poll, id).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	poll.DislikeCount++
	db.Save(&poll)

	c.JSON(200, poll)
}
