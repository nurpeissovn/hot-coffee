package models

import "errors"

type Order struct {
	ID           string      `json:"order_id"`
	CustomerName string      `json:"customer_name"`
	Items        []OrderItem `json:"items"`
	Status       string      `json:"status"`
	CreatedAt    string      `json:"created_at"`
}

type OrderItem struct {
	ProductId string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

var (
	ErrorNotFound     = errors.New("ProductId Not Found")
	ErrorQuantity     = errors.New("Ingredients Quantity")
	ErrorConflict     = errors.New("Already updated")
	ErrorQuantityLess = errors.New("Quantity less")
)

type Error struct {
	error_ string `json:"error"`
}
