package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/datarhei/core-cli-go/ui/components/footer"
	"github.com/datarhei/core-cli-go/ui/components/nodes"
	"github.com/datarhei/core-cli-go/ui/components/processes"
	"github.com/datarhei/core-cli-go/ui/messages"
	coreclient "github.com/datarhei/core-client-go/v16"
)

type model struct {
	client coreclient.RestClient

	id       string
	name     string
	numNodes int

	processes processes.Model
	nodes     nodes.Model

	footer footer.Model

	content string

	ready bool

	width  int
	height int
}

func newModel(client coreclient.RestClient) (*model, error) {
	m := &model{
		client:    client,
		nodes:     nodes.New(client),
		processes: processes.New(client),
		footer:    footer.New(),
		content:   "nodes",
	}

	return m, nil
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.footer.Init(),
		m.nodes.Init(),
		m.processes.Init(),
		ticker(),
		m.About(),
		tea.EnterAltScreen,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.nodes, cmd = m.nodes.Update(msg)
	cmds = append(cmds, cmd)

	m.processes, cmd = m.processes.Update(msg)
	cmds = append(cmds, cmd)

	m.footer, cmd = m.footer.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			cmds = append(cmds, tea.Quit)
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "q":
				cmds = append(cmds, tea.Quit)
			case "n":
				m.content = "nodes"
			case "p":
				m.content = "processes"
			}
		}
	case messages.AboutMsg:
		m.id = msg.ID
		m.name = msg.Name
		m.numNodes = len(msg.Nodes)

		cmds = append(cmds, m.About())
	case tea.WindowSizeMsg:
		m.ready = true
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, tea.Batch(cmds...)
}

func height(s string) int {
	return len(strings.Split(s, "\n"))
}

func (m model) View() string {
	if !m.ready {
		return "Initializing ..."
	}

	header := lipgloss.PlaceVertical(2, lipgloss.Top, fmt.Sprintf("datarheiCore Cluster\n%d nodes, connected to %s (%s)", m.numNodes, m.id, m.name))

	content := "Unknown page selected"

	if m.content == "nodes" {
		content = lipgloss.PlaceVertical(m.height-3, lipgloss.Top, m.nodes.View())
	} else if m.content == "processes" {
		content = lipgloss.PlaceVertical(m.height-3, lipgloss.Top, m.processes.View())
	}

	footer := lipgloss.PlaceVertical(1, lipgloss.Bottom, m.footer.View())

	filler := strings.Repeat("\n", m.height-3-height(content))

	return lipgloss.JoinVertical(lipgloss.Left, header, content+filler, footer)
}

func (m model) About() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		about, err := m.client.Cluster()
		if err != nil {
			return messages.ErrorMsg(err)
		}

		a := messages.AboutMsg{
			Nodes:       []messages.AboutNode{},
			ID:          about.ID,
			Name:        about.Name,
			Version:     about.Version,
			Degraded:    about.Degraded,
			DegradedErr: about.DegradedErr,
		}

		for _, node := range about.Nodes {
			n := messages.AboutNode{
				ID:          node.ID,
				Name:        node.Name,
				Uptime:      time.Duration(node.Uptime) * time.Second,
				LastContact: time.Duration(node.LastContact) * time.Millisecond,
				CPUUsage:    node.Resources.CPU / node.Resources.CPULimit,
				MemoryUsage: float64(node.Resources.Mem) / float64(node.Resources.MemLimit),
				Leader:      node.Leader,
				Version:     node.Version,
			}

			a.Nodes = append(a.Nodes, n)
		}

		return a
	})
}

func Run(client coreclient.RestClient) error {
	m, err := newModel(client)
	if err != nil {
		return err
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err = p.Run()

	return err
}

type tickMsg time.Time

func ticker() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
