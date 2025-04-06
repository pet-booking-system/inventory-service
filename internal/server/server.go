package server

import (
	"context"
	"errors"
	"strings"

	"invservice/internal/logger"
	"invservice/internal/service"

	inventorypb "github.com/pet-booking-system/proto-definitions/inventory"
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
	logger.Info("Received request to create resource: ", req)
	resource, err := s.invService.CreateResource(req.Name, req.Type, req.Description)
	if err != nil {
		if strings.Contains(err.Error(), "required") {
			logger.Error("Invalid argument while creating resource: ", err)
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		logger.Error("Failed to create resource: ", err)
		return nil, status.Error(codes.Internal, "failed to create resource")
	}

	logger.Info("Resource created successfully: ", resource)
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
	logger.Info("Received request to list resources: ", req)
	resources, err := s.invService.ListResources()
	if err != nil {
		logger.Error("Failed to list resources: ", err)
		return nil, err
	}

	var pbResources []*inventorypb.Resource
	for _, r := range resources {
		pbResources = append(pbResources, &inventorypb.Resource{
			ResourceId:  r.ResourceID.String(),
			Name:        r.Name,
			Status:      r.Status,
			Type:        r.Type,
			Description: r.Description,
		})
	}

	return &inventorypb.ListResourcesResponse{
		Resources: pbResources,
	}, nil
}

func (s *InventoryServer) GetResource(ctx context.Context, req *inventorypb.GetResourceRequest) (*inventorypb.GetResourceResponse, error) {
	logger.Info("Received request to get resource: ", req)
	resource, err := s.invService.GetResource(req.ResourceId)
	if err != nil {
		if strings.Contains(err.Error(), "invalid uuid format") {
			logger.Error("Invalid UUID format for resource ID: ", req.ResourceId)
			return nil, status.Error(codes.InvalidArgument, "invalid resource id format")
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("Resource not found with ID: ", req.ResourceId)
			return nil, status.Error(codes.NotFound, "resource not found")
		}
		logger.Error("Failed to get resource with ID ", req.ResourceId, ": ", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}
	logger.Info("Resource retrieved successfully: ", resource)

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
			logger.Error("Invalid UUID format for resource ID: ", req.ResourceId)
			return nil, status.Error(codes.InvalidArgument, "invalid resource id format")
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("Resource not found with ID: ", req.ResourceId)
			return nil, status.Error(codes.NotFound, "resource not found")
		}
		logger.Error("Failed to check availability for resource ID ", req.ResourceId, ": ", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &inventorypb.CheckAvailabilityResponse{
		IsAvailable: available,
	}, nil
}

func (s *InventoryServer) UpdateResourceStatus(ctx context.Context, req *inventorypb.UpdateResourceStatusRequest) (*inventorypb.UpdateResourceStatusResponse, error) {
	logger.Info("Received request to update resource status: ", req)
	updatedResource, err := s.invService.UpdateResourceStatus(req.ResourceId, req.NewStatus)
	if err != nil {
		if strings.Contains(err.Error(), "invalid uuid format") {
			logger.Error("Invalid UUID format in UpdateResourceStatus for resource ID: ", req.ResourceId)
			return nil, status.Error(codes.InvalidArgument, "invalid resource id format")
		}
		if strings.Contains(err.Error(), "invalid status") {
			logger.Error("Invalid status value in UpdateResourceStatus: ", req.NewStatus)
			return nil, status.Error(codes.InvalidArgument, "invalid status value")
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("Resource not found in UpdateResourceStatus with ID: ", req.ResourceId)
			return nil, status.Error(codes.NotFound, "resource not found")
		}
		logger.Error("Failed to update status for resource ID ", req.ResourceId, ": ", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}
	logger.Info("Resource status updated successfully: ", updatedResource)

	return &inventorypb.UpdateResourceStatusResponse{
		ResourceId: updatedResource.ResourceID.String(),
		Status:     updatedResource.Status,
	}, nil
}

func (s *InventoryServer) DeleteResource(ctx context.Context, req *inventorypb.DeleteResourceRequest) (*inventorypb.DeleteResourceResponse, error) {
	logger.Info("Received request to delete resource: ", req)
	err := s.invService.DeleteResource(req.ResourceId)
	if err != nil {
		if strings.Contains(err.Error(), "invalid uuid format") {
			logger.Error("Invalid UUID format in DeleteResource for ID: ", req.ResourceId)
			return nil, status.Error(codes.InvalidArgument, "invalid resource id format")
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("Resource not found in DeleteResource with ID: ", req.ResourceId)
			return nil, status.Error(codes.NotFound, "resource not found")
		}
		logger.Error("Failed to delete resource with ID ", req.ResourceId, ": ", err)
		return nil, status.Error(codes.Internal, "failed to delete resource")
	}
	logger.Info("Resource deleted successfully: ", req.ResourceId)

	return &inventorypb.DeleteResourceResponse{
		ResourceId: req.ResourceId,
		Deleted:    true,
	}, nil
}
