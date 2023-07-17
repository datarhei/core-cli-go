package footer

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/datarhei/core-cli-go/ui/messages"
)

var baseStyle = lipgloss.NewStyle()

type Model struct {
	ready bool

	time  time.Time
	width int

	degraded    bool
	degradedErr string

	version string
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		ticker(),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	//var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tickMsg:
		m.time = time.Time(msg)
		cmds = append(cmds, ticker())
	case tea.WindowSizeMsg:
		m.ready = true
		m.width = msg.Width
	case messages.AboutMsg:
		m.degraded = msg.Degraded
		m.degradedErr = msg.DegradedErr
		m.version = msg.Version
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return ""
	}

	t := m.time.Format(time.DateTime)

	style := baseStyle.Copy().Width(m.width)

	b := strings.Builder{}
	if !m.degraded {
		b.WriteString("🍻 Operational, version: " + m.version)
		style = style.Background(lipgloss.Color("#008"))
	} else {
		b.WriteString("🙈 ")
		err := m.degradedErr
		if len(err) > m.width-b.Len()-len(t) {
			err = err[:m.width-b.Len()-len(t)-10]
		}
		b.WriteString(err)
		style = style.Background(lipgloss.Color("#800"))
	}
	b.WriteString(strings.Repeat(" ", m.width-b.Len()-len(t)))
	b.WriteString(t)

	return style.Render(b.String())
}

func New() Model {
	return Model{
		time: time.Now(),
	}
}

type tickMsg time.Time

func ticker() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
