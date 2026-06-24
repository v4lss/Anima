// Package mongodb - MongoAlertRepository implements domain/alert.Repository.
package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/v4lss/animas/internal/domain/alert"
)

// AlertRepository is the MongoDB adapter for alert.Repository.
type AlertRepository struct {
	col *mongo.Collection
}

func NewAlertRepository(db *mongo.Database) *AlertRepository {
	return &AlertRepository{col: db.Collection("alerts")}
}

func (r *AlertRepository) Save(ctx context.Context, a *alert.Alert) error {
	a.ID = primitive.NewObjectID().Hex()
	_, err := r.col.InsertOne(ctx, a)
	return err
}

func (r *AlertRepository) FindByMonitorID(ctx context.Context, monitorID string) ([]*alert.Alert, error) {
	cursor, err := r.col.Find(ctx, bson.M{"monitorid": monitorID})
	if err != nil {
		return nil, err
	}
	var alerts []*alert.Alert
	return alerts, cursor.All(ctx, &alerts)
}

// AlertConfigRepository is the MongoDB adapter for alert.ConfigRepository.
type AlertConfigRepository struct {
	col *mongo.Collection
}

func NewAlertConfigRepository(db *mongo.Database) *AlertConfigRepository {
	return &AlertConfigRepository{col: db.Collection("alert_configs")}
}

func (r *AlertConfigRepository) Save(ctx context.Context, c *alert.Config) error {
	now := time.Now()
	if c.ID == "" {
		c.ID = primitive.NewObjectID().Hex()
		c.CreatedAt = now
	}
	c.UpdatedAt = now

	filter := bson.M{"id": c.ID}
	update := bson.M{"$set": c}
	opts := options.Update().SetUpsert(true)

	_, err := r.col.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *AlertConfigRepository) FindByMonitorID(ctx context.Context, monitorID string) ([]*alert.Config, error) {
	cursor, err := r.col.Find(ctx, bson.M{"monitorid": monitorID})
	if err != nil {
		return nil, err
	}
	var configs []*alert.Config
	return configs, cursor.All(ctx, &configs)
}

func (r *AlertConfigRepository) Delete(ctx context.Context, id string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"id": id})
	return err
}
