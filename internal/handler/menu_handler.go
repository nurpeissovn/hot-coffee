package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type menuHand struct {
	ser service.MenuServiceInter
}

func ReturnMenuHand(s service.MenuServiceInter) *menuHand {
	return &menuHand{ser: s}
}

var Handler = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func (h *menuHand) GetAllMenu(w http.ResponseWriter, r *http.Request) {
	d, _ := h.ser.GetAllMenu()
	jsonData, _ := json.MarshalIndent(d, " ", "  ")
	w.Write(jsonData)
	Handler.Info("Success", "Method", "GET", "Status", 200)
}

func (h *menuHand) GetMenuId(w http.ResponseWriter, r *http.Request) {
	d, err := h.ser.MenuId(r.PathValue("ID"))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Item not found in the menu"))
		Handler.Error(err.Error(), "Method", "GET", "Status", 404)

	} else {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", " ")
		if err := encoder.Encode(d); err != nil {
			w.Write([]byte(err.Error()))
		}
		Handler.Info("Success", "Method", "GET", "Status", 200)

	}
}

func (h *menuHand) DeleteMenuId(w http.ResponseWriter, r *http.Request) {
	res := h.ser.DeleteId(r.PathValue("ID"))
	if res == models.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("The Item not exist"))
		Handler.Error(models.ErrNotFound.Error(), "Method", "GET", "Status", 404)

	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("successfulle deleted"))
		Handler.Info("Success", "Method", "GET", "Status", 200)

	}
}

func (h *menuHand) PostMenu(w http.ResponseWriter, r *http.Request) {
	data := r.Body
	var res models.MenuItem
	json.NewDecoder(data).Decode(&res)
	if err := h.ser.PostMenu(res); err == nil {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("The item successfully added"))
		Handler.Info("Success", "Method", "GET", "Status", 201)

	} else if err == models.ErrExists {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("The item already exists"))
		Handler.Error(models.ErrExists.Error(), "Method", "GET", "Status", 409)

	} else if err == models.ErrNotEnough {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Product ID is missing"))
		Handler.Error(err.Error(), "Method", "GET", "Status", 400)

	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Ingredients not found for menu"))
		Handler.Error(models.ErrNotFound.Error(), "Method", "GET", "Status", 404)

	}
}

func (h *menuHand) PutMenu(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("ID")
	bodyData := r.Body
	var res models.MenuItem
	json.NewDecoder(bodyData).Decode(&res)
	if h.ser.PutMenu(id, res) == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("The menu successfully updated"))
		Handler.Info("Success", "Method", "GET", "Status", 200)

	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("The item not found in menu to update"))
		Handler.Error(models.ErrNotFound.Error(), "Method", "GET", "Status", 404)

	}
}
