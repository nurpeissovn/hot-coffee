package dal

import (
	"encoding/json"
	"os"

	"hot-coffee/models"
)

type menuRepo struct {
	filepath string // menu.json
}

type MenuDalInter interface {
	ReadAllMenu() ([]models.MenuItem, error)
	WriteAllMenu([]models.MenuItem) error
}

func ReturnMenuRepo(path string) *menuRepo {
	return &menuRepo{filepath: path}
}

func (dm *menuRepo) ReadAllMenu() ([]models.MenuItem, error) {
	f, e := os.Open(dm.filepath)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	var menuItems []models.MenuItem
	if e := json.NewDecoder(f).Decode(&menuItems); e != nil {
		return nil, e
	}
	return menuItems, nil
}

func (p *menuRepo) WriteAllMenu(menus []models.MenuItem) error {
	f, e := os.Create(p.filepath)
	if e != nil {
		return e
	}
	defer f.Close()
	encode := json.NewEncoder(f)
	encode.SetIndent("", " ")
	return encode.Encode(&menus)
}
