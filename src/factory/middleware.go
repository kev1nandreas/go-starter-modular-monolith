package factory

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/kev1nandreas/go-rest-api-template/config"
	"github.com/kev1nandreas/go-rest-api-template/src/shared/auth"
	"github.com/kev1nandreas/go-rest-api-template/src/shared/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	xss "github.com/araujo88/gin-gonic-xss-middleware"
)

// APIKeyAuth validates the X-API-Key header
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == config.GetEnvString("API_SECRET_KEY", auth.GenerateRandomKey()) {
			c.Next()
		} else {
			response.NewErrorResponse(
				http.StatusUnauthorized,
				"Unauthorized",
				"Invalid API Key",
			).Send(c)
			c.Abort()
		}
	}
}

// JWTAuth validates the Bearer token in Authorization header
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		header := c.GetHeader("Authorization")
		if header == "" {
			response.NewErrorResponse(
				http.StatusUnauthorized,
				"Unauthorized",
				"Authorization header missing",
			).Send(c)
			c.Abort()
			return
		}

		if !strings.HasPrefix(header, BearerSchema) {
			response.NewErrorResponse(
				http.StatusUnauthorized,
				"Unauthorized",
				"Invalid authorization header format",
			).Send(c)
			c.Abort()
			return
		}

		tokenStr := header[len(BearerSchema):]
		claims := &auth.Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return auth.JwtKey, nil
		})

		if err != nil {
			response.NewErrorResponse(
				http.StatusUnauthorized,
				"Unauthorized",
				"Invalid token: "+err.Error(),
			).Send(c)
			c.Abort()
			return
		}

		if !token.Valid {
			response.NewErrorResponse(
				http.StatusUnauthorized,
				"Unauthorized",
				"Token is not valid",
			).Send(c)
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}

// Cors configures CORS middleware
func Cors() gin.HandlerFunc {
	appURL := config.GetEnvString("APP_URL_PROD", "http://localhost:8080")
	allowedOrigins := []string{
		"http://127.0.0.1",
		"http://localhost",
		"http://localhost:3000",
		"http://localhost:8080",
		"http://localhost:5173",
		appURL,
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// Security configures security headers
func Security() gin.HandlerFunc {
	return secure.New(secure.Config{
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		IENoOpen:              true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
	})
}

// Xss provides XSS protection
func Xss() gin.HandlerFunc {
	var xssMdlwr xss.XssMw
	return xssMdlwr.RemoveXss()
}

// RateLimiter limits request rate
func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(r, b)
	return func(c *gin.Context) {
		if !limiter.AllowN(time.Now(), 1) {
			c.String(http.StatusTooManyRequests, "Rate limit exceeded")
			c.Abort()
			return
		}
		c.Next()
	}
}

// Logger logs requests to MongoDB
func Logger(logger *zap.Logger, collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		logger.Info("Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		)

		logEntry := bson.M{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   duration,
			"ip":         c.ClientIP(),
			"user-agent": c.Request.UserAgent(),
			"errors":     c.Errors.ByType(gin.ErrorTypePrivate).String(),
		}

		if _, err := collection.InsertOne(context.TODO(), logEntry); err != nil {
			logger.Error("Failed to log to MongoDB", zap.Error(err))
		}
	}
}
