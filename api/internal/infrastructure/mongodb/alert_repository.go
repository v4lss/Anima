// Package mongodb - MongoAlertRepository implements domain/alert.Repository.
package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

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
