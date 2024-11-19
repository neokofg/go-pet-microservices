package service

import (
	"context"
	"github.com/neokofg/go-pet-microservices/catalog-service/api/proto"
	"github.com/neokofg/go-pet-microservices/catalog-service/ent"
	"github.com/neokofg/go-pet-microservices/catalog-service/ent/item"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type CatalogService struct {
	proto.UnimplementedCatalogServiceServer
	client *ent.Client
	logger *zap.Logger
}

func NewCatalogService(client *ent.Client, logger *zap.Logger) *CatalogService {
	return &CatalogService{
		client: client,
		logger: logger,
	}
}

func (s *CatalogService) GetItems(ctx context.Context, req *proto.GetItemsRequest) (*proto.GetItemsResponse, error) {
	query := s.client.Item.Query()

	if len(req.Tags) > 0 {
		query = query.Where(item.TagsNotNil())
	}

	total, err := query.Count(ctx)
	if err != nil {
		s.logger.Error("Failed to count items", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to count items")
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	offset := int(req.Page-1) * limit

	items, err := query.
		Limit(limit).
		Offset(offset).
		Order(ent.Desc(item.FieldCreatedAt)). // По умолчанию сортируем по дате создания
		All(ctx)

	if err != nil {
		s.logger.Error("Failed to fetch items", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to fetch items")
	}

	// Преобразуем элементы в proto формат
	protoItems := make([]*proto.Item, len(items))
	for i, itm := range items {
		protoItems[i] = &proto.Item{
			Id:          itm.ID,
			Title:       itm.Title,
			Description: itm.Description,
			Tags:        itm.Tags,
			ImageUrl:    itm.ImageURL,
			Rating:      itm.Rating,
			ReviewCount: int32(itm.ReviewCount),
			CreatedAt:   itm.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   itm.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &proto.GetItemsResponse{
		Items:      protoItems,
		Total:      int32(total),
		Page:       req.Page,
		TotalPages: int32((total + limit - 1) / limit),
	}, nil
}

func (s *CatalogService) GetItem(ctx context.Context, req *proto.GetItemRequest) (*proto.Item, error) {
	itm, err := s.client.Item.Get(ctx, req.Id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "itm not found")
		}
		s.logger.Error("Failed to get itm", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get itm")
	}

	return &proto.Item{
		Id:          itm.ID,
		Title:       itm.Title,
		Description: itm.Description,
		Tags:        itm.Tags,
		ImageUrl:    itm.ImageURL,
		Rating:      itm.Rating,
		ReviewCount: int32(itm.ReviewCount),
		CreatedAt:   itm.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   itm.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *CatalogService) CreateItem(ctx context.Context, req *proto.CreateItemRequest) (*proto.Item, error) {
	itm, err := s.client.Item.Create().
		SetTitle(req.Title).
		SetDescription(req.Description).
		SetTags(req.Tags).
		SetImageURL(req.ImageUrl).
		Save(ctx)

	if err != nil {
		s.logger.Error("Failed to create itm", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create itm")
	}

	return &proto.Item{
		Id:          itm.ID,
		Title:       itm.Title,
		Description: itm.Description,
		Tags:        itm.Tags,
		ImageUrl:    itm.ImageURL,
		Rating:      itm.Rating,
		ReviewCount: int32(itm.ReviewCount),
		CreatedAt:   itm.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   itm.UpdatedAt.Format(time.RFC3339),
	}, nil
}
