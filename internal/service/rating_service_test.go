package service

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	pb "github.com/NumberMuncherGeorge/sandbox--go-grpc-microservice/api/gen/go/rating/v1"
	"github.com/NumberMuncherGeorge/sandbox--go-grpc-microservice/internal/repository"
)

func setupTestService(t *testing.T) (*RatingService, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to test MongoDB (requires MongoDB to be running)
	mongoURI := "mongodb://localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Skipf("Skipping test: MongoDB not available: %v", err)
		return nil, nil
	}

	if err := client.Ping(ctx, nil); err != nil {
		t.Skipf("Skipping test: MongoDB not available: %v", err)
		return nil, nil
	}

	// Use a test database
	db := client.Database("rating_service_test")

	resourceRepo := repository.NewResourceRepository(db)
	ratingRepo := repository.NewRatingRepository(db)
	svc := NewRatingService(resourceRepo, ratingRepo)

	cleanup := func() {
		// Clean up test database
		_ = db.Drop(context.Background())
		_ = client.Disconnect(context.Background())
	}

	return svc, cleanup
}

func TestCreateResource(t *testing.T) {
	svc, cleanup := setupTestService(t)
	if svc == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// Test successful resource creation
	req := &pb.CreateResourceRequest{
		Name:        "Test Product",
		Description: "A test product for testing",
		Categories:  []string{"quality", "value", "design"},
	}

	resp, err := svc.CreateResource(ctx, req)
	if err != nil {
		t.Fatalf("CreateResource failed: %v", err)
	}

	if resp.Resource.Name != req.Name {
		t.Errorf("Expected name %q, got %q", req.Name, resp.Resource.Name)
	}
	if resp.Resource.Id == "" {
		t.Error("Expected resource ID to be set")
	}
}

func TestCreateResourceValidation(t *testing.T) {
	svc, cleanup := setupTestService(t)
	if svc == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// Test missing name
	req := &pb.CreateResourceRequest{
		Description: "A test product",
		Categories:  []string{"quality"},
	}

	_, err := svc.CreateResource(ctx, req)
	if err == nil {
		t.Error("Expected error for missing name")
	}

	// Test missing categories
	req = &pb.CreateResourceRequest{
		Name:        "Test Product",
		Description: "A test product",
	}

	_, err = svc.CreateResource(ctx, req)
	if err == nil {
		t.Error("Expected error for missing categories")
	}
}

func TestGetResource(t *testing.T) {
	svc, cleanup := setupTestService(t)
	if svc == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// Create a resource first
	createReq := &pb.CreateResourceRequest{
		Name:        "Test Product",
		Description: "A test product for testing",
		Categories:  []string{"quality", "value"},
	}

	createResp, err := svc.CreateResource(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateResource failed: %v", err)
	}

	// Get the resource
	getReq := &pb.GetResourceRequest{
		Id: createResp.Resource.Id,
	}

	getResp, err := svc.GetResource(ctx, getReq)
	if err != nil {
		t.Fatalf("GetResource failed: %v", err)
	}

	if getResp.Resource.Name != createReq.Name {
		t.Errorf("Expected name %q, got %q", createReq.Name, getResp.Resource.Name)
	}
}

func TestListResources(t *testing.T) {
	svc, cleanup := setupTestService(t)
	if svc == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// Create multiple resources
	for i := 0; i < 3; i++ {
		req := &pb.CreateResourceRequest{
			Name:        "Test Product",
			Description: "A test product",
			Categories:  []string{"quality"},
		}
		_, err := svc.CreateResource(ctx, req)
		if err != nil {
			t.Fatalf("CreateResource failed: %v", err)
		}
	}

	// List resources
	listReq := &pb.ListResourcesRequest{
		PageSize: 10,
	}

	listResp, err := svc.ListResources(ctx, listReq)
	if err != nil {
		t.Fatalf("ListResources failed: %v", err)
	}

	if len(listResp.Resources) != 3 {
		t.Errorf("Expected 3 resources, got %d", len(listResp.Resources))
	}
}

