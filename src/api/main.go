package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Config struct {
	CatalogUrl string
	InventoryUrl string
}

type Stock struct {
	Product any
	Count any
}


var cfg = Config {
	CatalogUrl: os.Getenv("CATALOG_URL"),
	InventoryUrl: os.Getenv("INVENTORY_URL"),
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/", func(r chi.Router) {
		r.Use(productsCtx)
		r.Get("/", getProducts)
	})

  	r.Route("/{product}", func(r chi.Router) {
	 	r.Use(productCtx)
		r.Use(inventoryCtx)
		r.Get("/", getProduct)
	})
	http.ListenAndServe(":3000", r)
}

func productsCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r * http.Request) {

		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/", cfg.CatalogUrl), nil)
		if err != nil {
			fmt.Print(err.Error())
		}

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")

		for name, values := range r.Header {
			if strings.Contains(name, "X-Acme-") {
				for _, value := range values {
					req.Header.Add(name, value)
				}
			}

		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Print(err.Error())
		}
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Print(err.Error())
		}

		var responseObject []string
		json.Unmarshal(bodyBytes, &responseObject)

		ctx := context.WithValue(r.Context(), "products", responseObject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func productCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r * http.Request) {
		productName := chi.URLParam(r, "product")

		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/%s", cfg.CatalogUrl, productName), nil)
		if err != nil {
			fmt.Print(err.Error())
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")

		for name, values := range r.Header {
			if strings.Contains(name, "X-Acme-") {
				for _, value := range values {
					req.Header.Add(name, value)
				}
			}

		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Print(err.Error())
		}
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Print(err.Error())
		}

		var responseObject map[string]interface{}
		json.Unmarshal(bodyBytes, &responseObject)

		ctx := context.WithValue(r.Context(), "product", responseObject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func inventoryCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r * http.Request) {
		productName := chi.URLParam(r, "product")

		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/%s", cfg.InventoryUrl, productName), nil)
		if err != nil {
			fmt.Print(err.Error())
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Print(err.Error())
		}
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Print(err.Error())
		}

		var responseObject int
		json.Unmarshal(bodyBytes, &responseObject)

		ctx := context.WithValue(r.Context(), "inventory", responseObject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getProducts(w http.ResponseWriter, r * http.Request) {
	ctx := r.Context()
	products := ctx.Value("products")
	data, err := json.Marshal(products)
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	w.Write(data)
}

func getProduct(w http.ResponseWriter, r * http.Request) {
	ctx := r.Context()
	stock := Stock{
		Product: ctx.Value("product").(any),
		Count: ctx.Value("inventory").(any),
	}

	data, err := json.Marshal(stock)
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	w.Write(data)
}
