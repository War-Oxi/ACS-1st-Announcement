package main

import (
	"net"
	"net/http"
	"strconv"
	"time"
	// "strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Post entity represents a post in the database
type Post struct {
	// gorm.Model
	PostID       uint      `gorm:"primary_key;autoIncrement" json:"post_id"`
	PostTitle    string    `json:"post_title"`
	PostContents string    `json:"post_contents"`
	PostCreateAt time.Time `json:"post_create_at"`
}

var db *gorm.DB

func init() {
	// Initialize the database connection
	var err error
	dsn := "acs_master:1234@tcp(211.183.3.100:33306)/acs_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Post{})
}

func createPost(c *gin.Context) {
	var newPost Post
	if err := c.BindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPost.PostCreateAt = time.Now()
	if result := db.Create(&newPost); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"post_id": newPost.PostID})
}

func getPosts(c *gin.Context) {
	var posts []Post
	var totalCount int64
	page := c.DefaultQuery("page", "1")

	// Calculate the total number of posts
	db.Model(&Post{}).Count(&totalCount)

	// Parse the page number
	var pageNumber int
	if num, err := strconv.Atoi(page); err == nil && num > 0 {
		pageNumber = num
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	// Define the pagination variables
	offset := (pageNumber - 1) * 15
	limit := 15

	if result := db.Order("post_create_at desc").Offset(offset).Limit(limit).Find(&posts); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Respond with posts and total post count
	c.JSON(http.StatusOK, gin.H{"posts": posts, "total_count": totalCount})
}

func getContents(c *gin.Context) {
	var post Post
	if err := db.Where("post_id = ?", c.Param("id")).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func deleteAllPosts(c *gin.Context) {
	// Delete all records from the Post table
	// result := db.Where("1 = 1").Delete(&Post{})
	result := db.Exec("TRUNCATE TABLE posts")

	if result.Error != nil {
		// If there was an error during the deletion
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Respond with the number of deleted records
	c.JSON(http.StatusOK, gin.H{"deleted": result.RowsAffected})
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func ipHandler(c *gin.Context) {
	ip := getLocalIP()
	if ip == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get IP"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ip": ip})
}

func main2() {
	router := gin.Default()

	// Set CORS middleware
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // 여기서 모든 오리진을 허용하도록 설정합니다.
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Set the router to serve static files from a specified directory
	router.Static("/frontend", "./frontend")

	// Serve HTML directly by registering a handler for the index route
	router.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

	router.GET("/ip", ipHandler)

	// Define your routes here
	router.DELETE("/posts", deleteAllPosts)
	router.POST("/posts", createPost)
	router.GET("/posts", getPosts)
	router.GET("/posts/:id", getContents)

	// Start the server
	router.Run(":8989")
}
