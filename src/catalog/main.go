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
	Description string
}

var catalog = []Product {
	Product{
		Name: "Milk",
		Size: 9,
		Unit: "ounces",
		Description: "A creamy, whole-fat treat. Produced by the finest cows in California.",
	},
	Product{
		Name: "Lemonade",
		Size: 750,
		Unit: "milliliters",
		Description: "Summer's favorite refreshment. Fresh lemon sourness with just the right amount of sweetness. This will take you back to the lemonade stands of your childhood.",
	},
	Product{
		Name: "Mexican Coke",
		Size: 12,
		Unit: "ounces",
		Description: "The classic. Made with real sugar in a glass bottle.",
	},
	Product{
		Name: "Water",
		Size: 200,
		Unit: "milliliters",
		Description: "Fresh, cold, calorie-free refreshment. Bottled at the source.",
	},
	Product{
		Name: "Fanta",
		Size: 12,
		Unit: "ounces",
		Description: "The international favorite. Bright flavor and bright color.",
	},
	Product{
		Name: "Sprite",
		Size: 16,
		Unit: "ounces",
		Description: "Lemon-lime refreshment. A tasty beverage perfect for washing down a cheeseburger.",
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
