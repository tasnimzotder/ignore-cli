package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	normalStyle   = lipgloss.NewStyle()
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	statusStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

type TemplateItem struct {
	Name string
}

func (t TemplateItem) FilterValue() string { return t.Name }
func (t TemplateItem) Title() string       { return t.Name }
func (t TemplateItem) Description() string { return "" }

type pickerDelegate struct {
	selected map[string]struct{}
}

func (d pickerDelegate) Height() int                             { return 1 }
func (d pickerDelegate) Spacing() int                            { return 0 }
func (d pickerDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d pickerDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(TemplateItem)
	if !ok {
		return
	}

	cursor := "  "
	if index == m.Index() {
		cursor = "> "
	}

	check := "[ ]"
	style := normalStyle
	if _, ok := d.selected[item.Name]; ok {
		check = "[x]"
		style = selectedStyle
	}

	fmt.Fprint(w, style.Render(fmt.Sprintf("%s%s %s", cursor, check, item.Name)))
}

type PickerModel struct {
	list     list.Model
	selected map[string]struct{}
	loading  bool
	spinner  spinner.Model
	err      error
	done     bool
	fetchFn  func() ([]string, error)
}

type PickerResult struct {
	Selected []string
	Canceled bool
}

type templateListMsg struct {
	names []string
}

type templateListErrMsg struct {
	err error
}

func NewPicker(fetchFn func() ([]string, error)) PickerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return PickerModel{
		selected: make(map[string]struct{}),
		loading:  true,
		spinner:  s,
		fetchFn:  fetchFn,
	}
}

func (m PickerModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchTemplates())
}

func (m PickerModel) fetchTemplates() tea.Cmd {
	fn := m.fetchFn
	return func() tea.Msg {
		names, err := fn()
		if err != nil {
			return templateListErrMsg{err: err}
		}
		return templateListMsg{names: names}
	}
}

func (m PickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case templateListMsg:
		items := make([]list.Item, len(msg.names))
		for i, name := range msg.names {
			items[i] = TemplateItem{Name: name}
		}

		delegate := pickerDelegate{selected: m.selected}
		l := list.New(items, delegate, 40, 20)
		l.Title = "Select templates (space=toggle, enter=confirm)"
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
			m.list.SetSize(msg.Width, msg.Height-2)
		}
		return m, nil

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}

		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case " ":
			if item, ok := m.list.SelectedItem().(TemplateItem); ok {
				if _, exists := m.selected[item.Name]; exists {
					delete(m.selected, item.Name)
				} else {
					m.selected[item.Name] = struct{}{}
				}
				m.list.SetDelegate(pickerDelegate{selected: m.selected})
			}
			return m, nil
		case "enter":
			m.done = true
			return m, tea.Quit
		case "q", "esc":
			m.selected = nil
			m.done = true
			return m, tea.Quit
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

func (m PickerModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	if m.loading {
		return fmt.Sprintf("\n  %s Fetching templates...\n", m.spinner.View())
	}

	count := len(m.selected)
	status := ""
	if count > 0 {
		names := make([]string, 0, count)
		for name := range m.selected {
			names = append(names, name)
		}
		status = statusStyle.Render(fmt.Sprintf("\n  %d selected: %s", count, strings.Join(names, ", ")))
	}

	return m.list.View() + status + "\n"
}

func (m PickerModel) Result() PickerResult {
	if m.selected == nil {
		return PickerResult{Canceled: true}
	}

	names := make([]string, 0, len(m.selected))
	for name := range m.selected {
		names = append(names, name)
	}
	return PickerResult{Selected: names}
}
