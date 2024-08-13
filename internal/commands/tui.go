package commands

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/superboomer/mtiled/internal/options"
)

func TUI(opts *options.Opts) error {
	service, err := createService(opts)
	if err != nil {
		return fmt.Errorf("cant create service: %w", err)
	}

	if err := render(service); err != nil {
		return fmt.Errorf("render error: %w", err)
	}

	return nil
}

func render(service *service) error {

	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	m := model{
		service:  service,
		spinner:  s,
		progress: p,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		return fmt.Errorf("error occurred when creating new program error: %w", err)
	}

	return nil
}
