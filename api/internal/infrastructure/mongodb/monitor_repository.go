// Package mongodb - MongoMonitorRepository implements domain/monitor.Repository.
package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/v4lss/animas/internal/domain/monitor"
)

// MonitorRepository is the MongoDB adapter for monitor.Repository.
type MonitorRepository struct {
	col *mongo.Collection
}

func NewMonitorRepository(db *mongo.Database) *MonitorRepository {
	return &MonitorRepository{col: db.Collection("monitors")}
}

func (r *MonitorRepository) Create(ctx context.Context, m *monitor.Monitor) error {
	m.ID = primitive.NewObjectID().Hex()
	_, err := r.col.InsertOne(ctx, m)
	return err
}

func (r *MonitorRepository) FindByID(ctx context.Context, id string) (*monitor.Monitor, error) {
	var m monitor.Monitor
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&m)
	return &m, err
}

func (r *MonitorRepository) FindByUserID(ctx context.Context, userID string) ([]*monitor.Monitor, error) {
	cursor, err := r.col.Find(ctx, bson.M{"userid": userID})
	if err != nil {
		return nil, err
	}
	var monitors []*monitor.Monitor
	return monitors, cursor.All(ctx, &monitors)
}

func (r *MonitorRepository) Update(ctx context.Context, m *monitor.Monitor) error {
	_, err := r.col.ReplaceOne(ctx, bson.M{"_id": m.ID}, m)
	return err
}

func (r *MonitorRepository) Delete(ctx context.Context, id string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
