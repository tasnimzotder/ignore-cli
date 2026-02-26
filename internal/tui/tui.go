package tui

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	statusStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	addedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	normalStyle   = lipgloss.NewStyle()
	previewBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)
	previewTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
)

// TemplateItem represents a template in the list.
type TemplateItem struct {
	Name string
}

func (t TemplateItem) FilterValue() string { return t.Name }
func (t TemplateItem) Title() string       { return t.Name }
func (t TemplateItem) Description() string { return "" }

// unifiedDelegate renders list items with selection and status indicators.
type unifiedDelegate struct {
	selected map[string]struct{}
	added    map[string]struct{}
}

func (d unifiedDelegate) Height() int                             { return 1 }
func (d unifiedDelegate) Spacing() int                            { return 0 }
func (d unifiedDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d unifiedDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(TemplateItem)
	if !ok {
		return
	}

	cursor := "  "
	if index == m.Index() {
		cursor = "> "
	}

	var check string
	var style lipgloss.Style

	_, isSelected := d.selected[item.Name]
	_, isAdded := d.added[item.Name]

	switch {
	case isSelected:
		check = "[x]"
		style = selectedStyle
	case isAdded:
		check = "[+]"
		style = addedStyle
	default:
		check = "[ ]"
		style = normalStyle
	}

	fmt.Fprint(w, style.Render(fmt.Sprintf("%s%s %s", cursor, check, item.Name)))
}

// FetchFn fetches the list of template names.
type FetchFn func() ([]string, error)

// ContentFn fetches the content of a template by name.
type ContentFn func(name string) (string, error)

// Result holds the user's selections after the TUI exits.
type Result struct {
	Selected []string
	Canceled bool
}

// Messages
type templateListMsg struct{ names []string }
type templateListErrMsg struct{ err error }
type templateContentMsg struct {
	name    string
	content string
}
type templateContentErrMsg struct {
	name string
	err  error
}

// Model is the unified TUI model.
type Model struct {
	list     list.Model
	viewport viewport.Model
	selected map[string]struct{}
	added    map[string]struct{}

	loading        bool
	spinner        spinner.Model
	previewLoading bool
	previewName    string
	previewContent map[string]string
	err            error
	done           bool
	override       bool

	fetchFn   FetchFn
	contentFn ContentFn

	width  int
	height int
	ready  bool
}

// New creates a new unified TUI model.
func New(fetchFn FetchFn, contentFn ContentFn, added []string, override bool) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	addedMap := make(map[string]struct{}, len(added))
	for _, name := range added {
		addedMap[name] = struct{}{}
	}

	return Model{
		selected:       make(map[string]struct{}),
		added:          addedMap,
		loading:        true,
		spinner:        s,
		previewContent: make(map[string]string),
		fetchFn:        fetchFn,
		contentFn:      contentFn,
		override:       override,
	}
}

func (m Model) Init() tea.Cmd {
	fn := m.fetchFn
	return tea.Batch(m.spinner.Tick, func() tea.Msg {
		names, err := fn()
		if err != nil {
			return templateListErrMsg{err: err}
		}
		return templateListMsg{names: names}
	})
}

