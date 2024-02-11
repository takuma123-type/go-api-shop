package services

import (
	"go-test/dto"
	"go-test/models"
	"go-test/repositories"
)

type IItemService interface {
	FindAll() (*[]models.Item, error)
	FindById(itemId uint, userId uint) (*models.Item, error)
	Create(CreateItemInput dto.CreateItemInput, userId uint) (*models.Item, error)
	Update(itemId uint, userId uint, updateItemInput dto.UpdateItemInput) (*models.Item, error)
	Delete(itemId uint, userId uint) error
}

type ItemService struct {
	repository repositories.IItemRepository
}

func NewItemService(repository repositories.IItemRepository) IItemService {
	return &ItemService{repository: repository}
}

func (s *ItemService) FindAll() (*[]models.Item, error) {
	return s.repository.FindAll()
}

func (s *ItemService) FindById(itemId uint, userId uint) (*models.Item, error) {
	return s.repository.FindById(itemId, userId)
}

func (s *ItemService) Create(CreateItemInput dto.CreateItemInput, userId uint) (*models.Item, error) {
	newItem := models.Item{
		Name:        CreateItemInput.Name,
		Price:       CreateItemInput.Price,
		Description: CreateItemInput.Description,
		SoldOut:     false,
		UserID:      userId,
	}
	return s.repository.Create(newItem)
}

func (s *ItemService) Update(itemId uint, userId uint, updateItemInput dto.UpdateItemInput) (*models.Item, error) {
	targetItem, err := s.FindById(itemId, userId)
	if err != nil {
		return nil, err
	}

	if updateItemInput.Name != "" {
		targetItem.Name = updateItemInput.Name
	}
	if updateItemInput.Price != 0 {
		targetItem.Price = updateItemInput.Price
	}
	if updateItemInput.Description != "" {
		targetItem.Description = updateItemInput.Description
	}
	if updateItemInput.SoldOut != false {
		targetItem.SoldOut = updateItemInput.SoldOut
	}
	return s.repository.Update(*targetItem)
}

func (s *ItemService) Delete(itemId uint, userId uint) error {
	return s.repository.Delete(itemId, userId)
}
