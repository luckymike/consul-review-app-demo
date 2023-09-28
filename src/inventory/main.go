package main

import (
	"context"
	"net/http"
	"encoding/json"


	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var inventory = map[string] int {
	"Milk": 40,
	"Lemonade":  144,
	"Mexican Coke": 60,
	"Water": 20,
	"Fanta": 100,
	"Sprite": 50,
}


func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Route("/{item}", func(r chi.Router) {
		r.Use(itemCtx)
		r.Get("/", getInventory)
	})
	http.ListenAndServe(":3000", r)
}

func itemCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r * http.Request) {
		itemName := chi.URLParam(r, "item")
		count, ok := inventory[itemName]
		if !ok {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "count", count)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getInventory(w http.ResponseWriter, r * http.Request) {
	ctx := r.Context()
	count := ctx.Value("count")
	data, err := json.Marshal(count)
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	w.Write(data)
}
