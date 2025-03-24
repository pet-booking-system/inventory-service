package server

import (
	"context"
	"errors"
	"strings"

	"invservice/internal/service"

	inventorypb "github.com/azhaxyly/proto-definitions/inventory"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type InventoryServer struct {
	inventorypb.UnimplementedInventoryServiceServer
	invService service.InventoryService
}

func NewInventoryServer(invService service.InventoryService) *InventoryServer {
	return &InventoryServer{invService: invService}
}

func (s *InventoryServer) CreateResource(ctx context.Context, req *inventorypb.CreateResourceRequest) (*inventorypb.CreateResourceResponse, error) {
	resource, err := s.invService.CreateResource(req.Name, req.Type, req.Description)
	if err != nil {
		if strings.Contains(err.Error(), "required") {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to create resource")
	}

	return &inventorypb.CreateResourceResponse{
		Resource: &inventorypb.Resource{
			ResourceId:  resource.ResourceID.String(),
			Name:        resource.Name,
			Status:      resource.Status,
			Type:        resource.Type,
			Description: resource.Description,
		},
	}, nil
}

func (s *InventoryServer) ListResources(ctx context.Context, req *inventorypb.ListResourcesRequest) (*inventorypb.ListResourcesResponse, error) {
	resources, err := s.invService.ListResources()
	if err != nil {
		return nil, err
	}

	var pbResources []*inventorypb.Resource
	for _, r := range resources {
		pbResources = append(pbResources, &inventorypb.Resource{
			ResourceId: r.ResourceID.String(),
			Name:       r.Name,
			Status:     r.Status,
		})
	}

	return &inventorypb.ListResourcesResponse{
		Resources: pbResources,
	}, nil
}

func (s *InventoryServer) GetResource(ctx context.Context, req *inventorypb.GetResourceRequest) (*inventorypb.GetResourceResponse, error) {
	resource, err := s.invService.GetResource(req.ResourceId)
	if err != nil {
		if strings.Contains(err.Error(), "invalid uuid format") {
			return nil, status.Error(codes.InvalidArgument, "invalid resource id format")
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "resource not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &inventorypb.GetResourceResponse{
		Resource: &inventorypb.Resource{
			ResourceId:  resource.ResourceID.String(),
			Name:        resource.Name,
			Status:      resource.Status,
			Type:        resource.Type,
			Description: resource.Description,
		},
	}, nil
}

func (s *InventoryServer) CheckAvailability(ctx context.Context, req *inventorypb.CheckAvailabilityRequest) (*inventorypb.CheckAvailabilityResponse, error) {
	available, err := s.invService.CheckAvailability(req.ResourceId)
	if err != nil {
		if strings.Contains(err.Error(), "invalid uuid format") {
			return nil, status.Error(codes.InvalidArgument, "invalid resource id format")
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "resource not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &inventorypb.CheckAvailabilityResponse{
		IsAvailable: available,
	}, nil
}

func (s *InventoryServer) UpdateResourceStatus(ctx context.Context, req *inventorypb.UpdateResourceStatusRequest) (*inventorypb.UpdateResourceStatusResponse, error) {
	updatedResource, err := s.invService.UpdateResourceStatus(req.ResourceId, req.NewStatus)
	if err != nil {
		if strings.Contains(err.Error(), "invalid uuid format") {
			return nil, status.Error(codes.InvalidArgument, "invalid resource id format")
		}
		if strings.Contains(err.Error(), "invalid status") {
			return nil, status.Error(codes.InvalidArgument, "invalid status value")
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "resource not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &inventorypb.UpdateResourceStatusResponse{
		ResourceId: updatedResource.ResourceID.String(),
		Status:     updatedResource.Status,
	}, nil
}

func (s *InventoryServer) DeleteResource(ctx context.Context, req *inventorypb.DeleteResourceRequest) (*inventorypb.DeleteResourceResponse, error) {
	err := s.invService.DeleteResource(req.ResourceId)
	if err != nil {
		if strings.Contains(err.Error(), "invalid uuid format") {
			return nil, status.Error(codes.InvalidArgument, "invalid resource id format")
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "resource not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete resource")
	}

	return &inventorypb.DeleteResourceResponse{
		ResourceId: req.ResourceId,
		Deleted:    true,
	}, nil
}
