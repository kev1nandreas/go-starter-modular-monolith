package factory

import (
	"context"
	"strconv"

	"github.com/kev1nandreas/go-rest-api-template/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDB() *mongo.Collection {
	mongoHost := config.GetEnvString("MONGO_HOST", "localhost")
	mongoPort := config.GetEnvInt("MONGO_PORT", 27017)
	mongoDatabase := config.GetEnvString("MONGO_DATABASE", "logging")
	mongoCollection := config.GetEnvString("MONGO_COLLECTION", "logs")

	uri := "mongodb://" + mongoHost + ":" + strconv.Itoa(mongoPort)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	return client.Database(mongoDatabase).Collection(mongoCollection)
}
