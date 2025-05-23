package main

import (
	"net/http"
	// run go get . if this package is not installed
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(tokenAuthMiddleware())
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.DELETE("/albums/:id", deleteAlbumByID)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run("0.0.0.0:" + port)
}

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 20.00},
	{ID: "3", Title: "Sarah Vaughan", Artist: "Sarah Vaughan", Price: 49.00},
}

// getAlbums responds with the list of all albums as JSON
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums to add a new album
func postAlbums(c *gin.Context) {
	var newAlbum album

	// bind the request body to newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID -> fetch a specific album by its ID
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func deleteAlbumByID(c *gin.Context) {
	id := c.Param("id")

	for i, a := range albums {
		if a.ID == id {
			albums = append(albums[:i], albums[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "album deleted"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func tokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := os.Getenv("API_SECRET_TOKEN")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "server misconfiguration"})
			return
		}
		authHeader := c.GetHeader("Authorization")
		if authHeader != token {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
