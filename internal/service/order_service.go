package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type OrderServiceInt interface {
	Post(models.Order, map[string]string) error
	Get() []models.Order
	GetId(string) (models.Order, error)
	Put(string, models.Order, map[string]string) error
	Delete(string) error
	PostClose(string) error
	GetTotalSales(map[string]float64)
	GetPopularItems(map[string]float64) []map[string]float64
	PrintErrorMessage(http.ResponseWriter, *http.Request, int, string)
}

type OrderService struct {
	OrderRepository  dal.OrderRepository
	MenuRepository   dal.MenuDalInter
	InventRepository dal.InventDal
}

func checkMenu(order models.Order, menues []models.MenuItem) error {
	for i := range order.Items {
		var isExist bool
		for j := range menues {
			if order.Items[i].ProductId == menues[j].ID {
				isExist = true
				break
			}
		}
		if !isExist {
			return models.ErrorNotFound
		}
	}
	return nil
}

func checkIngredients(orderItem models.OrderItem, menueIngredients []models.MenuItemIngredient, ingredients *[]models.InventoryItem, m map[string]string) error {
	for i := range menueIngredients {
		for j := range *ingredients {
			if menueIngredients[i].IngredientID == (*ingredients)[j].IngredientID {
				quantity := float64(orderItem.Quantity) * menueIngredients[i].Quantity
				if (*ingredients)[j].Quantity-quantity < 0 {
					m["ingredientID"] = (*ingredients)[j].IngredientID
					m["quantityRequired"] = strconv.Itoa(int(quantity))
					m["quantityAvailable"] = strconv.Itoa(int((*ingredients)[j].Quantity))
					m["unit"] = (*ingredients)[j].Unit
					return models.ErrorQuantity
				}
				(*ingredients)[j].Quantity = (*ingredients)[j].Quantity - quantity
			}
		}
	}

	return nil
}

