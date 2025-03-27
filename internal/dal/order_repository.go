package dal

import (
	"encoding/json"
	"fmt"
	"os"

	"hot-coffee/models"
)

type OrderRepository interface {
	WriteData(orders []models.Order)
	GetData() []models.Order
}

type Order struct {
	path string
}

func InitOrder(s string) *Order {
	order := Order{path: s}

	return &order
}

func (o *Order) WriteData(orders []models.Order) {
	file, err := os.Create(o.path)
	defer file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error:", err)
		return
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	err_ := encoder.Encode(orders)
	if err_ != nil {
		fmt.Fprintf(os.Stderr, "Error:", err_)
		return
	}
}

func (o *Order) GetData() []models.Order {
	file, _ := os.Open(o.path)
	defer file.Close()

	decoder := json.NewDecoder(file)

	filteredData := []models.Order{}

	decoder.Token()

	for decoder.More() {
		data := models.Order{}
		decoder.Decode(&data)
		filteredData = append(filteredData, data)
	}

	return filteredData
}
