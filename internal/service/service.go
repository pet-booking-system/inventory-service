package service

import (
	"errors"
	"fmt"
	"invservice/internal/logger"
	"invservice/internal/models"
	"invservice/internal/repository"

	"gorm.io/gorm"
)

type InventoryService interface {
	CreateResource(name, resType, description string) (*models.Resource, error)
	ListResources() ([]models.Resource, error)
	GetResource(resourceID string) (*models.Resource, error)
	CheckAvailability(resourceID string) (bool, error)
	UpdateResourceStatus(resourceID string, newStatus string) (*models.Resource, error)
	DeleteResource(resourceID string) error
}

type inventoryService struct {
	repo repository.InventoryRepository
}

func NewInventoryService(repo repository.InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) CreateResource(name, resType, description string) (*models.Resource, error) {
	if name == "" || resType == "" {
		return nil, fmt.Errorf("name and type are required")
	}
	return s.repo.CreateResource(name, resType, description)
}

func (s *inventoryService) ListResources() ([]models.Resource, error) {
	return s.repo.ListResources()
}

func (s *inventoryService) GetResource(resourceID string) (*models.Resource, error) {
	return s.repo.GetResourceByID(resourceID)
}

func (s *inventoryService) CheckAvailability(resourceID string) (bool, error) {
	resource, err := s.repo.GetResourceByID(resourceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
		logger.Error("Failed to check availability: %w", err)
		return false, err
	}
	if resource.Status == "available" {
		return true, nil
	}
	return false, nil
}

func (s *inventoryService) UpdateResourceStatus(resourceID string, newStatus string) (*models.Resource, error) {
	allowedStatuses := map[string]bool{
		"available":   true,
		"booked":      true,
		"unavailable": true,
	}

	if !allowedStatuses[newStatus] {
		logger.Error("Invalid status: %w", newStatus)
		return nil, fmt.Errorf("invalid status: %s", newStatus)
	}

	updatedResource, err := s.repo.UpdateResourceStatus(resourceID, newStatus)
	if err != nil {
		return nil, err
	}

	return updatedResource, nil
}

func (s *inventoryService) DeleteResource(resourceID string) error {
	return s.repo.DeleteResource(resourceID)
}