func (m Model) fetchContent(name string) tea.Cmd {
	fn := m.contentFn
	return func() tea.Msg {
		content, err := fn(name)
		if err != nil {
			return templateContentErrMsg{name: name, err: err}
		}
		return templateContentMsg{name: name, content: content}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case templateListMsg:
		items := make([]list.Item, len(msg.names))
		for i, name := range msg.names {
			items[i] = TemplateItem{Name: name}
		}

		delegate := unifiedDelegate{selected: m.selected, added: m.added}
		listWidth := m.listWidth()
		l := list.New(items, delegate, listWidth, m.height-2)
		l.Title = "Select templates (space=toggle, enter=confirm)"
		l.SetShowStatusBar(true)
		l.SetFilteringEnabled(true)
		l.Styles.Title = titleStyle
		m.list = l
		m.loading = false

		vpWidth := m.previewWidth()
		m.viewport = viewport.New(viewport.WithWidth(vpWidth), viewport.WithHeight(m.height-4))
		m.viewport.SetContent("Navigate to a template to preview its contents.")
		m.ready = true

		if len(items) > 0 {
			if item, ok := items[0].(TemplateItem); ok {
				m.previewName = item.Name
				m.previewLoading = true
				return m, m.fetchContent(item.Name)
			}
		}
		return m, nil

	case templateListErrMsg:
		m.err = msg.err
		m.loading = false
		return m, tea.Quit

	case templateContentMsg:
		m.previewContent[msg.name] = msg.content
		if msg.name == m.previewName {
			m.previewLoading = false
			m.viewport.SetContent(msg.content)
			m.viewport.GotoTop()
		}
		return m, nil

	case templateContentErrMsg:
		if msg.name == m.previewName {
			m.previewLoading = false
			m.viewport.SetContent(fmt.Sprintf("Error loading preview: %v", msg.err))
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.ready {
			m.list.SetSize(m.listWidth(), m.height-2)
			m.viewport.SetWidth(m.previewWidth())
			m.viewport.SetHeight(m.height - 4)
		}
		return m, nil

	case tea.KeyPressMsg:
		if m.loading {
			return m, nil
		}

		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "space":
			if item, ok := m.list.SelectedItem().(TemplateItem); ok {
				if _, exists := m.selected[item.Name]; exists {
					delete(m.selected, item.Name)
				} else {
					m.selected[item.Name] = struct{}{}
				}
				m.list.SetDelegate(unifiedDelegate{selected: m.selected, added: m.added})
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

	if m.ready {
		var cmds []tea.Cmd

		var listCmd tea.Cmd
		m.list, listCmd = m.list.Update(msg)
		if listCmd != nil {
			cmds = append(cmds, listCmd)
		}

		if item, ok := m.list.SelectedItem().(TemplateItem); ok && item.Name != m.previewName {
			m.previewName = item.Name
			if content, cached := m.previewContent[item.Name]; cached {
				m.viewport.SetContent(content)
				m.viewport.GotoTop()
				m.previewLoading = false
			} else {
				m.previewLoading = true
				m.viewport.SetContent("Loading...")
				cmds = append(cmds, m.fetchContent(item.Name))
			}
		}

		var vpCmd tea.Cmd
		m.viewport, vpCmd = m.viewport.Update(msg)
		if vpCmd != nil {
			cmds = append(cmds, vpCmd)
		}

		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m Model) View() tea.View {
	var content string

	if m.err != nil {
		content = fmt.Sprintf("Error: %v\n", m.err)
	} else if m.loading {
		content = fmt.Sprintf("\n  %s Fetching templates...\n", m.spinner.View())
	} else {
		leftPane := m.list.View()

		var previewHeader string
		if m.previewLoading {
			previewHeader = previewTitle.Render(fmt.Sprintf("Preview: %s (loading...)", m.previewName))
		} else {
			previewHeader = previewTitle.Render(fmt.Sprintf("Preview: %s", m.previewName))
		}
		rightPane := previewBorder.Width(m.previewWidth()).Render(
			previewHeader + "\n" + m.viewport.View(),
		)

		joined := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

		count := len(m.selected)
		status := ""
		if count > 0 {
			names := make([]string, 0, count)
			for name := range m.selected {
				names = append(names, name)
			}
			status = statusStyle.Render(fmt.Sprintf("  %d selected: %s", count, strings.Join(names, ", ")))
		}

		content = joined + "\n" + status
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func (m Model) Result() Result {
	if m.selected == nil {
		return Result{Canceled: true}
	}
	names := make([]string, 0, len(m.selected))
	for name := range m.selected {
		names = append(names, name)
	}
	return Result{Selected: names}
}

func (m Model) listWidth() int {
	if m.width == 0 {
		return 35
	}
	return m.width * 2 / 5
}

func (m Model) previewWidth() int {
	if m.width == 0 {
		return 45
	}
	return m.width - m.listWidth() - 2
}
