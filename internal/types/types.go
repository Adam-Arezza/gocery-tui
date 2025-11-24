package types

import(
    "fmt"
)

type GroceryItem struct {
    Id       int     `json:"id"`
    Name     string  `json:"item_name"`
    Price    float32 `json:"unit_price"`
    Stock    int     `json:"stock"`
    Category int     `json:"category_id"`
}

func (g GroceryItem) Title() string       { return g.Name }
func (g GroceryItem) Description() string { return fmt.Sprintf("Stock: %d | Price: $%.2f", g.Stock, g.Price) }
func (g GroceryItem) FilterValue() string { return g.Name }

type PurchaseRequestItem struct {
    ItemId int `json:"item_id"`
    Stock int `json:"stock"`
}

type CartItem struct {
    Id int
    Name string  
    Price float32 
    Stock int     
    Quantity int
}

func (c CartItem) Title() string       { return c.Name }
func (c CartItem) Description() string { return fmt.Sprintf("Qty: %d | Price: $%.2f", c.Quantity, c.Price) }
func (c CartItem) FilterValue() string { return c.Name }

