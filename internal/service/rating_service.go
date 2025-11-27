// Package service provides gRPC service implementations for the rating service.
package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/NumberMuncherGeorge/sandbox--go-grpc-microservice/api/gen/go/rating/v1"
	"github.com/NumberMuncherGeorge/sandbox--go-grpc-microservice/internal/model"
	"github.com/NumberMuncherGeorge/sandbox--go-grpc-microservice/internal/repository"
)

const (
	defaultPageSize = 10
	maxPageSize     = 100
	minScore        = 1
	maxScore        = 5
)

// RatingService implements the gRPC RatingService
type RatingService struct {
	pb.UnimplementedRatingServiceServer
	resourceRepo *repository.ResourceRepository
	ratingRepo   *repository.RatingRepository
}

// NewRatingService creates a new RatingService
func NewRatingService(resourceRepo *repository.ResourceRepository, ratingRepo *repository.RatingRepository) *RatingService {
	return &RatingService{
		resourceRepo: resourceRepo,
		ratingRepo:   ratingRepo,
	}
}

// CreateResource creates a new resource that can be rated
func (s *RatingService) CreateResource(ctx context.Context, req *pb.CreateResourceRequest) (*pb.CreateResourceResponse, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if len(req.GetCategories()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one category is required")
	}

	resource := &model.Resource{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Categories:  req.GetCategories(),
	}

	created, err := s.resourceRepo.Create(ctx, resource)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create resource: %v", err)
	}

	return &pb.CreateResourceResponse{
		Resource: modelResourceToProto(created),
	}, nil
}

// GetResource retrieves a resource by its ID
func (s *RatingService) GetResource(ctx context.Context, req *pb.GetResourceRequest) (*pb.GetResourceResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	resource, err := s.resourceRepo.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "resource not found: %v", err)
	}

	return &pb.GetResourceResponse{
		Resource: modelResourceToProto(resource),
	}, nil
}

// ListResources lists all available resources
func (s *RatingService) ListResources(ctx context.Context, req *pb.ListResourcesRequest) (*pb.ListResourcesResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	resources, nextPageToken, err := s.resourceRepo.List(ctx, pageSize, req.GetPageToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list resources: %v", err)
	}

	protoResources := make([]*pb.Resource, len(resources))
	for i, r := range resources {
		protoResources[i] = modelResourceToProto(r)
	}

	return &pb.ListResourcesResponse{
		Resources:     protoResources,
		NextPageToken: nextPageToken,
	}, nil
}

// SubmitRating submits a rating for a resource in a specific category
func (s *RatingService) SubmitRating(ctx context.Context, req *pb.SubmitRatingRequest) (*pb.SubmitRatingResponse, error) {
	if req.GetResourceId() == "" {
		return nil, status.Error(codes.InvalidArgument, "resource_id is required")
	}
	if req.GetCategory() == "" {
		return nil, status.Error(codes.InvalidArgument, "category is required")
	}
	if req.GetScore() < minScore || req.GetScore() > maxScore {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("score must be between %d and %d", minScore, maxScore))
	}
	if req.GetRaterId() == "" {
		return nil, status.Error(codes.InvalidArgument, "rater_id is required")
	}

	// Verify the resource exists
	resource, err := s.resourceRepo.GetByID(ctx, req.GetResourceId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "resource not found: %v", err)
	}

	// Verify the category is valid for this resource
	categoryValid := false
	for _, c := range resource.Categories {
		if c == req.GetCategory() {
			categoryValid = true
			break
		}
	}
	if !categoryValid {
		return nil, status.Errorf(codes.InvalidArgument, "invalid category %q for this resource", req.GetCategory())
	}

	rating := &model.Rating{
		ResourceID: req.GetResourceId(),
		Category:   req.GetCategory(),
		Score:      req.GetScore(),
		Comment:    req.GetComment(),
		RaterID:    req.GetRaterId(),
	}

	created, err := s.ratingRepo.Create(ctx, rating)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create rating: %v", err)
	}

	return &pb.SubmitRatingResponse{
		Rating: modelRatingToProto(created),
	}, nil
}

// GetResourceRatings retrieves all ratings for a resource
func (s *RatingService) GetResourceRatings(ctx context.Context, req *pb.GetResourceRatingsRequest) (*pb.GetResourceRatingsResponse, error) {
	if req.GetResourceId() == "" {
		return nil, status.Error(codes.InvalidArgument, "resource_id is required")
	}

	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	// Verify resource exists
	_, err := s.resourceRepo.GetByID(ctx, req.GetResourceId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "resource not found: %v", err)
	}

	ratings, nextPageToken, err := s.ratingRepo.GetByResourceID(ctx, req.GetResourceId(), req.GetCategory(), pageSize, req.GetPageToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get ratings: %v", err)
	}

	summaries, err := s.ratingRepo.GetCategoryRatingSummaries(ctx, req.GetResourceId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get rating summaries: %v", err)
	}

	protoRatings := make([]*pb.Rating, len(ratings))
	for i, r := range ratings {
		protoRatings[i] = modelRatingToProto(r)
	}

	protoSummaries := make([]*pb.CategoryRatingSummary, len(summaries))
	for i, s := range summaries {
		protoSummaries[i] = &pb.CategoryRatingSummary{
			Category:     s.Category,
			AverageScore: s.AverageScore,
			TotalRatings: s.TotalRatings,
		}
	}

	return &pb.GetResourceRatingsResponse{
		Ratings:       protoRatings,
		Summaries:     protoSummaries,
		NextPageToken: nextPageToken,
	}, nil
}

func modelResourceToProto(r *model.Resource) *pb.Resource {
	return &pb.Resource{
		Id:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Categories:  r.Categories,
		CreatedAt:   timestamppb.New(r.CreatedAt),
		UpdatedAt:   timestamppb.New(r.UpdatedAt),
	}
}

func modelRatingToProto(r *model.Rating) *pb.Rating {
	return &pb.Rating{
		Id:         r.ID,
		ResourceId: r.ResourceID,
		Category:   r.Category,
		Score:      r.Score,
		Comment:    r.Comment,
		RaterId:    r.RaterID,
		CreatedAt:  timestamppb.New(r.CreatedAt),
	}
}
