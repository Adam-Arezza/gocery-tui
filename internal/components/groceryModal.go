package components

import (
    "fmt"
	tea "github.com/charmbracelet/bubbletea"
    "github.com/Adam-Arezza/gocery-tui/internal/styles"
    "github.com/Adam-Arezza/gocery-tui/internal/types"
    "github.com/charmbracelet/lipgloss"
)


type GroceryModal struct{
    GroceryItem types.CartItem
    Height int
    Width int
}

func (gm *GroceryModal) Init() tea.Cmd{
    return nil
}

func (gm *GroceryModal) Update(msg tea.Msg) (tea.Model, tea.Cmd){
    var cmd tea.Cmd
    switch msg := msg.(type){
    case tea.KeyMsg:
        switch msg.String(){
        case "up", "k":
            if gm.GroceryItem.Quantity < gm.GroceryItem.Stock{
                gm.GroceryItem.Quantity++
                return gm, cmd
            }

        case "down", "j":
            if gm.GroceryItem.Quantity > 1{
                gm.GroceryItem.Quantity--
            }
            return gm,cmd
        }
    }
    return gm, cmd
}

func (gm *GroceryModal) View() string{
    header := styles.HeaderStyle.Render("Select Quantity")
    item := lipgloss.NewStyle(). 
                Foreground(lipgloss.Color("200")). 
                Render(fmt.Sprintf("Item: %s\n", gm.GroceryItem.Name))
    details := fmt.Sprintf(
        "Price: $%.2f\nStock: %d\n\nQuantity: %d",
        gm.GroceryItem.Price,
        gm.GroceryItem.Stock,
        gm.GroceryItem.Quantity,
    )
    instructions := "↑/k ↓/j to change quantity\nEnter to confirm\nEsc to cancel"
    content := styles.ModalStyle.Render(fmt.Sprintf("%s\n%s\n%s\n\n%s", header, item,  details, instructions))
    modalScreen := lipgloss.Place(
        gm.Width*2,//width
        gm.Height*2,//height
        lipgloss.Center,//horizontalpos
        lipgloss.Center,//vertpos
        content,//content string
    )
    return modalScreen
}

