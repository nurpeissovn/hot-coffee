package service

import (
	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type inventoryService struct {
	in dal.InventDal
}

func ReturnInvent(i dal.InventDal) *inventoryService {
	return &inventoryService{in: i}
}

type InventServiceInter interface {
	PostInvent(item models.InventoryItem) error
	GetInvent() ([]models.InventoryItem, error)
	GetInventId(string) (models.InventoryItem, error)
	PutInventId(string, models.InventoryItem) error
	DeleteInventId(string) error
}

func (i *inventoryService) PostInvent(item models.InventoryItem) error {
	data, _ := i.in.ReadAllInvent()

	for _, k := range data {
		if k.IngredientID == item.IngredientID {
			return models.ErrExists
		} else if item.IngredientID == "" {
			return models.ErrNotEnough
		}
	}
	return i.in.WriteAllInvent(append(data, item))
}

func (i *inventoryService) GetInvent() ([]models.InventoryItem, error) {
	data, err := i.in.ReadAllInvent()

	if err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func (i *inventoryService) GetInventId(id string) (models.InventoryItem, error) {
	data, err := i.in.ReadAllInvent()
	if err != nil {
		return models.InventoryItem{}, err
	}

	for _, k := range data {
		if k.IngredientID == id {
			return k, nil
		}
	}

	return models.InventoryItem{}, models.ErrNotFound
}

func (i *inventoryService) PutInventId(id string, invent models.InventoryItem) error {
	data, err := i.in.ReadAllInvent()
	if err != nil {
		return err
	}
	var newinvent []models.InventoryItem
	for j, k := range data {
		if k.IngredientID == id {
			invent.IngredientID = k.IngredientID
			newinvent = append(data[:j], data[j+1:]...)
			return i.in.WriteAllInvent(append(newinvent, invent))
		}
	}

	return models.ErrNotFound
}

func (i *inventoryService) DeleteInventId(id string) error {
	data, err := i.in.ReadAllInvent()
	if err != nil {
		return err
	}

	for j, k := range data {
		if k.IngredientID == id {
			return i.in.WriteAllInvent(append(data[:j], data[j+1:]...))
		}
	}

	return models.ErrNotFound
}
