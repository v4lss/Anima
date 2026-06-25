// Package mongodb - MongoCheckRepository implements domain/monitor.CheckRepository.
package mongodb

import (
	"context"
	"time"

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
	err = cursor.All(ctx, &checks)
	return checks, err
}

func (r *CheckRepository) FindByMonitorIDPaginated(ctx context.Context, monitorID string, page, limit int, status string, fromDate, toDate string) ([]*monitor.Check, int64, error) {
	filter := bson.M{"monitorid": monitorID}

	// Add status filter if provided
	if status != "" {
		filter["status"] = status
	}

	// Add date range filters if provided
	if fromDate != "" || toDate != "" {
		dateFilter := bson.M{}
		if fromDate != "" {
			if t, err := time.Parse(time.RFC3339, fromDate); err == nil {
				dateFilter["$gte"] = t
			}
		}
		if toDate != "" {
			if t, err := time.Parse(time.RFC3339, toDate); err == nil {
				dateFilter["$lte"] = t
			}
		}
		if len(dateFilter) > 0 {
			filter["checkedat"] = dateFilter
		}
	}

	// Get total count for pagination
	total, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Calculate skip for pagination
	skip := int64((page - 1) * limit)
	opts := options.Find().
		SetSort(bson.D{{Key: "checkedat", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(limit))

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	var checks []*monitor.Check
	err = cursor.All(ctx, &checks)
	return checks, total, err
}
