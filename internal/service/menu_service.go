package service

import (
	"errors"

	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type menuService struct {
	dm dal.MenuDalInter
	di dal.InventDal
}

func ReturnMenuService(d dal.MenuDalInter, d2 dal.InventDal) *menuService {
	return &menuService{dm: d, di: d2}
}

type MenuServiceInter interface {
	GetAllMenu() ([]models.MenuItem, error)
	MenuId(id string) (models.MenuItem, error)
	DeleteId(id string) error
	PostMenu(menu models.MenuItem) error
	PutMenu(id string, updateMenu models.MenuItem) error
}

func (m *menuService) GetAllMenu() ([]models.MenuItem, error) {
	return m.dm.ReadAllMenu()
}

func (m *menuService) MenuId(id string) (models.MenuItem, error) {
	menu, _ := m.dm.ReadAllMenu()
	var updated models.MenuItem
	for _, k := range menu {
		if k.ID == id {
			return k, nil
		}
	}

	return updated, errors.New("item not found in menu")
}

func (m *menuService) DeleteId(id string) error {
	menu, _ := m.dm.ReadAllMenu()
	for i, k := range menu {
		if k.ID == id {
			return m.dm.WriteAllMenu(append(menu[:i], menu[i+1:]...))
		}
	}

	return models.ErrNotFound
}

func (m *menuService) PostMenu(menu models.MenuItem) error {
	d, _ := m.dm.ReadAllMenu()

	for _, k := range d {
		if k.ID == menu.ID {
			return models.ErrExists
		} else if menu.ID == "" {
			return models.ErrNotEnough
		}
	}
	ins, _ := m.di.ReadAllInvent()
	for _, k := range menu.Ingredients {
		var ithas bool
		for _, j := range ins {
			if j.IngredientID == k.IngredientID {
				ithas = true
				break
			}
		}
		if !ithas {
			return models.ErrNotFound
		}
	}

	return m.dm.WriteAllMenu(append(d, menu))
}

func (m *menuService) PutMenu(id string, newMenu models.MenuItem) error {
	menu, _ := m.dm.ReadAllMenu()

	for i, k := range menu {
		if k.ID == id {
			newMenu.ID = k.ID
			n := append(menu[:i], menu[i+1:]...)
			return m.dm.WriteAllMenu(append(n, newMenu))

		}
	}

	return models.ErrNotFound
}
