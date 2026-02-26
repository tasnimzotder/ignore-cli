package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type BrowserModel struct {
	list    list.Model
	loading bool
	spinner spinner.Model
	err     error
	fetchFn func() ([]string, error)
}

func NewBrowser(fetchFn func() ([]string, error)) BrowserModel {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return BrowserModel{
		loading: true,
		spinner: s,
		fetchFn: fetchFn,
	}
}

func (m BrowserModel) Init() tea.Cmd {
	fn := m.fetchFn
	return tea.Batch(m.spinner.Tick, func() tea.Msg {
		names, err := fn()
		if err != nil {
			return templateListErrMsg{err: err}
		}
		return templateListMsg{names: names}
	})
}

func (m BrowserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case templateListMsg:
		items := make([]list.Item, len(msg.names))
		for i, name := range msg.names {
			items[i] = TemplateItem{Name: name}
		}

		delegate := list.NewDefaultDelegate()
		l := list.New(items, delegate, 40, 20)
		l.Title = "Available .gitignore templates"
		l.SetShowStatusBar(true)
		l.SetFilteringEnabled(true)
		l.Styles.Title = titleStyle
		m.list = l
		m.loading = false
		return m, nil

	case templateListErrMsg:
		m.err = msg.err
		m.loading = false
		return m, tea.Quit

	case tea.WindowSizeMsg:
		if !m.loading {
			m.list.SetSize(msg.Width, msg.Height)
		}
		return m, nil

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}
		if m.list.FilterState() != list.Filtering {
			switch msg.String() {
			case "q", "esc":
				return m, tea.Quit
			}
		}

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	if !m.loading {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m BrowserModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	if m.loading {
		return fmt.Sprintf("\n  %s Fetching templates...\n", m.spinner.View())
	}

	return m.list.View()
}