func TestSubmitRating(t *testing.T) {
	svc, cleanup := setupTestService(t)
	if svc == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// Create a resource first
	createReq := &pb.CreateResourceRequest{
		Name:        "Test Product",
		Description: "A test product for testing",
		Categories:  []string{"quality", "value"},
	}

	createResp, err := svc.CreateResource(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateResource failed: %v", err)
	}

	// Submit a rating
	ratingReq := &pb.SubmitRatingRequest{
		ResourceId: createResp.Resource.Id,
		Category:   "quality",
		Score:      5,
		Comment:    "Excellent quality!",
		RaterId:    "user123",
	}

	ratingResp, err := svc.SubmitRating(ctx, ratingReq)
	if err != nil {
		t.Fatalf("SubmitRating failed: %v", err)
	}

	if ratingResp.Rating.Score != 5 {
		t.Errorf("Expected score 5, got %d", ratingResp.Rating.Score)
	}
	if ratingResp.Rating.Category != "quality" {
		t.Errorf("Expected category 'quality', got %q", ratingResp.Rating.Category)
	}
}

func TestSubmitRatingValidation(t *testing.T) {
	svc, cleanup := setupTestService(t)
	if svc == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// Create a resource first
	createReq := &pb.CreateResourceRequest{
		Name:        "Test Product",
		Description: "A test product for testing",
		Categories:  []string{"quality"},
	}

	createResp, err := svc.CreateResource(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateResource failed: %v", err)
	}

	// Test invalid score (too low)
	ratingReq := &pb.SubmitRatingRequest{
		ResourceId: createResp.Resource.Id,
		Category:   "quality",
		Score:      0,
		RaterId:    "user123",
	}

	_, err = svc.SubmitRating(ctx, ratingReq)
	if err == nil {
		t.Error("Expected error for invalid score")
	}

	// Test invalid score (too high)
	ratingReq.Score = 6
	_, err = svc.SubmitRating(ctx, ratingReq)
	if err == nil {
		t.Error("Expected error for invalid score")
	}

	// Test invalid category
	ratingReq.Score = 5
	ratingReq.Category = "nonexistent"
	_, err = svc.SubmitRating(ctx, ratingReq)
	if err == nil {
		t.Error("Expected error for invalid category")
	}
}

func TestGetResourceRatings(t *testing.T) {
	svc, cleanup := setupTestService(t)
	if svc == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// Create a resource
	createReq := &pb.CreateResourceRequest{
		Name:        "Test Product",
		Description: "A test product for testing",
		Categories:  []string{"quality", "value"},
	}

	createResp, err := svc.CreateResource(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateResource failed: %v", err)
	}

	// Submit multiple ratings
	ratings := []struct {
		category string
		score    int32
	}{
		{"quality", 5},
		{"quality", 4},
		{"value", 3},
	}

	for _, r := range ratings {
		ratingReq := &pb.SubmitRatingRequest{
			ResourceId: createResp.Resource.Id,
			Category:   r.category,
			Score:      r.score,
			RaterId:    "user123",
		}
		_, err := svc.SubmitRating(ctx, ratingReq)
		if err != nil {
			t.Fatalf("SubmitRating failed: %v", err)
		}
	}

	// Get ratings
	getRatingsReq := &pb.GetResourceRatingsRequest{
		ResourceId: createResp.Resource.Id,
		PageSize:   10,
	}

	getRatingsResp, err := svc.GetResourceRatings(ctx, getRatingsReq)
	if err != nil {
		t.Fatalf("GetResourceRatings failed: %v", err)
	}

	if len(getRatingsResp.Ratings) != 3 {
		t.Errorf("Expected 3 ratings, got %d", len(getRatingsResp.Ratings))
	}

	if len(getRatingsResp.Summaries) != 2 {
		t.Errorf("Expected 2 category summaries, got %d", len(getRatingsResp.Summaries))
	}

	// Check summaries
	for _, summary := range getRatingsResp.Summaries {
		switch summary.Category {
		case "quality":
			if summary.TotalRatings != 2 {
				t.Errorf("Expected 2 quality ratings, got %d", summary.TotalRatings)
			}
			expectedAvg := 4.5
			if summary.AverageScore != expectedAvg {
				t.Errorf("Expected quality average %.1f, got %.1f", expectedAvg, summary.AverageScore)
			}
		case "value":
			if summary.TotalRatings != 1 {
				t.Errorf("Expected 1 value rating, got %d", summary.TotalRatings)
			}
		}
	}
}