func checkMenuPost(order models.Order, menues []models.MenuItem, ingredients *[]models.InventoryItem, m map[string]string) error {
	for i := range order.Items {
		for j := range menues {
			if order.Items[i].ProductId == menues[j].ID {
				err := checkIngredients(order.Items[i], menues[j].Ingredients, ingredients, m)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (o *OrderService) Post(order models.Order, m map[string]string) error {
	orders := o.OrderRepository.GetData()
	var id string = strconv.Itoa(len(orders) + 1)
	order.ID = id
	order.Status = "open"
	order.CreatedAt = time.Now().Format(time.RFC3339)
	orders = append(orders, order)

	menues, _ := o.MenuRepository.ReadAllMenu()
	ingredients, _ := o.InventRepository.ReadAllInvent()

	err := checkMenu(order, menues)
	if err != nil {
		return err
	}

	err = checkMenuPost(order, menues, &ingredients, m)
	if err != nil {
		return err
	}
	o.InventRepository.WriteAllInvent(ingredients)
	o.OrderRepository.WriteData(orders)

	return nil
}

func (o *OrderService) Get() []models.Order {
	data := o.OrderRepository.GetData()

	return data
}

func (o *OrderService) GetId(id string) (models.Order, error) {
	orders := o.OrderRepository.GetData()
	var order models.Order

	for i := range orders {
		order = orders[i]
		if order.ID == id {
			return order, nil
		}
	}

	var order_1 models.Order
	return order_1, models.ErrorNotFound
}

func checkIngredientsPut(orderItem, orderNewItem models.OrderItem, menueIngredients []models.MenuItemIngredient, ingredients *[]models.InventoryItem, m map[string]string) error {
	for i := range menueIngredients {
		for j := range *ingredients {
			if menueIngredients[i].IngredientID == (*ingredients)[j].IngredientID {
				if orderItem.Quantity > orderNewItem.Quantity {
					quantity := float64(orderItem.Quantity-orderNewItem.Quantity) * menueIngredients[i].Quantity
					(*ingredients)[j].Quantity = (*ingredients)[j].Quantity + quantity
				} else if orderItem.Quantity < orderNewItem.Quantity {
					quantity := float64(orderNewItem.Quantity-orderItem.Quantity) * menueIngredients[i].Quantity
					if (*ingredients)[j].Quantity-quantity < 0 {
						m["ingredientID"] = (*ingredients)[j].IngredientID
						m["quantityRequired"] = strconv.Itoa(int(quantity))
						m["quantityAvailable"] = strconv.Itoa(int((*ingredients)[j].Quantity))
						m["unit"] = (*ingredients)[j].Unit
						return models.ErrorQuantity
					}
					(*ingredients)[j].Quantity = (*ingredients)[j].Quantity - quantity
				}
			}
		}
	}

	return nil
}

func checkIngredientsPutUpdate(orderItem models.OrderItem, menueIngredients []models.MenuItemIngredient, ingredients *[]models.InventoryItem, m map[string]string) error {
	for i := range menueIngredients {
		for j := range *ingredients {
			if menueIngredients[i].IngredientID == (*ingredients)[j].IngredientID {
				quantity := float64(orderItem.Quantity) * menueIngredients[i].Quantity
				(*ingredients)[j].Quantity = (*ingredients)[j].Quantity + quantity
			}
		}
	}

	return nil
}

func checkMenuPut(orderItem, orderNewItem models.OrderItem, menues []models.MenuItem, ingredients *[]models.InventoryItem, m map[string]string) error {
	for j := range menues {
		if orderItem.ProductId == menues[j].ID {
			err := checkIngredientsPut(orderItem, orderNewItem, menues[j].Ingredients, ingredients, m)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func checkMenuPutNew(orderNewItem models.OrderItem, menues []models.MenuItem, ingredients *[]models.InventoryItem, m map[string]string) error {
	for j := range menues {
		if orderNewItem.ProductId == menues[j].ID {
			err := checkIngredients(orderNewItem, menues[j].Ingredients, ingredients, m)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func checkMenuPutNewUpdate(orderItem models.OrderItem, menues []models.MenuItem, ingredients *[]models.InventoryItem, m map[string]string) error {
	for j := range menues {
		if orderItem.ProductId == menues[j].ID {
			err := checkIngredientsPutUpdate(orderItem, menues[j].Ingredients, ingredients, m)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func updateOrder(order, orderNew models.Order, menues []models.MenuItem, ingredients *[]models.InventoryItem, m map[string]string) ([]models.OrderItem, error) {
	var result []models.OrderItem

	for i := range orderNew.Items {
		var isExist bool
		for j := range order.Items {
			if orderNew.Items[i].ProductId == order.Items[j].ProductId {
				if orderNew.Items[i].Quantity <= 0 {
					return result, models.ErrorQuantityLess
				}
				err := checkMenuPut(order.Items[j], orderNew.Items[i], menues, ingredients, m)
				if err != nil {
					return order.Items, err
				}
				order.Items[j].Quantity = orderNew.Items[i].Quantity
				isExist = true
				result = append(result, orderNew.Items[i])
			}
		}
		if !isExist {
			if orderNew.Items[i].Quantity <= 0 {
				return result, models.ErrorQuantityLess
			}
			err := checkMenuPutNew(orderNew.Items[i], menues, ingredients, m)
			if err != nil {
				return order.Items, err
			}
			order.Items = append(order.Items, orderNew.Items[i])
			result = append(result, orderNew.Items[i])
		}
	}

	for i := range order.Items {
		var isExist bool
		for j := range orderNew.Items {
			if orderNew.Items[j].ProductId == order.Items[i].ProductId {
				isExist = true
				break
			}
		}
		if !isExist {
			err := checkMenuPutNewUpdate(order.Items[i], menues, ingredients, m)
			if err != nil {
				return order.Items, err
			}
		}
	}

	return result, nil
}

func (o *OrderService) Put(id string, order models.Order, m map[string]string) error {
	orders := o.OrderRepository.GetData()
	menues, _ := o.MenuRepository.ReadAllMenu()
	ingredients, _ := o.InventRepository.ReadAllInvent()

	id_, _ := strconv.Atoi(id)
	err := checkMenu(order, menues)
	if err != nil || id_ < 0 {
		return err
	}

	for i := range orders {
		if orders[i].ID == id && orders[i].Status == "open" {
			orders[i].Items, err = updateOrder(orders[i], order, menues, &ingredients, m)
			if err != nil {
				return err
			}
			o.InventRepository.WriteAllInvent(ingredients)
			o.OrderRepository.WriteData(orders)
			return nil
		} else if orders[i].ID == id && orders[i].Status == "closed" {
			return models.ErrorConflict
		}
	}

	return models.ErrorNotFound
}

func checkIngredientsDelete(orderItem models.OrderItem, menueIngredients []models.MenuItemIngredient, ingredients *[]models.InventoryItem) {
	for i := range menueIngredients {
		for j := range *ingredients {
			if menueIngredients[i].IngredientID == (*ingredients)[j].IngredientID {
				quantity := float64(orderItem.Quantity) * menueIngredients[i].Quantity
				(*ingredients)[j].Quantity = (*ingredients)[j].Quantity + quantity
			}
		}
	}
}

func checkMenuDelete(order models.Order, menues []models.MenuItem, ingredients *[]models.InventoryItem) {
	for i := range order.Items {
		for j := range menues {
			if order.Items[i].ProductId == menues[j].ID {
				checkIngredientsDelete(order.Items[i], menues[j].Ingredients, ingredients)
			}
		}
	}
}

func (o *OrderService) Delete(id string) error {
	orders := o.OrderRepository.GetData()
	menues, _ := o.MenuRepository.ReadAllMenu()
	ingredients, _ := o.InventRepository.ReadAllInvent()
	var result []models.Order
	var order models.Order
	var isExist bool

	for i := range orders {
		order = orders[i]
		if order.ID != id {
			result = append(result, order)
		}
		if order.ID == id {
			isExist = true
			if order.Status == "open" {
				checkMenuDelete(order, menues, &ingredients)
				o.InventRepository.WriteAllInvent(ingredients)
			}
		}
	}

	if isExist {
		o.OrderRepository.WriteData(result)
		return nil
	}
	return models.ErrorNotFound
}

func (o *OrderService) PostClose(id string) error {
	orders := o.OrderRepository.GetData()
	var result []models.Order
	var order models.Order
	var isExist bool

	for i := range orders {
		order = orders[i]
		if order.ID == id && order.Status == "open" {
			order.Status = "closed"
			isExist = true
		} else if order.ID == id && order.Status == "closed" {
			return models.ErrorConflict
		}
		result = append(result, order)
	}

	if isExist {
		o.OrderRepository.WriteData(result)
		return nil
	}

	return models.ErrorNotFound
}

func getMenuPrice(order models.Order, menues []models.MenuItem, m map[string]float64) {
	for i := range order.Items {
		for j := range menues {
			if order.Items[i].ProductId == menues[j].ID {
				m[order.Items[i].ProductId] += float64(order.Items[i].Quantity) * menues[j].Price
			}
		}
	}
}

func (o *OrderService) GetTotalSales(m map[string]float64) {
	orders := o.OrderRepository.GetData()
	menues, _ := o.MenuRepository.ReadAllMenu()

	for i := range orders {
		if orders[i].Status == "closed" {
			getMenuPrice(orders[i], menues, m)
		}
	}
}

func sortItems(m map[string]float64) []map[string]float64 {
	var array []map[string]float64

	for k, v := range m {
		m1 := map[string]float64{k: v}
		array = append(array, m1)
	}

	var length int = len(array)

	for i := range array {
		for j := i + 1; j < length; j++ {
			var v float64
			var v1 float64
			for _, v2 := range array[i] {
				v = v2
			}
			for _, v3 := range array[j] {
				v1 = v3
			}
			if v < v1 {
				array[i], array[j] = array[j], array[i]
			}
		}
	}

	return array
}

func (o *OrderService) GetPopularItems(m map[string]float64) []map[string]float64 {
	orders := o.OrderRepository.GetData()

	for i := range orders {
		for j := range orders[i].Items {
			m[orders[i].Items[j].ProductId] += float64(orders[i].Items[j].Quantity)
		}
	}

	array := sortItems(m)
	return array
}

func (o *OrderService) PrintErrorMessage(w http.ResponseWriter, req *http.Request, h int, s string) {
	w.WriteHeader(h)
	w.Header().Set("Content-Type", "application/json")
	m := map[string]string{"error": s}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err := encoder.Encode(m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error:", err)
	}
}
