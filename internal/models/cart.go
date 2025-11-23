package models

import (
	"fmt"
	"net/http"
    "encoding/json"
    "bytes"
	"github.com/Adam-Arezza/gocery-tui/config"
	"github.com/Adam-Arezza/gocery-tui/internal/components"
	"github.com/Adam-Arezza/gocery-tui/internal/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type CartModel struct{
    Items list.Model
    Focused bool
    Width int
    Height int
    Selected bool
    Total float32
    Wallet float32
    ItemStyles *list.DefaultDelegate
    PurchaseModal *components.CartModal
    ShowPurchaseModal bool
    ServerConfig *config.ServerConfig
}

type CompletePurchaseMsg struct {}
type PurchaseError struct {}
type PurchaseRequestItem struct {
    ItemId int `json:"item_id"`
    Stock int `json:"stock"`
}

func (cart *CartModel) Init() tea.Cmd {
    cart.PurchaseModal = cart.NewPurchaseModal()
    return nil
}

func (cart *CartModel) Update(msg tea.Msg) (tea.Model, tea.Cmd){
    var cmd tea.Cmd
    switch msg := msg.(type){

        case tea.WindowSizeMsg:
            cart.Width = msg.Width / 2
            cart.Height = msg.Height
            cart.Items.SetSize(cart.Width, cart.Height-2)
            cart.PurchaseModal.Width = cart.Width
            cart.PurchaseModal.Height = cart.Height

        case NewCartItemMsg:
            var item components.CartItem
            item = msg.item
            cartItems := cart.Items.Items() 
            cartItems = append(cartItems, item)
            cart.Items.SetItems(cartItems)
            cart.Total = cart.cartSum()
            cart.PurchaseModal.TotalPrice = cart.Total
            return cart, cmd

        case components.PurchaseMsg:
            purchaseItems := cart.getPurchaseItems()
            cmd, err := cart.makePurchase(purchaseItems)
            if err != nil {
                return cart, func() tea.Msg{
                    return PurchaseError{}
                }
            }
            cart.ShowPurchaseModal = false
            cart.Wallet = cart.Wallet - cart.Total
            cart.updateTotal(0.00)
            cart.Items.SetItems([]list.Item{})
            cart.PurchaseModal.Wallet = cart.Wallet
            return cart, cmd

        case components.CloseModalMsg:
            cart.ShowPurchaseModal = false
            cart.PurchaseModal.InsufficientFunds = false
            return cart, nil

        case tea.KeyMsg:
            if cart.Focused{
                if cart.ShowPurchaseModal{
                    modal, cmd := cart.PurchaseModal.Update(msg)
                    if modal, ok := modal.(*components.CartModal); ok {
                        cart.PurchaseModal = modal
                    }
                    return cart, cmd
                }

                switch msg.String(){

                case "enter":
                    cart.Selected = !cart.Selected
                    return cart, cmd

                case "up", "k":
                    if cart.Selected{
                        item := cart.Items.SelectedItem().(components.CartItem)
                        idx := cart.Items.Index()
                        if item.Quantity >= 1 && item.Quantity < item.Stock{
                            item.Quantity++
                            cart.Items.Items()[idx] = item
                            cart.Items.SetItems(cart.Items.Items())
                            cart.updateTotal(cart.cartSum())
                            return cart, cmd
                        }
                    }

                case "down","j":
                    if cart.Selected{
                        item := cart.Items.SelectedItem().(components.CartItem)
                        idx := cart.Items.Index()
                        if item.Quantity >= 1{
                            item.Quantity--
                            cart.Items.Items()[idx] = item
                            cart.Items.SetItems(cart.Items.Items())
                            cart.updateTotal(cart.cartSum())
                            return cart, cmd
                        }
                    }

                case "d":
                    if cart.Selected{
                        idx := cart.Items.Index()
                        cart.Items.RemoveItem(idx)
                        cart.Selected = false
                        cart.updateTotal(cart.cartSum())
                        return cart, cmd
                    }

                case " ":
                    cart.ShowPurchaseModal = true
                    return cart, cmd

                case "esc":
                    if cart.ShowPurchaseModal{
                        cart.ShowPurchaseModal = false
                    }
                    return cart, cmd
            }
        }
    }

    if cart.Focused{
        cart.Items, cmd = cart.Items.Update(msg)
        return cart, cmd
    }

    return cart, cmd
}

