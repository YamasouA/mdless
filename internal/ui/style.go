package ui

import "github.com/charmbracelet/lipgloss"

var (
	Header = lipgloss.NewStyle().
		Background(lipgloss.Color("238")).
		Padding(0, 1)

	ActiveTab = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("62")).
			Padding(0, 1)

	InactiveTab = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Background(lipgloss.Color("238")).
			Padding(0, 1)

	TabSeparator = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Background(lipgloss.Color("238"))

	Footer = lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("236")).
		Padding(0, 1)
)
