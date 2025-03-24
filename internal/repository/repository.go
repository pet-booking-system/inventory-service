package repository

import (
	"fmt"
	"invservice/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	CreateResource(name, resType, description string) (*models.Resource, error)
	ListResources() ([]models.Resource, error)
	GetResourceByID(id string) (*models.Resource, error)
	UpdateResourceStatus(resourceID string, newStatus string) (*models.Resource, error)
	DeleteResource(resourceID string) error
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) CreateResource(name, resType, description string) (*models.Resource, error) {
	newResource := &models.Resource{
		Name:        name,
		Type:        resType,
		Status:      "available", // default status
		Description: description,
	}

	if err := r.db.Create(newResource).Error; err != nil {
		return nil, err
	}
	return newResource, nil
}

func (r *inventoryRepository) ListResources() ([]models.Resource, error) {
	var resources []models.Resource
	if err := r.db.Find(&resources).Error; err != nil {
		return nil, err
	}
	return resources, nil
}

func (r *inventoryRepository) GetResourceByID(id string) (*models.Resource, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid format: %w", err)
	}
	var resource models.Resource
	if err := r.db.First(&resource, "resource_id = ?", uid).Error; err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r *inventoryRepository) UpdateResourceStatus(resourceID string, newStatus string) (*models.Resource, error) {
	uid, err := uuid.Parse(resourceID)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid format: %w", err)
	}

	var resource models.Resource
	if err := r.db.First(&resource, "resource_id = ?", uid).Error; err != nil {
		return nil, err
	}

	resource.Status = newStatus
	if err := r.db.Save(&resource).Error; err != nil {
		return nil, err
	}

	return &resource, nil
}

func (r *inventoryRepository) DeleteResource(resourceID string) error {
	uid, err := uuid.Parse(resourceID)
	if err != nil {
		return fmt.Errorf("invalid uuid format: %w", err)
	}

	result := r.db.Delete(&models.Resource{}, "resource_id = ?", uid)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
