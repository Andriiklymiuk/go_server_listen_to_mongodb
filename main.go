package main

import (
	"andriiklymiuk/go_server_listen_to_mongodb/utils"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	envConfig, err := utils.LoadConnectionConfig()
	if err != nil {
		color.Red("Couldn't load env variables: \n%v", err)
		os.Exit(1)
	}
	time.Sleep(3 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoDbUrl := fmt.Sprintf("mongodb://%s:%d/%s",
		envConfig.DbHost,
		envConfig.DbPort,
		envConfig.DbName,
	)
	credential := options.Credential{
		Username: envConfig.DbUser,
		Password: envConfig.DbPassword,
	}
	client, err := mongo.Connect(
		ctx,
		options.
			Client().
			ApplyURI(mongoDbUrl).
			SetAuth(credential),
	)

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = pingDatabaseUntilConnected(ctx, client, 10, 1*time.Second)
	if err != nil {
		fmt.Println(err, "Pinging mongo db failed: ", err.Error())
		return
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		color.Red(err.Error())
	} else {
		color.Green("Connection to mongodb established")
	}
}

func pingDatabaseUntilConnected(ctx context.Context, dbClient *mongo.Client, maxRetries int, retryInterval time.Duration) error {
	for i := 0; i < maxRetries; i++ {
		err := dbClient.Ping(ctx, readpref.Primary())
		if err == nil {
			return nil
		}
		color.Red("Error pinging the database: %s. Retrying in %s...\n", err, retryInterval)
		time.Sleep(retryInterval)
	}

	return fmt.Errorf("exceeded maximum number of retries")
}
