package messages

import (
    "github.com/Adam-Arezza/gocery-tui/internal/types"
)

type NewCartItemMsg struct{
    Item types.CartItem
}

type CompletePurchaseMsg struct {}

type PurchaseError struct {}

type PurchaseMsg struct{}

type CloseModalMsg struct{}


