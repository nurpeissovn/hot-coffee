package dal

import (
	"encoding/json"
	"os"

	"hot-coffee/models"
)

type inventoryRepo struct {
	path string
}

func ReturnInvent(p string) *inventoryRepo {
	return &inventoryRepo{path: p}
}

type InventDal interface {
	ReadAllInvent() ([]models.InventoryItem, error)
	WriteAllInvent([]models.InventoryItem) error
}

func (p *inventoryRepo) ReadAllInvent() ([]models.InventoryItem, error) {
	file, err := os.Open(p.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var invent []models.InventoryItem
	if err := json.NewDecoder(file).Decode(&invent); err != nil {
		return nil, err
	}

	return invent, nil
}

func (p *inventoryRepo) WriteAllInvent(inventories []models.InventoryItem) error {
	file, err := os.Create(p.path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(&inventories)
}
