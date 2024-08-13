package commands

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/superboomer/mtiled/internal/downloader"
)

var (
	currentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle           = lipgloss.NewStyle().Margin(1, 2)
	checkMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	errorMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("#DB4035")).SetString("✖")
)

type model struct {
	service       *service
	index         int
	providerIndex int
	countSuccess  int
	count         int
	width         int
	height        int
	spinner       spinner.Model
	progress      progress.Model
	done          bool
}

// Init starts loading files
func (m model) Init() tea.Cmd {
	return tea.Batch(m.downloadAndSave(), m.spinner.Tick)
}

// downloadTileMsg msg for Update
type downloadedTileMsg struct {
	name string
	err  error
}

// downloadAndSave download and save file/ send Msg for Update
func (m *model) downloadAndSave() tea.Cmd {
	point := m.service.points[m.index]
	provider := m.service.providers[m.providerIndex]

	dErr := m.service.downloader.Download(
		&downloader.DownloadRequest{
			Provider: provider,
			Zoom:     m.service.zoom,
			Side:     m.service.side,
			Point:    &point,
		})

	return func() tea.Msg {
		return downloadedTileMsg{name: fmt.Sprintf("%s [%s] [%s]", point.Name, point.ID, provider), err: dErr}
	}
}

// Update update cli by specified *Msg
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	case downloadedTileMsg:
		m.count++

		printMsg := tea.Printf("%s %s", checkMark, msg.name)
		if msg.err != nil {
			printMsg = tea.Printf("%s %s", errorMark, fmt.Sprintf("%s \n%s", msg.name, msg.err.Error()))
		} else {
			m.countSuccess++
		}

		if m.index >= len(m.service.points)-1 && m.providerIndex >= len(m.service.providers)-1 {
			m.done = true
			return m, tea.Sequence(
				printMsg,
				tea.Quit, // exit the program
			)
		}

		if m.providerIndex >= len(m.service.providers)-1 {
			m.providerIndex = 0
			m.index++
		} else {
			m.providerIndex++
		}

		// Update progress bar
		progressCmd := m.progress.SetPercent(float64(m.count) / float64(len(m.service.points)*len(m.service.providers)))

		if msg.err != nil {
			return m, tea.Batch(
				progressCmd,
				printMsg,
				m.downloadAndSave(),
			)
		}

		return m, tea.Batch(
			progressCmd,
			printMsg,
			m.downloadAndSave(), // download the next tile
		)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd
	}

	return m, nil
}

// View draw cli
func (m model) View() string {
	n := len(m.service.points) * len(m.service.providers)
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	if m.done {
		return doneStyle.Render(fmt.Sprintf("Done!\nDownloaded %d tiles\nFailed: %d\n", m.countSuccess, n-m.countSuccess))
	}

	pkgCount := fmt.Sprintf(" %*d/%*d|%*d", w, m.count, w, n, w, m.count-m.countSuccess)

	spin := m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+pkgCount))

	text := currentPkgNameStyle.Render(fmt.Sprintf("%s [%s] [%s]", m.service.points[m.index].Name, m.service.points[m.index].ID, m.service.providers[m.providerIndex]))
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Downloading... ", text)

	cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+pkgCount))
	gap := strings.Repeat(" ", cellsRemaining)

	return spin + info + gap + prog + pkgCount
}
