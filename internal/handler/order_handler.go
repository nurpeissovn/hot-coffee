package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type OrderHandler struct {
	OrderService service.OrderServiceInt
	Error        models.Error
}

func (h *OrderHandler) PostHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var order models.Order
	m := make(map[string]string)

	err := decoder.Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err_ := h.OrderService.Post(order, m)
	if err_ == nil {
		m := map[string]string{"message": "Order successfully created"}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", " ")
		err_ := encoder.Encode(m)
		if err_ != nil {
			fmt.Fprintf(os.Stderr, "Error:", err_)
		}

		logger.Info("Order successfully created", "method", "POST", "status", 201)
	} else if errors.Is(err_, models.ErrorNotFound) {
		h.OrderService.PrintErrorMessage(w, req, http.StatusBadRequest, "Invalid product id in order items")
		logger.Error("Invalid product id", "method", "POST", "status", 400)
	} else if errors.Is(err_, models.ErrorQuantity) {
		h.OrderService.PrintErrorMessage(w, req, http.StatusBadRequest, "Insufficient inventory for ingredient"+" '"+m["ingredientID"]+"'. Required: "+m["quantityRequired"]+m["unit"]+", Available: "+m["quantityAvailable"]+m["unit"]+".")
		logger.Error("Insufficient inventory for ingredient", "method", "POST", "status", 400)
	}
}

func (h *OrderHandler) GetHandler(w http.ResponseWriter, req *http.Request) {
	orders := h.OrderService.Get()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err := encoder.Encode(orders)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error:", err)
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Retrieve all orders", "method", "GET", "status", 200)
}

func (h *OrderHandler) GetIdHandler(w http.ResponseWriter, req *http.Request) {
	var id string = req.PathValue("id")
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	order, err := h.OrderService.GetId(id)
	if errors.Is(err, models.ErrorNotFound) {
		h.OrderService.PrintErrorMessage(w, req, http.StatusBadRequest, "Invalid order id")
		logger.Error("Invalid order id", "method", "GET", "status", 400)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err_ := encoder.Encode(order)
	if err_ != nil {
		fmt.Fprintf(os.Stderr, "Error:", err_)
		return
	}

	logger.Info("Retrieve a specific order", "method", "GET", "status", 200)
}

func (h *OrderHandler) PutHandler(w http.ResponseWriter, req *http.Request) {
	var id string = req.PathValue("id")

	decoder := json.NewDecoder(req.Body)
	var order models.Order

	err := decoder.Decode(&order)
	if err != nil {
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	m := make(map[string]string)
	err_ := h.OrderService.Put(id, order, m)
	if err_ == nil {
		m := map[string]string{"message": "Order successfully updated"}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", " ")
		err_ := encoder.Encode(m)
		if err_ != nil {
			fmt.Fprintf(os.Stderr, "Error:", err_)
		}

		logger.Info("Order updated", "method", "PUT", "status", 200)
	} else if errors.Is(err_, models.ErrorNotFound) {
		h.OrderService.PrintErrorMessage(w, req, http.StatusBadRequest, "Invalid order id or invalid product id in order items")
		logger.Error("Invalid order id or invalid product id", "method", "PUT", "status", 400)
	} else if errors.Is(err_, models.ErrorQuantity) {
		h.OrderService.PrintErrorMessage(w, req, http.StatusBadRequest, "Insufficient inventory for ingredient"+" '"+m["ingredientID"]+"'. Required: "+m["quantityRequired"]+m["unit"]+", Available: "+m["quantityAvailable"]+m["unit"]+".")
		logger.Error("Insufficient inventory for ingredient", "method", "PUT", "status", 400)
	} else if errors.Is(err_, models.ErrorConflict) {
		h.OrderService.PrintErrorMessage(w, req, http.StatusConflict, "Order already closed")
		logger.Error("Failed to update order", "method", "PUT", "status", 409)
	} else if errors.Is(err_, models.ErrorQuantityLess) {
		h.OrderService.PrintErrorMessage(w, req, http.StatusBadRequest, "Invalid quantity")
		logger.Error("Failed to update order", "method", "PUT", "status", 400)
	}
}

func (h *OrderHandler) DeleteHandler(w http.ResponseWriter, req *http.Request) {
	var id string = req.PathValue("id")
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := h.OrderService.Delete(id)
	if errors.Is(err, models.ErrorNotFound) {
		h.OrderService.PrintErrorMessage(w, req, http.StatusBadRequest, "Invalid order id")
		logger.Error("Failed to delete order", "method", "DELETE", "status", 400)
		return
	}

	m := map[string]string{"message": "Order successfully deleted"}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err_ := encoder.Encode(m)
	if err_ != nil {
		fmt.Fprintf(os.Stderr, "Error:", err_)
	}

	logger.Info("Order deleted", "method", "DELETE", "status", 200)
}

func (h *OrderHandler) PostCloseHandler(w http.ResponseWriter, req *http.Request) {
	var id string = req.PathValue("id")
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := h.OrderService.PostClose(id)
	if errors.Is(err, models.ErrorConflict) {
		h.OrderService.PrintErrorMessage(w, req, http.StatusConflict, "Order already closed")
		logger.Error("Failed to close order", "method", "POST", "status", 409)
		return
	} else if errors.Is(err, models.ErrorNotFound) {
		h.OrderService.PrintErrorMessage(w, req, http.StatusBadRequest, "Invalid order id")
		logger.Error("Failed to close order", "method", "POST", "status", 400)
		return
	}

	m := map[string]string{"message": "Order successfully closed"}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err_ := encoder.Encode(m)
	if err_ != nil {
		fmt.Fprintf(os.Stderr, "Error:", err_)
	}

	logger.Info("Order closed", "method", "POST", "status", 201)
}

func (h *OrderHandler) GetTotalSalesHandler(w http.ResponseWriter, req *http.Request) {
	m := make(map[string]float64)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	h.OrderService.GetTotalSales(m)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err := encoder.Encode(m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error:", err)
	}

	logger.Info("Total Sales", "method", "GET", "status", 200)
}

func (h *OrderHandler) GetPopularItemsHandler(w http.ResponseWriter, req *http.Request) {
	m := make(map[string]float64)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	array := h.OrderService.GetPopularItems(m)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err := encoder.Encode(array)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error:", err)
	}

	logger.Info("Popular Items", "method", "GET", "status", 200)
}
