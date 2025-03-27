package handler

import (
	"encoding/json"
	"net/http"

	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type inventHand struct {
	ser service.InventServiceInter
}

func ReturnInventHand(s service.InventServiceInter) *inventHand {
	return &inventHand{ser: s}
}

func (i *inventHand) PostInvent(w http.ResponseWriter, r *http.Request) {
	dataBody := r.Body
	var invent models.InventoryItem
	json.NewDecoder(dataBody).Decode(&invent)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	if err := i.ser.PostInvent(invent); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(map[string]string{"error": err.Error()})
		Handler.Error(err.Error(), "Method", "GET", "Status", 400)

	} else {
		w.WriteHeader(http.StatusCreated)
		encoder.Encode(map[string]string{"Message": models.SuccesMsg})
		Handler.Info("Success", "Method", "GET", "Status", 200)

	}
}

func (i *inventHand) GetInvent(w http.ResponseWriter, r *http.Request) {
	res, err := i.ser.GetInvent()
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	if err != nil {
		encoder.Encode(&err)
		Handler.Error(err.Error(), "Method", "GET", "Status", 404)

	} else {
		encoder.Encode(&res)
		Handler.Info("Success", "Method", "GET", "Status", 200)

	}
}

func (i *inventHand) GetInventId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("ID")
	data, err := i.ser.GetInventId(id)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		encoder.Encode(map[string]string{"error": err.Error()})
		Handler.Error(err.Error(), "Method", "GET", "Status", 404)

	} else {
		w.WriteHeader(http.StatusOK)
		encoder.Encode(&data)
		Handler.Info("Success", "Method", "GET", "Status", 200)

	}
}

func (i *inventHand) PutInventId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("ID")
	item := r.Body
	var invent models.InventoryItem
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	if err := json.NewDecoder(item).Decode(&invent); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(map[string]string{"error": err.Error()})
		Handler.Error(err.Error(), "Method", "GET", "Status", 500)

	} else if err := i.ser.PutInventId(id, invent); err != nil {
		w.WriteHeader(http.StatusNotFound)
		encoder.Encode(map[string]string{"error": err.Error()})
		Handler.Error(err.Error(), "Method", "GET", "Status", 404)

	} else {
		w.WriteHeader(http.StatusOK)
		encoder.Encode(map[string]string{"Message": models.SuccesMsg})
		Handler.Info("Success", "Method", "GET", "Status", 200)

	}
}

func (i *inventHand) DeleteInventId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("ID")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	if err := i.ser.DeleteInventId(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(map[string]string{"error": err.Error()})
		Handler.Error(err.Error(), "Method", "GET", "Status", 404)

	} else {
		w.WriteHeader(http.StatusOK)
		encoder.Encode(map[string]string{"Message": models.SuccesDeleteMsg})
		Handler.Info("Success", "Method", "GET", "Status", 200)

	}
}
