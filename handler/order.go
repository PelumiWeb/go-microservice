package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"uuid"

	"github.com/dreamsofcode-io/orders-api/model"
	"github.com/dreamsofcode-io/orders-api/repository/order"
	"github.com/go-chi/chi/v5"
)



type Order struct {
	Repo *order.RedisRepo
}


func (h *Order) Create(w http.ResponseWriter, r *http.Request ) {
	var body struct {
		CustomerID uuid.UUID `json:"customer_id"`
		LineItems []model.LineItem `json:"line_items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	order := model.Order{
		OrderId: rand.Uint64(),
		CustomerID: body.CustomerID,
		LineItems: body.LineItems,
		CreatedAt: &now,
	}

	if err := h.Repo.Insert(r.Context(), order); err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (h *Order) List(w http.ResponseWriter, r *http.Request ) {
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"

	}
	const descimal = 10
	const bitSize = 64

	cursor, err := strconv.ParseUint(cursorStr, descimal, bitSize)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50
	res, err := h.Repo.FindAll(r.Context(), order.FindAllPage{
		Offset: uint(cursor),
		Size: size,
	})

	if err != nil {
		fmt.Println("failed to find all:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Items []model.Order `json:"items"`
		Next uint64 `json:"next,omitempty"`
	}

	response.Items = res.Orders
	response.Next = res.Cursor

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed to marshall", err)
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}


	w.Write(data)



	fmt.Println("List all orders")
}

func (h *Order) GetById(w http.ResponseWriter, r *http.Request ) {

	idParam := chi.URLParam(r, "id")
	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	o, err := h.Repo.FindById(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}

	if err := json.NewEncoder(w).Encode(o); err != nil {
		fmt.Println("failed to marshall")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}



	fmt.Println("Get by ID")
}

func (h *Order) UpdateById(w http.ResponseWriter, r *http.Request ) {
	var body struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idParam := chi.URLParam(r, "id")


	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	theOrder, err := h.Repo.FindById(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	const completedStatus = "completed"
	const shippedStatus = "shipped"
	now := time.Now().UTC()

	switch body.Status {
	case shippedStatus:
		if theOrder.ShippedAt != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.ShippedAt = &now
	case completedStatus:
		if theOrder.CompletedAt != nil  || theOrder.ShippedAt == nil{
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.CompletedAt = &now
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.Repo.Update(r.Context(), theOrder)
	if err != nil {
		fmt.Println("failed to insert", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(theOrder); err != nil {
		fmt.Println("failed to marshall:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


	



	fmt.Println("Update by ID")
}

func (h *Order) DeleteById(w http.ResponseWriter, r *http.Request ) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.Repo.DeleteById(r.Context(), orderID) 
	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	}else if err != nil {
		fmt.Println("falied to fing id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}



	fmt.Println("Delete by ID")
}


