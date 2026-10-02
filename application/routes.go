package application

import (
	"net/http"

	"github.com/dreamsofcode-io/orders-api/handler"
	"github.com/dreamsofcode-io/orders-api/repository/order"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


func (a *App) loadRoutes()  {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/orders", a.loadOrderRoutes)


	a.router = router
}

func (a *App) loadOrderRoutes(router chi.Router) {
	orderHander := &handler.Order{
		Repo: &order.RedisRepo{
			Client: a.rdb,
		},

	}

	router.Post("/", orderHander.Create)
	router.Get("/", orderHander.List)
	router.Get("/{id}", orderHander.GetById)
	router.Put("/{id}", orderHander.UpdateById)
	router.Delete("/{id}", orderHander.DeleteById)



}

