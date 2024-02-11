package repositories

import (
	"errors"
	"go-test/models"

	"gorm.io/gorm"
)

type IItemRepository interface {
	FindAll() (*[]models.Item, error)
	FindById(itemId uint, userId uint) (*models.Item, error)
	Create(newItem models.Item) (*models.Item, error)
	Update(updateItem models.Item) (*models.Item, error)
	Delete(itemId uint, userId uint) error
}

type ItemMemoryRepository struct {
	items []models.Item
}

func NewItemMemoryRepository(items []models.Item) IItemRepository {
	return &ItemMemoryRepository{items: items}
}

func (r *ItemMemoryRepository) FindAll() (*[]models.Item, error) {
	return &r.items, nil
}

func (r *ItemMemoryRepository) FindById(itemId uint, userId uint) (*models.Item, error) {
	for _, item := range r.items {
		if item.ID == itemId {
			return &item, nil
		}
	}
	return nil, nil
}

func (r *ItemMemoryRepository) Create(newItem models.Item) (*models.Item, error) {
	newItem.ID = uint(len(r.items) + 1)
	r.items = append(r.items, newItem)
	return &newItem, nil
}

func (r *ItemMemoryRepository) Update(updateItem models.Item) (*models.Item, error) {
	for i, item := range r.items {
		if item.ID == updateItem.ID {
			r.items[i] = updateItem
			return &updateItem, nil
		}
	}
	return nil, nil
}

func (r *ItemMemoryRepository) Delete(itemId uint, userId uint) error {
	for i, v := range r.items {
		if v.ID == itemId {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return nil
		}
	}
	return errors.New("Item not found")
}

type ItemDBRepository struct {
	db *gorm.DB
}

// Create implements IItemRepository.
func (r *ItemDBRepository) Create(newItem models.Item) (*models.Item, error) {
	result := r.db.Create(&newItem)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newItem, nil
}

// FindAll implements IItemRepository.
func (r *ItemDBRepository) FindAll() (*[]models.Item, error) {
	var items []models.Item
	result := r.db.Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return &items, nil
}

// FindById implements IItemRepository.
func (r *ItemDBRepository) FindById(itemId uint, userId uint) (*models.Item, error) {
	var item models.Item
	result := r.db.First(&item, "id = ? AND user_id = ?", itemId, userId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

// Update implements IItemRepository.
func (r *ItemDBRepository) Update(updateItem models.Item) (*models.Item, error) {
	result := r.db.Save(&updateItem)
	if result.Error != nil {
		return nil, result.Error
	}
	return &updateItem, nil
}

// Delete implements IItemRepository.
func (r *ItemDBRepository) Delete(itemId uint, userId uint) error {
	result := r.db.Delete(&models.Item{}, "id = ? AND user_id = ?", itemId, userId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("Item not found")
	}
	return nil
}

func NewItemRepository(db *gorm.DB) IItemRepository {
	return &ItemDBRepository{db: db}
}
