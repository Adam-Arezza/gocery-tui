package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Adam-Arezza/gocery-tui/config"
	"github.com/Adam-Arezza/gocery-tui/internal/components"
	"github.com/Adam-Arezza/gocery-tui/internal/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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

type GroceryStore struct {
    List     list.Model
    Focused  bool
    Height   int
    Width    int
    Debug    string
    ShowQuantityModal bool
    SelectedItem *GroceryItem
    ItemModal components.GroceryModal
    GroceryServer config.ServerConfig
}

type NewCartItemMsg struct{
    item components.CartItem
}

func NewGroceryStore(cfg config.ServerConfig) *GroceryStore {
    delegate := styles.GroceryDelegate{
        StyleDelegate: list.NewDefaultDelegate(),
    }
    delegate.StyleDelegate.Styles.NormalTitle = styles.GroceryItemNormalTitle
    delegate.StyleDelegate.Styles.NormalDesc = styles.GroceryItemNormalDesc
    delegate.StyleDelegate.Styles.SelectedTitle = styles.GroceryItemSelectedTitle
    delegate.StyleDelegate.Styles.SelectedDesc = styles.GroceryItemSelectedDesc

    groceryList := list.New([]list.Item{}, delegate, 0, 0)
    groceryList.Styles.Title = styles.GroceryListTitle
    groceryList.Title = groceryList.Styles.Title.Render("GROCERY STORE")

    //Pagination style
    groceryList.Styles.PaginationStyle = styles.PaginatorStyle
    groceryList.Styles.InactivePaginationDot = styles.PaginatorStyle
    groceryList.Paginator.InactiveDot = "-" 
    groceryList.Paginator.ActiveDot = "*"

    //Help and HelpStyle style
    groceryList.Styles.HelpStyle = styles.HelpStyle
    groceryList.Help.Styles.ShortDesc = styles.HelpMenuStyle
    groceryList.Help.Styles.ShortKey = styles.HelpMenuStyle 

    return &GroceryStore{
        List: groceryList,
        GroceryServer: cfg,
    }
}

func (g *GroceryStore) Init() tea.Cmd {
    return g.LoadItems()
}

type GroceryApiResponse []GroceryItem

func (g *GroceryStore) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {

    case tea.WindowSizeMsg:
        g.Height = msg.Height
        g.Width = msg.Width / 2
        g.List.SetSize(g.Width, g.Height-2) // Leave room for title
        return g, cmd

    case GroceryApiResponse:
        var items []list.Item
        for _, item := range msg {
            items = append(items, item)
        }
        g.List.SetItems(items)
        return g, cmd

    case tea.KeyMsg:

        switch msg.String(){

        case "enter", " ":
            if !g.ShowQuantityModal && g.Focused{
                item, ok := g.List.SelectedItem().(GroceryItem) 
                if ok{
                    g.ItemModal = components.GroceryModal{
                        GroceryItem: components.CartItem{
                            Id: item.Id,
                            Name: item.Title(),
                            Price: item.Price,
                            Stock: item.Stock,
                            Quantity: 1,
                        },
                        Height: g.Height/2,
                        Width: g.Width/2,
                    }
                    g.ShowQuantityModal = true
                    return g,cmd
                }else{
                    g.Debug += fmt.Sprintf("Error: %v\n", msg)
                    return g,cmd
                }
            }else{
                //get the item and add it to the cart
                if g.Focused{
                    newCartItem := g.ItemModal.GroceryItem
                    msg := func() (tea.Msg){
                        return NewCartItemMsg{
                            item: newCartItem,
                        }
                    }
                    g.ShowQuantityModal = false
                    return g, msg
                }
            }

        case "esc":
            if g.ShowQuantityModal{
                g.ShowQuantityModal = false
                return g, cmd
            } 
        }
    }

    if g.Focused && !g.ShowQuantityModal{
        g.List, cmd = g.List.Update(msg)
        return g,cmd
    }

    if g.ShowQuantityModal{
        modal, cmd := g.ItemModal.Update(msg)
        if modal, ok := modal.(*components.GroceryModal); ok {
            g.ItemModal = *modal
        }
        return g,cmd
    }
    return g, cmd
}

func (g *GroceryStore) View() string {
    listView := g.List.View()

    if g.ShowQuantityModal {
        return g.ItemModal.View()
    }

    if g.Focused {
        return styles.FocusedStyle.Height(g.Height).Width(g.Width).Render(listView)
    } else {
        return styles.UnFocusedStyle.Height(g.Height).Width(g.Width).Render(listView)
    }
}

func (g *GroceryStore) LoadItems() tea.Cmd {
    return func() tea.Msg {
        url := "http://" + g.GroceryServer.Host + ":" + g.GroceryServer.Port + "/grocery_items"        
        resp, err := http.Get(url)
        if err != nil {
            return err
        }

        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            return fmt.Errorf("API error: %d", resp.StatusCode)
        }

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            return err
        }

        var groceryItems []GroceryItem
        err = json.Unmarshal(body, &groceryItems)
        if err != nil {
            return err
        }

        return GroceryApiResponse(groceryItems)
    }
}

