package graphql

import "time"

type Account struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Orders []*Order `json:"orders"`
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

type OrderedProduct struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

type Order struct {
	ID         string            `json:"id"`
	CreatedAt  time.Time         `json:"createdAt"`
	TotalPrice float64           `json:"totalPrice"`
	Products   []*OrderedProduct `json:"products"`
}

type PaginationInput struct {
	Skip int `json:"skip"`
	Take int `json:"take"`
}

type AccountInput struct {
	Name string `json:"name"`
}

type ProductInput struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

type OrderedProductInput struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type OrderInput struct {
	AccountID string                 `json:"accountId"`
	Products  []*OrderedProductInput `json:"products"`
}
type Query struct {
}
