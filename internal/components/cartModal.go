package components

import (
	"fmt"

	"github.com/Adam-Arezza/gocery-tui/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


type CartModal struct {
    TotalPrice float32
    Confirm bool
    Height int
    Width int
    Wallet float32
    InsufficientFunds bool
}


func (c *CartModal)Init() tea.Cmd{
    return nil
}

type PurchaseMsg struct{}
type CloseModalMsg struct{}

func (c *CartModal) Update(msg tea.Msg)(tea.Model, tea.Cmd){
    var cmd tea.Cmd
    switch msg := msg.(type){
    case tea.KeyMsg:
        switch msg.String(){
        case " ", "c":
            if c.TotalPrice > c.Wallet{
                c.InsufficientFunds = true
                return c, cmd
            }
            c.Confirm = true
            cmd := func () tea.Msg{return PurchaseMsg{}}
            return c, cmd

        case "esc":
            cmd := func() tea.Msg{return CloseModalMsg{}}
            return c, cmd
        }
    }
    return c,cmd 
}

func (c CartModal) View() string {
    price := styles.HeaderStyle.Render(fmt.Sprintf("Total Price: $%.2f\n", c.TotalPrice))
    wallet := styles.HeaderStyle.Render(fmt.Sprintf("Wallet: $%.2f\n", c.Wallet))
    confirm := styles.HelpStyle.Render(fmt.Sprintf("Confirm with space or 'c' | esc to cancel"))
    modal := styles.ModalStyle.Render(wallet, price, confirm)
    if c.InsufficientFunds{
        fundsMsg := lipgloss.NewStyle(). 
        Foreground(lipgloss.Color("160")).
        Render(fmt.Sprint("Not enough funds in wallet for purchase!"))        
        modal := styles.ModalStyle.Render(wallet, price, confirm + "\n", fundsMsg)
        return modal
    }

    return modal
}

