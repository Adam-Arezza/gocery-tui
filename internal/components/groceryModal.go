package components

import (
    "fmt"
	tea "github.com/charmbracelet/bubbletea"
    "github.com/Adam-Arezza/gocery-tui/internal/styles"
    "github.com/charmbracelet/lipgloss"
)

type CartItem struct {
    Name     string  
    Price    float32 
    Stock    int     
    Quantity int
}

type GroceryModal struct{
    GroceryItem CartItem
}

func (gm *GroceryModal) Init() tea.Cmd{
    return nil
}

func (c CartItem) Title() string       { return c.Name }
func (c CartItem) Description() string { return fmt.Sprintf("Qty: %d | Price: $%.2f", c.Quantity, c.Price) }
func (c CartItem) FilterValue() string { return c.Name }

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
    content := fmt.Sprintf("%s\n%s\n%s\n\n%s", header, item,  details, instructions)
    return styles.ModalStyle.Render(content)
}

