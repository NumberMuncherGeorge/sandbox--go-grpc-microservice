// Package model contains data models for the rating service.
package model

import (
	"time"
)

// Resource represents an entity that can be rated
type Resource struct {
	ID          string    `bson:"_id,omitempty"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	Categories  []string  `bson:"categories"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

// Rating represents a rating submission for a resource
type Rating struct {
	ID         string    `bson:"_id,omitempty"`
	ResourceID string    `bson:"resource_id"`
	Category   string    `bson:"category"`
	Score      int32     `bson:"score"` // 1-5
	Comment    string    `bson:"comment"`
	RaterID    string    `bson:"rater_id"`
	CreatedAt  time.Time `bson:"created_at"`
}

// CategoryRatingSummary provides aggregated rating info for a category
type CategoryRatingSummary struct {
	Category     string  `bson:"_id"`
	AverageScore float64 `bson:"average_score"`
	TotalRatings int32   `bson:"total_ratings"`
}
