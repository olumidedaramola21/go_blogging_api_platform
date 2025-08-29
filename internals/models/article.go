package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive" // includes mongodb types like ObjectID
)

// Article represents a blog article
type Article struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title string `json:"title" bson:"title" validate:"required,max=200"`
	Content string `json:"content" bson:"content" validate:"required,min=10"`
	Author string `json:"author" bson:"author" validate:"required"`
	Tags []string `json:"tags" bson:"tags"`
	PublishedDate time.Time `json:"publishedDate" bson:"publishedDate"`
	IsPublished bool `json:"isPublished" bson:"isPublished"`
	CreatedAT time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAT time.Time `json:"updatedAt" bson:"updatedAt"`
}

// CreateArticleRequest	represents the request body for creating articles
type CreateArticleRequest struct {
	Title string `json:"title" bson:"title" validate:"required,max=200"`
	Content string `json:"content" bson:"content" validate:"required,min=10"`
	Author string `json:"author" bson:"author" validate:"required"`
	Tags []string `json:"tags" bson:"tags"`	
}

// UpdateArticleRequest	represents the request body for updating articles
type UpdateArticleRequest struct {
	Title *string `json:"title" bson:"title" validate:"required,max=200"`
	Content *string `json:"content" bson:"content" validate:"required,min=10"`
	Author *string `json:"author" bson:"author" validate:"required"`
	Tags []string `json:"tags" bson:"tags"`	
	IsPublished bool `json:"isPublished" bson:"isPublished,omitempty"`
}

// ArticleResponse represents the API response format
type ArticleResponse struct {
	Success bool `json:"success"`
	Message string `json:"message,omitempty"`
	Data interface{} `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

// ArticleListResponse represents the API response format
type ArticleListResponse struct {
	Success bool `json:"success"`
	Data interface{} `json:"data,omitempty"`
	Pagination PaginationInfo `json:"pagination"`
}

type PaginationInfo struct {
	CurrentPage int `json:"currentPage"`
	TotalPages int `json:"totalPages"`
	TotalArticles int64 `json:"totalArticles"`
	HasNextPage bool `json:"hasNextPage"`
	HasPrevPage bool `json:"hasPrevPage"`
}

func NewArticle(req CreateArticleRequest) Article {
	now := time.Now()
	return Article{
		ID: primitive.NewObjectID(),
		Title: req.Title,
		Content: req.Content,
		Author: req.Author,
		Tags: req.Tags,
		PublishedDate: now,
		IsPublished: true,
		CreatedAT: now,
		UpdatedAT: now,
	}
}
