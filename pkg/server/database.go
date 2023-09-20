package server

import (
	"chat/pkg/api/domain"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg Config) (mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(cfg.MongoURI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}
	fmt.Println("Connected to MongoDB")

	return *client, nil
}

func ConnectPsqlDB(cfg Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DbHost, cfg.DbUser, cfg.DbName, cfg.DbPort, cfg.DbPassword)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	db.AutoMigrate(domain.PrivateChat{}, domain.ChatHistory{}, domain.PrivateChatHistory{})
	return db, nil
}
