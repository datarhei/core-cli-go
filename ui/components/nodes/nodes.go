package nodes

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/datarhei/core-cli-go/ui/messages"
	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/evertras/bubble-table/table"
)

type Model struct {
	ready bool

	client coreclient.RestClient

	nodes []messages.AboutNode
	table table.Model

	width int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ready = true
		m.width = msg.Width
		m.table = m.table.WithTargetWidth(msg.Width)
	case messages.AboutMsg:
		m.nodes = msg.Nodes
		m.updateTable()
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "Initializing ..."
	}

	return m.table.View()
}

const (
	columnKeyID          = "id"
	columnKeyName        = "name"
	columnKeyUptime      = "uptime"
	columnKeyLastContact = "lastcontact"
	columnKeyStatus      = "status"
	columnKeyCPU         = "cpu"
	columnKeyMemory      = "memory"
	columnKeyVersion     = "version"
)

func (m *Model) updateTable() {
	rows := []table.Row{}

	for _, n := range m.nodes {
		pcpu := progress.New(progress.WithWidth(20), progress.WithGradient("#00FF00", "#FF0000"))
		pmem := progress.New(progress.WithWidth(20), progress.WithGradient("#00FF00", "#FF0000"))

		status := "follower"
		colorStr := "#fa0"
		if n.Leader {
			status = "leader"
			colorStr = "#8b8"
		}

		rows = append(rows, table.NewRow(table.RowData{
			columnKeyID:          n.ID,
			columnKeyName:        n.Name,
			columnKeyVersion:     n.Version,
			columnKeyUptime:      n.Uptime.String(),
			columnKeyLastContact: n.LastContact,
			columnKeyStatus:      table.NewStyledCell(status, lipgloss.NewStyle().Foreground(lipgloss.Color(colorStr))),
			columnKeyCPU:         pcpu.ViewAs(n.CPUUsage),
			columnKeyMemory:      pmem.ViewAs(n.MemoryUsage),
		}))
	}

	m.table = m.table.WithRows(rows).SortByAsc(columnKeyID)
}

func New(client coreclient.RestClient) Model {
	m := Model{
		client: client,
	}

	m.table = table.New([]table.Column{
		table.NewFlexColumn(columnKeyID, "ID", 1),
		table.NewFlexColumn(columnKeyName, "Name", 1),
		table.NewFlexColumn(columnKeyVersion, "Version", 1),
		table.NewFlexColumn(columnKeyUptime, "Uptime", 1),
		table.NewFlexColumn(columnKeyLastContact, "Last Contact", 1),
		table.NewFlexColumn(columnKeyStatus, "Status", 1),
		table.NewFlexColumn(columnKeyCPU, "CPU", 1),
		table.NewFlexColumn(columnKeyMemory, "Memory", 1),
	}).WithRows([]table.Row{}).
		Focused(true).
		WithBaseStyle(lipgloss.NewStyle().
			//Foreground(lipgloss.Color("#a7a")).
			BorderForeground(lipgloss.Color("#555")).
			Align(lipgloss.Right)).
		HeaderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)).
		WithPageSize(5)

	return m
}
