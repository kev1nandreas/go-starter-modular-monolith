package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kev1nandreas/go-rest-api-template/config"
	"github.com/kev1nandreas/go-rest-api-template/src/factory"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	docs "github.com/kev1nandreas/go-rest-api-template/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(
	logger *zap.Logger,
	mongoCollection *mongo.Collection,
	bookHandler *BookHandler,
	userHandler *UserHandler,
) *gin.Engine {
	isLogging := config.GetEnvBool("APP_MONGO_LOGGING", false)

	r := gin.Default()

	if isLogging {
		logger.Info("MongoDB logging is enabled")
		r.Use(factory.Logger(logger, mongoCollection))
	} else {
		logger.Info("MongoDB logging is disabled")
		r.Use(gin.Logger())
	}

	if gin.Mode() == gin.ReleaseMode {
		r.Use(factory.Security())
		r.Use(factory.Xss())
	}
	r.Use(factory.Cors())
	r.Use(factory.RateLimiter(rate.Every(1*time.Minute), 60))

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		v1.GET("/", bookHandler.Healthcheck)
		v1.GET("/books", factory.APIKeyAuth(), bookHandler.FindBooks)
		v1.POST("/books", factory.APIKeyAuth(), factory.JWTAuth(), bookHandler.CreateBook)
		v1.GET("/books/:id", factory.APIKeyAuth(), bookHandler.FindBook)
		v1.PUT("/books/:id", factory.APIKeyAuth(), bookHandler.UpdateBook)
		v1.DELETE("/books/:id", factory.APIKeyAuth(), bookHandler.DeleteBook)

		v1.POST("/login", factory.APIKeyAuth(), userHandler.LoginHandler)
		v1.POST("/register", factory.APIKeyAuth(), userHandler.RegisterHandler)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}
