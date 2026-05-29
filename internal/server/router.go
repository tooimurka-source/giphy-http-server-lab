package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GIFSearcher interface {
	Search(query string) (string, error)
}

type searchRequest struct {
	Query string `json:"query" binding:"required"`
}

func NewRouter(searcher GIFSearcher) *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "giphy-http-server",
			"routes": []string{
				"GET /health",
				"GET /api/gifs/search?q=cat",
				"POST /api/gifs/search",
			},
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	{
		api.GET("/gifs/search", func(c *gin.Context) {
			query := strings.TrimSpace(c.Query("q"))
			respondWithGIF(c, searcher, query)
		})

		api.POST("/gifs/search", func(c *gin.Context) {
			var request searchRequest
			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "query is required"})
				return
			}

			respondWithGIF(c, searcher, strings.TrimSpace(request.Query))
		})
	}

	return router
}

func respondWithGIF(c *gin.Context, searcher GIFSearcher, query string) {
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query is required"})
		return
	}

	gifURL, err := searcher.Search(query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "gif not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query": query,
		"url":   gifURL,
	})
}
