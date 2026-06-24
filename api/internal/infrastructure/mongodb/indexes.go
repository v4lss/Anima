// Package mongodb - database index initialization for performance.
package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureIndexes creates all necessary indexes for the collections.
func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	// Users collection
	userIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	if _, err := db.Collection("users").Indexes().CreateMany(ctx, userIndexes); err != nil {
		return err
	}

	// Monitors collection
	monitorIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "userid", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "enabled", Value: 1}},
		},
	}
	if _, err := db.Collection("monitors").Indexes().CreateMany(ctx, monitorIndexes); err != nil {
		return err
	}

	// Monitor checks collection
	checkIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "monitorid", Value: 1}, {Key: "checkedat", Value: -1}},
		},
	}
	if _, err := db.Collection("monitor_checks").Indexes().CreateMany(ctx, checkIndexes); err != nil {
		return err
	}

	// Alert configs collection
	alertConfigIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "monitorid", Value: 1}},
		},
	}
	if _, err := db.Collection("alert_configs").Indexes().CreateMany(ctx, alertConfigIndexes); err != nil {
		return err
	}

	// Alerts collection
	alertIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "monitorid", Value: 1}, {Key: "sentat", Value: -1}},
		},
	}
	if _, err := db.Collection("alerts").Indexes().CreateMany(ctx, alertIndexes); err != nil {
		return err
	}

	log.Println("MongoDB indexes created successfully")
	return nil
}
