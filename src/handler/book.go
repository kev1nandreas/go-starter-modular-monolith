package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/kev1nandreas/go-rest-api-template/src/modules/book/dto"
	"github.com/kev1nandreas/go-rest-api-template/src/modules/book/usecase"
	"github.com/kev1nandreas/go-rest-api-template/src/shared/response"
)

type BookHandler struct {
	usecase     usecase.BookUsecase
	redisClient *redis.Client
	ctx         context.Context
}

func NewBookHandler(uc usecase.BookUsecase, redisClient *redis.Client, ctx context.Context) *BookHandler {
	return &BookHandler{
		usecase:     uc,
		redisClient: redisClient,
		ctx:         ctx,
	}
}

// @BasePath /api/v1

// Healthcheck godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} ok
// @Router / [get]
func (h *BookHandler) Healthcheck(c *gin.Context) {
	response.NewSuccessResponse("ok", nil).Send(c)
}

// FindBooks godoc
// @Summary Get all books with pagination
// @Description Get a list of all books with optional pagination
// @Tags books
// @Security ApiKeyAuth
// @Produce json
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {array} entity.Book "Successfully retrieved list of books"
// @Router /books [get]
func (h *BookHandler) FindBooks(c *gin.Context) {
	offsetQuery := c.DefaultQuery("offset", "0")
	limitQuery := c.DefaultQuery("limit", "10")

	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		response.NewErrorResponse(http.StatusBadRequest, "Invalid offset format", err.Error()).Send(c)
		return
	}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		response.NewErrorResponse(http.StatusBadRequest, "Invalid limit format", err.Error()).Send(c)
		return
	}

	// Try cache first
	cacheKey := "books_offset_" + offsetQuery + "_limit_" + limitQuery
	cachedBooks, err := h.redisClient.Get(h.ctx, cacheKey).Result()
	if err == nil {
		var books interface{}
		if err := json.Unmarshal([]byte(cachedBooks), &books); err != nil {
			response.NewErrorResponse(http.StatusInternalServerError, "Failed to unmarshal cached data", err.Error()).Send(c)
			return
		}
		response.NewSuccessResponse("Books retrieved from cache", books).Send(c)
		return
	}

	books, err := h.usecase.FindAll(offset, limit)
	if err != nil {
		response.NewErrorResponse(http.StatusInternalServerError, "Failed to retrieve books", err.Error()).Send(c)
		return
	}

	// Cache the result
	serializedBooks, err := json.Marshal(books)
	if err == nil {
		h.redisClient.Set(h.ctx, cacheKey, serializedBooks, time.Minute)
	}

	response.NewSuccessResponse("Books retrieved successfully", books).Send(c)
}

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book with the given input data
// @Tags books
// @Security ApiKeyAuth
// @Security JwtAuth
// @Accept  json
// @Produce  json
// @Param   input     body   dto.CreateBookRequest   true   "Create book object"
// @Success 201 {object} entity.Book "Successfully created book"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	var input dto.CreateBookRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		response.NewErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()).Send(c)
		return
	}

	book, err := h.usecase.Create(input)
	if err != nil {
		response.NewErrorResponse(http.StatusInternalServerError, "Failed to create book", err.Error()).Send(c)
		return
	}

	// Invalidate cache
	keysPattern := "books_offset_*"
	keys, err := h.redisClient.Keys(h.ctx, keysPattern).Result()
	if err == nil {
		for _, key := range keys {
			h.redisClient.Del(h.ctx, key)
		}
	}

	response.Response{
		StatusCode: http.StatusCreated,
		Success:    true,
		Message:    "Book created successfully",
		Data:       book,
	}.Send(c)
}

// FindBook godoc
// @Summary Find a book by ID
// @Description Get details of a book by its ID
// @Tags books
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} entity.Book "Successfully retrieved book"
// @Failure 404 {string} string "Book not found"
// @Router /books/{id} [get]
func (h *BookHandler) FindBook(c *gin.Context) {
	book, err := h.usecase.FindByID(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(http.StatusNotFound, "Book not found", err.Error()).Send(c)
		return
	}

	response.NewSuccessResponse("Book retrieved successfully", book).Send(c)
}

// UpdateBook godoc
// @Summary Update a book by ID
// @Description Update the book details for the given ID
// @Tags books
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param id path string true "Book ID"
// @Param input body dto.UpdateBookRequest true "Update book object"
// @Success 200 {object} entity.Book "Successfully updated book"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "book not found"
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	var input dto.UpdateBookRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		response.NewErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()).Send(c)
		return
	}

	book, err := h.usecase.Update(c.Param("id"), input)
	if err != nil {
		response.NewErrorResponse(http.StatusNotFound, "Book not found", err.Error()).Send(c)
		return
	}

	response.NewSuccessResponse("Book updated successfully", book).Send(c)
}

// DeleteBook godoc
// @Summary Delete a book by ID
// @Description Delete the book with the given ID
// @Tags books
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Book ID"
// @Success 204 {string} string "Successfully deleted book"
// @Failure 404 {string} string "book not found"
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
	if err := h.usecase.Delete(c.Param("id")); err != nil {
		response.NewErrorResponse(http.StatusNotFound, "Book not found", err.Error()).Send(c)
		return
	}

	response.Response{
		StatusCode: http.StatusNoContent,
		Success:    true,
		Message:    "Book deleted successfully",
		Data:       true,
	}.Send(c)
}
