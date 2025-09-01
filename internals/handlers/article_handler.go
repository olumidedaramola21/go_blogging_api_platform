package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go_blog_api/internals/database"
	"go_blog_api/internals/middleware"
	"go_blog_api/internals/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ArticleHandler handles article-related HTTP request
type ArticleHandler struct {
	collection *mongo.Collection
}

// NewArticleHandler creates a new article handler
func NewArticleHandler(databaseName, collectionName string) *ArticleHandler {
	return &ArticleHandler{
		collection: database.GetCollection(databaseName, collectionName),
	}
}

// GetArticles handles GET	/articles - retrieves all articles with filtering
func (h *ArticleHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse query parameters
	query := r.URL.Query()
	tags := query.Get("tags")
	author := query.Get("author")
	limitStr := query.Get("limit")
	pageStr := query.Get("page")

	// set default values
	limit := int64(10)
	page := int64(1)

	// Parse limit and page
	if limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 64); err == nil && l > 0 {
			limit = l
		}
	}
	if pageStr != "" {
		if p, err := strconv.ParseInt(pageStr, 10, 64); err == nil && p > 0 {
			page = p
		}
	}

	// Build filter
	filter := bson.M("isPublished": true)

	if tags != "" {
		tagList := strings.Split(tags, ",")
		// Trim whitespace from each tag
		for i, tag := range tagList {
			tagList[i] = strings.TrimSpace(tag)
		}
	}
}
