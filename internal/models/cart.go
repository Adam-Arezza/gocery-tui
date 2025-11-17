package models

import (
	"fmt"
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
}

func (cart *CartModel) Init() tea.Cmd {
return nil
}

func (cart *CartModel) Update(msg tea.Msg) (tea.Model, tea.Cmd){
    var cmd tea.Cmd
    switch msg := msg.(type){

        case tea.WindowSizeMsg:
            cart.Width = msg.Width / 2
            cart.Height = msg.Height
            cart.Items.SetSize(cart.Width, cart.Height-2)

        case NewCartItemMsg:
            var item components.CartItem
            item = msg.item
            cartItems := cart.Items.Items() 
            cartItems = append(cartItems, item)
            cart.Items.SetItems(cartItems)
            cart.Total = cart.CartSum()
            return cart, cmd

        case tea.KeyMsg:
            switch msg.String(){

                case "enter":
                if cart.Focused{
                    cart.Selected = !cart.Selected
                    return cart, cmd
                }

                case "up", "k":
                    if cart.Focused && cart.Selected{
                        item := cart.Items.SelectedItem().(components.CartItem)
                        idx := cart.Items.Index()
                        if item.Quantity >= 1 && item.Quantity < item.Stock{
                            item.Quantity++
                            cart.Items.Items()[idx] = item
                            cart.Items.SetItems(cart.Items.Items())
                            cart.Total = cart.CartSum()
                            return cart, cmd
                        }
                    }

                case "down","j":
                    if cart.Focused && cart.Selected{
                        item := cart.Items.SelectedItem().(components.CartItem)
                        idx := cart.Items.Index()
                        if item.Quantity >= 1{
                            item.Quantity--
                            cart.Items.Items()[idx] = item
                            cart.Items.SetItems(cart.Items.Items())
                            cart.Total = cart.CartSum()
                            return cart, cmd
                        }
                    }

                case "d":
                    if cart.Focused && cart.Selected{
                        idx := cart.Items.Index()
                        cart.Items.RemoveItem(idx)
                        cart.Selected = false
                        cart.Total = cart.CartSum()
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
    if cart.Focused{
        return styles.FocusedStyle.Height(cart.Height).Width(cart.Width).Render(list, cartTotal)
    }else{
        return styles.UnFocusedStyle.Height(cart.Height).Width(cart.Width).Render(list, cartTotal)
    }
}

func NewGroceryCart() *CartModel {
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
    }
}

func (cart *CartModel) CartSum() float32 {
    items := cart.Items.Items()
    var sum float32
    for _, item := range items{
        if cartItem, ok := item.(components.CartItem); ok {
            sum += cartItem.Price * float32(cartItem.Quantity)
        }
    }

    return sum
}
