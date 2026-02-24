package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var randomWords = []string{
	"serendipity", "ephemeral", "melancholy", "luminous", "cascade",
	"whisper", "labyrinth", "solitude", "twilight", "enigma",
	"quixotic", "nebula", "reverie", "zephyr", "kaleidoscope",
	"phantasm", "gossamer", "bioluminescent", "transcendent", "ineffable",
	"petrichor", "sonder", "hiraeth", "eudaimonia", "vellichor",
}

type message struct {
	text   string
	isUser bool
}

type model struct {
	messages []message
	input    string
	width    int
}

var (
	userStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true)

	botStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("213")).
			Bold(true)

	userMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	botMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Italic(true)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("213")).
			Bold(true).
			Padding(0, 1)

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("237"))
)

func randomResponse() string {
	count := rand.Intn(4) + 1
	words := make([]string, count)
	for i := range words {
		words[i] = randomWords[rand.Intn(len(randomWords))]
	}
	return strings.Join(words, " ")
}

func initialModel() model {
	return model{
		messages: []message{
			{text: "こんにちは！何でも話しかけてください。", isUser: false},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			if strings.TrimSpace(m.input) == "" {
				return m, nil
			}
			userMsg := message{text: m.input, isUser: true}
			botMsg := message{text: randomResponse(), isUser: false}
			m.messages = append(m.messages, userMsg, botMsg)
			m.input = ""

		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}

		default:
			if msg.Type == tea.KeyRunes {
				m.input += string(msg.Runes)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	return m, nil
}

func (m model) View() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("  TUI Chat  "))
	sb.WriteString("\n")

	divider := dividerStyle.Render(strings.Repeat("─", max(m.width, 20)))
	sb.WriteString(divider)
	sb.WriteString("\n")

	for _, msg := range m.messages {
		if msg.isUser {
			sb.WriteString(userStyle.Render("You: "))
			sb.WriteString(userMsgStyle.Render(msg.text))
		} else {
			sb.WriteString(botStyle.Render("Bot: "))
			sb.WriteString(botMsgStyle.Render(msg.text))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(divider)
	sb.WriteString("\n")
	sb.WriteString(inputStyle.Render("> "))
	sb.WriteString(m.input)
	sb.WriteString("█")
	sb.WriteString("\n")
	sb.WriteString(dividerStyle.Render("Ctrl+C / Esc: 終了"))

	return sb.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
}
