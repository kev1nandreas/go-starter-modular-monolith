package main

import (
	"context"
	"log"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/kev1nandreas/go-rest-api-template/config"
	"github.com/kev1nandreas/go-rest-api-template/src/factory"
	"github.com/kev1nandreas/go-rest-api-template/src/handler"
	bookFactory "github.com/kev1nandreas/go-rest-api-template/src/modules/book/factory"
	userFactory "github.com/kev1nandreas/go-rest-api-template/src/modules/user/factory"
	"go.mongodb.org/mongo-driver/mongo"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8001
// @BasePath  /api/v1

// @securityDefinitions.apikey JwtAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	appPort := config.GetEnvInt("APP_PORT", 8080)
	isLogging := config.GetEnvBool("APP_MONGO_LOGGING", false)
	isDebugging := config.GetEnvBool("APP_DEBUG", false)

	if isDebugging {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize infrastructure
	redisClient := factory.NewRedisClient()
	db := factory.NewDatabase()

	var mongoCollection *mongo.Collection
	if isLogging {
		mongoCollection = factory.NewMongoDB()
	}

	ctx := context.Background()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Initialize module usecases via factories
	bookUsecase := bookFactory.NewBookUsecase(db)
	userUsecase := userFactory.NewUserUsecase(db)

	// Initialize handlers
	bookHandler := handler.NewBookHandler(bookUsecase, redisClient, ctx)
	userHandler := handler.NewUserHandler(userUsecase)

	// Setup router
	r := handler.NewRouter(logger, mongoCollection, bookHandler, userHandler)

	if err := r.Run(":" + strconv.Itoa(appPort)); err != nil {
		log.Fatal(err)
	}
}
