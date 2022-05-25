package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

func Connect(config *Config) (MongoInstance, error) {
	dsn := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		config.User, config.Password, config.Host, config.Port)

	client, err := mongo.NewClient(options.Client().ApplyURI(dsn))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db := client.Database(config.DBName)

	err = client.Connect(ctx)
	if err != nil {
		return MongoInstance{}, err
	}
	mgi := MongoInstance{
		Client: client,
		Db:     db,
	}
	return mgi, nil
}
