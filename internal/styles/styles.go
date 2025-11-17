package styles

import(
    "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/bubbles/list"
    "strings"
    "fmt"
    "io"
    tea "github.com/charmbracelet/bubbletea"
)


type GroceryDelegate struct {
    StyleDelegate list.DefaultDelegate
}

func (gd GroceryDelegate) Height() int{return 2}
func (gd GroceryDelegate) Spacing() int{return 1}
func (gd GroceryDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {return nil}

func (d GroceryDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var b strings.Builder
	d.StyleDelegate.Render(&b, m, index, listItem)
	row := b.String()
	width := m.Width()
	styled := lipgloss.
		NewStyle().
        Width(width).
		Render(row)
	fmt.Fprint(w, styled)
}

var (
    backgroundColor = lipgloss.Color("62")
    primaryText = lipgloss.Color("230")
    borderColor = lipgloss.Color("#7aeb98")    
    secondaryText = lipgloss.Color("244")
    titleColor = lipgloss.Color("165")

    HeaderStyle = lipgloss.NewStyle(). 
        Foreground(lipgloss.Color("#abafb8")). 
        Bold(true)

    UnFocusedStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()). 
        Padding(0,4)

    FocusedStyle = lipgloss.NewStyle().
	    Border(lipgloss.DoubleBorder()).
        BorderForeground(borderColor).
        Padding(0,4)

    GroceryListTitle = lipgloss.NewStyle().
        Foreground(primaryText).
        Bold(true)
    
    GroceryListPagination = lipgloss.NewStyle().
        Foreground(secondaryText).
        Padding(1,0)
    
    GroceryItemNormalTitle = lipgloss.NewStyle().
        Foreground(primaryText)
    
    GroceryItemNormalDesc = lipgloss.NewStyle().
        Foreground(secondaryText)
    
    GroceryItemSelectedTitle = lipgloss.NewStyle().
        Foreground(titleColor).
        PaddingLeft(1).
        Bold(true)
    
    GroceryItemSelectedDesc = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        PaddingLeft(1).
        Bold(true)

    GroceryTotal = lipgloss.NewStyle(). 
        Foreground(primaryText). 
        Background(backgroundColor). 
        Bold(true)

    ModalStyle = lipgloss.NewStyle(). 
        Border(lipgloss.RoundedBorder()). 
        BorderForeground(lipgloss.Color("#7aeb98")).
        Padding(3)

    ActiveCartItemStyle = lipgloss.NewStyle(). 
        Foreground(primaryText).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(borderColor).
        PaddingLeft(1)

    FocusedPanelTitle = lipgloss.NewStyle(). 
        Background(backgroundColor).
        Foreground(primaryText).
        Padding(0,5).
        Bold(true)

    PaginatorStyle = lipgloss.NewStyle(). 
        Foreground(primaryText)

    HelpStyle = lipgloss.NewStyle(). 
        Foreground(borderColor).
        Border(lipgloss.ASCIIBorder()).
        BorderForeground(borderColor)

    HelpMenuStyle = lipgloss.NewStyle(). 
        Foreground(secondaryText)
)

