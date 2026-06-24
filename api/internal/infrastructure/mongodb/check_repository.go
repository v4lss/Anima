// Package mongodb - MongoCheckRepository implements domain/monitor.CheckRepository.
package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/v4lss/animas/internal/domain/monitor"
)

// CheckRepository is the MongoDB adapter for monitor.CheckRepository.
type CheckRepository struct {
	col *mongo.Collection
}

func NewCheckRepository(db *mongo.Database) *CheckRepository {
	return &CheckRepository{col: db.Collection("monitor_checks")}
}

func (r *CheckRepository) Save(ctx context.Context, c *monitor.Check) error {
	c.ID = primitive.NewObjectID().Hex()
	_, err := r.col.InsertOne(ctx, c)
	return err
}

func (r *CheckRepository) FindByMonitorID(ctx context.Context, monitorID string, limit int) ([]*monitor.Check, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "checkedat", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.col.Find(ctx, bson.M{"monitorid": monitorID}, opts)
	if err != nil {
		return nil, err
	}
	var checks []*monitor.Check
	return checks, cursor.All(ctx, &checks)
}
