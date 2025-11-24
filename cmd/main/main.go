package main

import (
	"fmt"
	"github.com/Adam-Arezza/gocery-tui/config"
	"github.com/Adam-Arezza/gocery-tui/internal/models"
    "github.com/Adam-Arezza/gocery-tui/internal/messages"
	"github.com/Adam-Arezza/gocery-tui/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
    groceryStore *models.GroceryStore
    groceryCart *models.CartModel
    Focus bool
    Height int
    Width int
}

func (m *model) Init() tea.Cmd{
    return tea.Batch(
        m.groceryStore.Init(),
        m.groceryCart.Init(),
    )
}

func (m *model) Update(msg tea.Msg)(tea.Model, tea.Cmd){
    var cmds []tea.Cmd
    var cmd tea.Cmd

    switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	    m.Width = msg.Width
	    m.Height = msg.Height
        sizeMsg := tea.WindowSizeMsg{
            Height: msg.Height - 5,
            Width: msg.Width - 3,
        }

        groceryModel, _ := m.groceryStore.Update(sizeMsg)
        if groceryStore, ok := groceryModel.(*models.GroceryStore); ok {
            m.groceryStore = groceryStore
        }

        groceryCart, _ := m.groceryCart.Update(sizeMsg)
        if groceryCart, ok := groceryCart.(*models.CartModel); ok {
            m.groceryCart = groceryCart
        }
        return m, nil

    case messages.CompletePurchaseMsg:
        return m, m.groceryStore.LoadItems()

	case tea.KeyMsg:
        switch msg.String() {
		case "ctrl+c", "q":
            return m, tea.Quit

        case "tab":
            if !m.groceryCart.Selected && !m.groceryStore.ShowQuantityModal{
                m.Focus = !m.Focus
                m.groceryStore.Focused = m.Focus
                m.groceryCart.Focused = !m.Focus

                if m.groceryStore.Focused {
                    m.groceryStore.List.Styles.Title = styles.FocusedPanelTitle
                }else{
                    m.groceryStore.List.Styles.Title = styles.GroceryListTitle
                }

                if m.groceryCart.Focused{
                    m.groceryCart.Items.Styles.Title = styles.FocusedPanelTitle
                }else{
                    m.groceryCart.Items.Styles.Title = styles.GroceryListTitle
                }
            }
        }

    case messages.NewCartItemMsg:
        m.groceryCart.Update(msg)
        return m, cmd
    }

    groceryCart, cmd := m.groceryCart.Update(msg)
        if groceryCart, ok := groceryCart.(*models.CartModel); ok {
            m.groceryCart = groceryCart
        }
    cmds = append(cmds, cmd)

    model, cmd := m.groceryStore.Update(msg)
        if groceryStore, ok := model.(*models.GroceryStore); ok {
            m.groceryStore = groceryStore
        }
    cmds = append(cmds, cmd)
    return m,tea.Batch(cmds...)
}

func (m *model) View()string{
    cart := m.groceryCart.View()
    store := m.groceryStore.View()
    mainView := lipgloss.JoinHorizontal(lipgloss.Top, cart, store)
    return mainView
}

func main(){
    serverConfig, err := config.Load("./config.json")
    if err != nil {
        fmt.Printf(err.Error())
    }
    groceryStore := models.NewGroceryStore(*serverConfig)
    groceryCart := models.NewGroceryCart(*&serverConfig)
    p := tea.NewProgram(&model{groceryStore: groceryStore, groceryCart: groceryCart, Focus: false}, tea.WithAltScreen())
	if _,err := p.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
