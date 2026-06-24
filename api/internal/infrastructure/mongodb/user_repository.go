// Package mongodb - MongoUserRepository implements domain/user.Repository.
package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/v4lss/animas/internal/domain/user"
)

// UserRepository is the MongoDB adapter for user.Repository.
type UserRepository struct {
	col *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{col: db.Collection("users")}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	u.ID = primitive.NewObjectID().Hex()
	_, err := r.col.InsertOne(ctx, u)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	var u user.User
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	return &u, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	return &u, err
}
