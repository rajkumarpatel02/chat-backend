package config

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var UserCollection *mongo.Collection // ✅ Add this line

func ConnectMongo() error {
    username := "Rajkumarpatel"
    password := "35t8DElpuPB94kvc" // If this contains special characters, encode it!
    cluster := "cluster0.bkfnv.mongodb.net"
    dbName := "chaty" // use a valid DB name (no spaces like "Project 0")

    // ✅ Correct way to build URI using fmt.Sprintf and variables
    uri := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority&appName=Cluster0", username, password, cluster)


	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOpts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return fmt.Errorf("Mongo connect error: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("Mongo ping error: %v", err)
	}

	fmt.Println("✅ Connected to MongoDB Atlas")

	DB = client.Database(dbName)
	UserCollection = DB.Collection("users") // ✅ Add this line

	return nil
}