func (cart *CartModel) View() string{
    if cart.Selected{
        cart.ItemStyles.Styles.SelectedDesc = styles.ActiveCartItemStyle
    }else{
        cart.ItemStyles.Styles.SelectedDesc = styles.GroceryItemSelectedDesc
    }
    list := cart.Items.View()
    cartTotal := styles.GroceryTotal.Render(fmt.Sprintf("TOTAL: $ %.2f",cart.Total))
    wallet := styles.GroceryTotal.Render(fmt.Sprintf("WALLET: $ %.2f", cart.Wallet))
    if cart.ShowPurchaseModal{
        return cart.PurchaseModal.View()
    }
    if cart.Focused{
        return styles.FocusedStyle.Height(cart.Height).Width(cart.Width).Render(list, cartTotal, wallet)
    }else{
        return styles.UnFocusedStyle.Height(cart.Height).Width(cart.Width).Render(list, cartTotal, wallet)
    }
}

func NewGroceryCart(cfg *config.ServerConfig) *CartModel {
    delegate := styles.GroceryDelegate{
        StyleDelegate: list.NewDefaultDelegate(),
     }

    delegate.StyleDelegate.Styles.NormalTitle = styles.GroceryItemNormalTitle
    delegate.StyleDelegate.Styles.NormalDesc  = styles.GroceryItemNormalDesc
    delegate.StyleDelegate.Styles.SelectedTitle = styles.GroceryItemSelectedTitle
    delegate.StyleDelegate.Styles.SelectedDesc  = styles.GroceryItemSelectedDesc

    groceryList := list.New([]list.Item{}, &delegate, 0, 0)
    groceryList.Styles.Title = styles.GroceryListTitle
    groceryList.Title = groceryList.Styles.Title.Render("MY CART")

    //Pagination style
    groceryList.Styles.PaginationStyle = styles.PaginatorStyle
    groceryList.Styles.InactivePaginationDot = styles.PaginatorStyle
    groceryList.Paginator.InactiveDot = "-" 
    groceryList.Paginator.ActiveDot = "*"

    //Help and HelpStyle style
    groceryList.Styles.HelpStyle = styles.HelpStyle
    groceryList.Help.Styles.ShortDesc = styles.HelpMenuStyle
    groceryList.Help.Styles.ShortKey = styles.HelpMenuStyle 
    
    return &CartModel{
        Items: groceryList,
        ItemStyles: &delegate.StyleDelegate,
        Wallet: 25.00,
        Total: 0.00,
        ShowPurchaseModal: false,
        Focused: false,
        Selected: false,
        ServerConfig: cfg,
    }
}

func (cart *CartModel) cartSum() float32 {
    items := cart.Items.Items()
    var sum float32
    for _, item := range items{
        if cartItem, ok := item.(components.CartItem); ok {
            sum += cartItem.Price * float32(cartItem.Quantity)
        }
    }

    return sum
}

func (cart *CartModel) updateTotal(total float32){
    cart.Total = total
    cart.PurchaseModal.TotalPrice = total
}

func (cart *CartModel) NewPurchaseModal() (*components.CartModal){
    modal := components.CartModal{
        Height: cart.Height/2,
        Width: cart.Width/2,
        Wallet: cart.Wallet,
        TotalPrice: cart.Total,
        InsufficientFunds: false,
        Confirm: false,
    }
    return &modal
}

func (cart *CartModel) makePurchase(items []PurchaseRequestItem) (tea.Cmd, error){
    url := "http://" + cart.ServerConfig.Host + ":" + cart.ServerConfig.Port + "/grocery_items"
    body := items

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
        return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)

    return func()tea.Msg{
        return CompletePurchaseMsg{}
    }, nil
}

func (cart *CartModel)getPurchaseItems()[]PurchaseRequestItem{
    var items []PurchaseRequestItem
    for _, item := range cart.Items.Items(){
        cartItem := item.(components.CartItem)
        newStock := cartItem.Stock - cartItem.Quantity
        newItem := PurchaseRequestItem{
            ItemId: cartItem.Id,
            Stock: newStock,
        }
        items = append(items, newItem)
    }
    return items
}
