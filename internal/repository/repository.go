// Package repository provides MongoDB repository implementations for the rating service.
package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/NumberMuncherGeorge/sandbox--go-grpc-microservice/internal/model"
)

// ResourceRepository handles resource operations in MongoDB
type ResourceRepository struct {
	collection *mongo.Collection
}

// NewResourceRepository creates a new ResourceRepository
func NewResourceRepository(db *mongo.Database) *ResourceRepository {
	return &ResourceRepository{
		collection: db.Collection("resources"),
	}
}

// Create creates a new resource
func (r *ResourceRepository) Create(ctx context.Context, resource *model.Resource) (*model.Resource, error) {
	resource.ID = primitive.NewObjectID().Hex()
	resource.CreatedAt = time.Now()
	resource.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// GetByID retrieves a resource by ID
func (r *ResourceRepository) GetByID(ctx context.Context, id string) (*model.Resource, error) {
	var resource model.Resource
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// List retrieves all resources with pagination
func (r *ResourceRepository) List(ctx context.Context, pageSize int32, pageToken string) ([]*model.Resource, string, error) {
	filter := bson.M{}
	if pageToken != "" {
		filter["_id"] = bson.M{"$gt": pageToken}
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "_id", Value: 1}}).
		SetLimit(int64(pageSize + 1)) // Fetch one extra to determine if there are more results

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, "", err
	}
	defer cursor.Close(ctx)

	var resources []*model.Resource
	if err = cursor.All(ctx, &resources); err != nil {
		return nil, "", err
	}

	var nextPageToken string
	if len(resources) > int(pageSize) {
		// There are more results, use the last item's ID from the current page as the token
		resources = resources[:pageSize]
		nextPageToken = resources[pageSize-1].ID
	}

	return resources, nextPageToken, nil
}

// RatingRepository handles rating operations in MongoDB
type RatingRepository struct {
	collection *mongo.Collection
}

// NewRatingRepository creates a new RatingRepository
func NewRatingRepository(db *mongo.Database) *RatingRepository {
	return &RatingRepository{
		collection: db.Collection("ratings"),
	}
}

// Create creates a new rating
func (r *RatingRepository) Create(ctx context.Context, rating *model.Rating) (*model.Rating, error) {
	rating.ID = primitive.NewObjectID().Hex()
	rating.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, rating)
	if err != nil {
		return nil, err
	}

	return rating, nil
}

// GetByResourceID retrieves all ratings for a resource with optional category filter
func (r *RatingRepository) GetByResourceID(ctx context.Context, resourceID string, category string, pageSize int32, pageToken string) ([]*model.Rating, string, error) {
	filter := bson.M{"resource_id": resourceID}
	if category != "" {
		filter["category"] = category
	}
	if pageToken != "" {
		filter["_id"] = bson.M{"$gt": pageToken}
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "_id", Value: 1}}).
		SetLimit(int64(pageSize + 1))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, "", err
	}
	defer cursor.Close(ctx)

	var ratings []*model.Rating
	if err = cursor.All(ctx, &ratings); err != nil {
		return nil, "", err
	}

	var nextPageToken string
	if len(ratings) > int(pageSize) {
		// There are more results, use the last item's ID from the current page as the token
		ratings = ratings[:pageSize]
		nextPageToken = ratings[pageSize-1].ID
	}

	return ratings, nextPageToken, nil
}

// GetCategoryRatingSummaries retrieves aggregated rating summaries for a resource
func (r *RatingRepository) GetCategoryRatingSummaries(ctx context.Context, resourceID string) ([]*model.CategoryRatingSummary, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"resource_id": resourceID}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$category"},
			{Key: "average_score", Value: bson.M{"$avg": "$score"}},
			{Key: "total_ratings", Value: bson.M{"$sum": 1}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var summaries []*model.CategoryRatingSummary
	if err = cursor.All(ctx, &summaries); err != nil {
		return nil, err
	}

	return summaries, nil
}
