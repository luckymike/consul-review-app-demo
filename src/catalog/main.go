package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Product struct {
	Name string
	Size int
	Unit string
}

var catalog = []Product {
	Product{
		Name: "Milk",
		Size: 9,
		Unit: "ounces",
	},
	Product{
		Name: "Lemonade",
		Size: 750,
		Unit: "milliliters",
	},
	Product{
		Name: "Mexican Coke",
		Size: 12,
		Unit: "ounces",
	},
	Product{
		Name: "Water",
		Size: 200,
		Unit: "milliliters",
	},
	Product{
		Name: "Fanta",
		Size: 12,
		Unit: "ounces",
	},
	Product{
		Name: "Sprite",
		Size: 16,
		Unit: "ounces",
	},
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/", func(r chi.Router) {
		r.Get("/", getCatalog)
	})

	r.Route("/{product}", func(r chi.Router) {
		r.Use(productCtx)
		r.Get("/", getProduct)
	})
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), r)
}

func productCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r * http.Request) {
		productName := chi.URLParam(r, "product")
		idx := productByName(productName, catalog)
		if idx == -1  {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		product := catalog[idx]
		ctx := context.WithValue(r.Context(), "product", product)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getCatalog(w http.ResponseWriter, r * http.Request) {
	var names []string
	for _, product := range catalog {
		names = append(names, product.Name)
	}
	data, err := json.Marshal(names)
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	w.Write(data)
}

func getProduct(w http.ResponseWriter, r * http.Request) {
	ctx := r.Context()
	product, ok := ctx.Value("product").(Product)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
 	data, err := json.Marshal(product)
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	w.Write(data)
}

func productByName(n string, s []Product) int {
	var index int = -1
	for i, product := range s {
		if product.Name == n {
			index = i
			break
		}
	}
	return index
}
