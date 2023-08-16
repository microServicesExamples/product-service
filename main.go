package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
)

var products = make(map[string]Product)

type ProductCategory string

const (
	PremiumProduct ProductCategory = "premium"
	RegularProduct ProductCategory = "regular"
	BudgetProduct  ProductCategory = "budget"
)

type Product struct {
	ID          string // todo: convert to int later
	Description string
	Name        string
	Category    ProductCategory // todo: convert to enum
	Price       float64
	Quantity    int64  // todo: see if we can move it another table
	CreatedAt   string // todo: see if we have any datetime variables
	UpdatedAt   string // todo: see if we have any datetime variables
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Pong"))
}

type CreateProductRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Category    ProductCategory `json:"category"` // todo: convert to enum
	Price       float64         `json:"price"`
	Quantity    int64           `json:"quantity"`
}

func (cpReq *CreateProductRequest) Validate() (err error) {
	// Validate product name
	productName := strings.TrimSpace(cpReq.Name)
	matched, err := regexp.MatchString(`^[a-zA-Z]+[_a-zA-Z\s\d]*$`, productName)
	if err != nil || !matched {
		fmt.Println("invalid product name")
		return errors.New("invalid product name")
	}

	if !(len(productName) > 0 && len(productName) < 30) {
		fmt.Println("product name too short or long")
		return errors.New("product name too short or long")
	}
	cpReq.Name = productName

	// Validate product description
	productDescription := strings.TrimSpace(cpReq.Description)
	if !(len(productDescription) > 5 && len(productDescription) < 50) {
		fmt.Println("product description too short or long")
		return errors.New("product description too short or long")
	}

	// Validate product category is one of the required product types
	if !(cpReq.Category == PremiumProduct || cpReq.Category == RegularProduct || cpReq.Category == BudgetProduct) {
		fmt.Println("invalid product category")
		return errors.New("invalid product category")
	}

	// Validate if product price is greater than 0
	if cpReq.Price <= 0 {
		fmt.Println("product price must be greater than 0")
		return errors.New("product price must be greater than 0")
	}

	// Validate if product quantity is non-negative
	if cpReq.Quantity < 0 {
		fmt.Println("product quantiy must be non-negative integer")
		return errors.New("product quantiy must be non-negative integer")
	}

	return nil
}

type CreateProductResponse struct {
	ID          string          `json:"id"` // todo: convert to int later
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Category    ProductCategory `json:"category"` // todo: convert to enum
	Price       float64         `json:"price"`
	Quantity    int64           `json:"quantity"`   // todo: see if we can move it another table
	CreatedAt   string          `json:"created_at"` // todo: see if we have any datetime variables
	UpdatedAt   string          `json:"updated_at"` // todo: see if we have any datetime variables
}

func AddProductHandler(w http.ResponseWriter, r *http.Request) {
	var cpReq CreateProductRequest
	err := json.NewDecoder(r.Body).Decode(&cpReq)
	if err != nil {
		fmt.Println("error unmashiling the request body, err:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request Body"))
		return
	}

	if err = cpReq.Validate(); err != nil {
		fmt.Println("error validating the request body, err:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Validate if the product already exisits
	for _, product := range products {
		if product.Name == cpReq.Name {
			fmt.Println("product already exisits in db")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Product already exisits in db"))
			return
		}
	}

	currentTime := time.Now().UTC().String()
	p := Product{
		ID:          uuid.New(),
		Name:        cpReq.Name,
		Description: cpReq.Description,
		Category:    cpReq.Category,
		Price:       cpReq.Price,
		Quantity:    cpReq.Quantity,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}
	products[p.ID] = p
	fmt.Println("success creating the product:", p)

	cpResp := CreateProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Category:    p.Category,
		Price:       p.Price,
		Quantity:    p.Quantity,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	resp, err := json.Marshal(cpResp)
	if err != nil {
		fmt.Println("error mashiling the response, err:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	var productList []CreateProductResponse

	for _, p := range products {
		productList = append(productList, CreateProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Category:    p.Category,
			Price:       p.Price,
			Quantity:    p.Quantity,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	resp, err := json.Marshal(productList)
	if err != nil {
		fmt.Println("error mashiling the response, err:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func GetProductDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["product_id"]

	p, ok := products[productId]
	if !ok {
		fmt.Println("product with id:", productId, "does not exist")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("product with id: %v does not exist", productId)))
		return
	}

	productDetails := CreateProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Category:    p.Category,
		Price:       p.Price,
		Quantity:    p.Quantity,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	resp, err := json.Marshal(productDetails)
	if err != nil {
		fmt.Println("error mashiling the response, err:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["product_id"]

	_, ok := products[productId]
	if !ok {
		fmt.Println("product with id:", productId, "does not exist")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("product with id: %v does not exist", productId)))
		return
	}

	delete(products, productId)
	w.WriteHeader(http.StatusNoContent)
}

type IncreseProductQuantityRequest struct {
	Quantity int64 `json:"quantity"`
}

func (iqReq *IncreseProductQuantityRequest) Validate() (err error) {
	// Validate if product quantity is non-negative
	if iqReq.Quantity < 0 {
		fmt.Println("product quantiy must be non-negative integer")
		return errors.New("product quantiy must be non-negative integer")
	}
	return nil
}

func IncreaseProductQuantityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["product_id"]

	var iqReq IncreseProductQuantityRequest
	err := json.NewDecoder(r.Body).Decode(&iqReq)
	if err != nil {
		fmt.Println("error unmashiling the request body, err:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request Body"))
		return
	}

	if err = iqReq.Validate(); err != nil {
		fmt.Println("error validating the request body, err:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	product, ok := products[productId]
	if !ok {
		fmt.Println("product with id:", productId, "does not exist")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("product with id: %v does not exist", productId)))
		return
	}

	updatedProduct := Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Category:    product.Category,
		Price:       product.Price,
		Quantity:    iqReq.Quantity,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   time.Now().UTC().String(),
	}
	products[product.ID] = updatedProduct
	fmt.Println("success updating the product:", product)

	iqResp := CreateProductResponse{
		ID:          updatedProduct.ID,
		Name:        updatedProduct.Name,
		Description: updatedProduct.Description,
		Category:    updatedProduct.Category,
		Price:       updatedProduct.Price,
		Quantity:    updatedProduct.Quantity,
		CreatedAt:   updatedProduct.CreatedAt,
		UpdatedAt:   updatedProduct.UpdatedAt,
	}
	resp, err := json.Marshal(iqResp)
	if err != nil {
		fmt.Println("error mashiling the response, err:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ping", PingHandler).Methods(http.MethodGet)

	s := r.PathPrefix("/products").Subrouter()
	s.HandleFunc("", AddProductHandler).Methods(http.MethodPost)
	s.HandleFunc("", GetProductsHandler).Methods(http.MethodGet)
	s.HandleFunc("/{product_id}", GetProductDetailsHandler).Methods(http.MethodGet)
	s.HandleFunc("/{product_id}", DeleteProductHandler).Methods(http.MethodDelete)
	s.HandleFunc("/{product_id}/quantity", IncreaseProductQuantityHandler).Methods(http.MethodPut)

	fmt.Println("Staring rest api server")
	go http.ListenAndServe(":8080", r)

	startGRPCServer()

}

