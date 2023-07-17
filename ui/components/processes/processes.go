package processes

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/datarhei/core-cli-go/ui/messages"
	coreclient "github.com/datarhei/core-client-go/v16"
	coreclientapi "github.com/datarhei/core-client-go/v16/api"
	"github.com/evertras/bubble-table/table"
)

type Model struct {
	ready bool

	client coreclient.RestClient

	processes  []coreclientapi.Process
	processMap coreclientapi.ClusterProcessMap
	table      table.Model

	width int
}

type processMsg struct {
	processes []coreclientapi.Process
	wishMap   coreclientapi.ClusterProcessMap
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.Processes(),
	)
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
	case processMsg:
		m.processes = msg.processes
		m.processMap = msg.wishMap
		m.updateTable()
		cmds = append(cmds, m.Processes())
	case messages.ErrorMsg:
		cmds = append(cmds, m.Processes())
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
	columnKeyID        = "id"
	columnKeyDomain    = "domain"
	columnKeyReference = "reference"
	columnKeyOrder     = "order"
	columnKeyState     = "state"
	columnKeyCPU       = "cpu"
	columnKeyMemory    = "memory"
	columnKeyRuntime   = "runtime"
	columnKeyNodeID    = "nodeid"
	columnKeyLastLog   = "lastlog"
)

func (m *Model) updateTable() {
	rows := []table.Row{}

	for _, p := range m.processes {
		runtime := p.State.Runtime
		if p.State.State != "running" {
			runtime = 0

			if p.State.Reconnect > 0 {
				runtime = -p.State.Reconnect
			}
		}

		order := strings.ToUpper(p.State.Order)
		switch order {
		case "START":
			order = lipgloss.NewStyle().Foreground(lipgloss.Color("#0f0")).Render(order)
		case "STOP":
			order = lipgloss.NewStyle().Foreground(lipgloss.Color("#888")).Render(order)
		}

		state := strings.ToUpper(p.State.State)
		switch state {
		case "RUNNING":
			state = lipgloss.NewStyle().Foreground(lipgloss.Color("#0f0")).Render(state)
		case "FINISHED":
			state = lipgloss.NewStyle().Foreground(lipgloss.Color("#888")).Render(state)
		case "FAILED":
			state = lipgloss.NewStyle().Foreground(lipgloss.Color("#f00")).Render(state)
		case "STARTING":
			state = lipgloss.NewStyle().Foreground(lipgloss.Color("#0ff")).Render(state)
		case "FINISHING":
			state = lipgloss.NewStyle().Foreground(lipgloss.Color("#0ff")).Render(state)
		case "KILLED":
			state = lipgloss.NewStyle().Foreground(lipgloss.Color("#800")).Render(state)
		}

		nodeid := m.processMap[coreclient.NewProcessID(p.ID, p.Domain).String()]
		if nodeid != p.CoreID {
			nodeid = "(" + nodeid + ")"

			if len(p.CoreID) != 0 {
				nodeid = p.CoreID + " " + nodeid
			}

			nodeid = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFAC1C")).Render(nodeid)
		}

		cpu := fmt.Sprintf("%.1f%%", p.State.Resources.CPU.Current)
		if p.State.Resources.CPU.IsThrottling {
			cpu = lipgloss.NewStyle().Foreground(lipgloss.Color("#800")).Render("* " + cpu)
		}

		rows = append(rows, table.NewRow(table.RowData{
			columnKeyID:        p.ID,
			columnKeyDomain:    p.Domain,
			columnKeyReference: p.Reference,
			columnKeyOrder:     order,
			columnKeyState:     state,
			columnKeyCPU:       cpu,
			columnKeyMemory:    formatByteCountBinary(p.State.Resources.Memory.Current),
			columnKeyRuntime:   time.Duration(runtime) * time.Second,
			columnKeyNodeID:    nodeid,
		}))
	}

	m.table = m.table.WithRows(rows).SortByAsc(columnKeyID)
}

func formatByteCountBinary(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d  B", b)
	}

	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func (m Model) Processes() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		list, err := m.client.ClusterProcessList(coreclient.ProcessListOptions{
			Filter: []string{"state"},
		})
		if err != nil {
			return messages.ErrorMsg(err)
		}

		wishMap, _ := m.client.ClusterDBProcessMap()

		return processMsg{
			processes: list,
			wishMap:   wishMap,
		}
	})
}

func New(client coreclient.RestClient) Model {
	m := Model{
		client: client,
	}

	m.table = table.New([]table.Column{
		table.NewFlexColumn(columnKeyID, "ID", 0),
		table.NewFlexColumn(columnKeyDomain, "Domain", 0),
		table.NewFlexColumn(columnKeyReference, "Reference", 0),
		table.NewFlexColumn(columnKeyOrder, "Order", 0),
		table.NewFlexColumn(columnKeyState, "State", 0),
		table.NewFlexColumn(columnKeyRuntime, "Runtime", 0),
		table.NewFlexColumn(columnKeyNodeID, "Node", 0),
		table.NewFlexColumn(columnKeyCPU, "CPU", 0),
		table.NewFlexColumn(columnKeyMemory, "Memory", 0),
	}).WithRows([]table.Row{}).
		Focused(true).
		WithBaseStyle(lipgloss.NewStyle().
			//Foreground(lipgloss.Color("#a7a")).
			BorderForeground(lipgloss.Color("#555")).
			Align(lipgloss.Right)).
		HeaderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)).
		WithPageSize(10)

	return m
}
