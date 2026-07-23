package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kev1nandreas/go-rest-api-template/src/modules/user/dto"
	"github.com/kev1nandreas/go-rest-api-template/src/modules/user/usecase"
	"github.com/kev1nandreas/go-rest-api-template/src/shared/response"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: uc,
	}
}

// @BasePath /api/v1

// LoginHandler godoc
// @Summary Authenticate a user
// @Schemes
// @Description Authenticates a user using username and password, returns a JWT token if successful
// @Tags user
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   user     body    dto.LoginRequest     true        "User login object"
// @Success 200 {string} string "JWT Token"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /login [post]
func (h *UserHandler) LoginHandler(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()).Send(c)
		return
	}

	loginResp, err := h.usecase.Login(req)
	if err != nil {
		response.NewErrorResponse(http.StatusUnauthorized, "Invalid username or password", err.Error()).Send(c)
		return
	}

	response.NewSuccessResponse("Login successful", gin.H{"token": loginResp.Token}).Send(c)
}

// RegisterHandler godoc
// @Summary Register a new user
// @Schemes http
// @Description Registers a new user with the given username and password
// @Tags user
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   user     body    dto.RegisterRequest     true        "User registration object"
// @Success 201 {string} string	"Successfully registered"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /register [post]
func (h *UserHandler) RegisterHandler(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()).Send(c)
		return
	}

	if err := h.usecase.Register(req); err != nil {
		response.NewErrorResponse(http.StatusInternalServerError, "Could not register user", err.Error()).Send(c)
		return
	}

	response.Response{
		StatusCode: http.StatusCreated,
		Success:    true,
		Message:    "Registration successful",
		Data:       nil,
	}.Send(c)
}
