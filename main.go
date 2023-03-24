package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/campbel/yoshi"
	"github.com/charmbracelet/lipgloss"
)

type Options struct {
	File string `yoshi:"FILE;Markdown file to view;"`
}

func main() {
	yoshi.New("md").Run(func(opts Options) error {
		data, err := os.ReadFile(opts.File)
		if err != nil {
			return err
		}
		fmt.Println(transform(string(data)))
		return nil
	})
}

func transform(data string) string {
	var out []string
	lines := strings.Split(data, "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		switch {
		case hasStringPrefix(line, "```"):
			line, i = transformCodeBlock(i, lines)
		case hasStringPrefix(line, "-"):
			line = transformList(line)
		case strings.HasPrefix(line, "###"):
			line = transformHHH(line)
		case strings.HasPrefix(line, "##"):
			line = transformHH(line)
		case strings.HasPrefix(line, "#"):
			line = transformH(line)
		default:
			line = lipgloss.NewStyle().Width(80).Render(line)
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

var (
	hStyle = lipgloss.NewStyle().
		Width(80).
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("63")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("7D56F4")).
		BorderTop(true).
		BorderBottom(true)

	hhStyle = lipgloss.NewStyle().
		Width(80).
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4"))

	hhhStyle = lipgloss.NewStyle().
			Width(80).
			Bold(false).
			Foreground(lipgloss.Color("#FAFAFA")).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("7D56F4")).
			BorderBottom(true)
)

func hasStringPrefix(s, r string) bool {
	return strings.HasPrefix(strings.TrimSpace(s), r)
}

func transformH(data string) string {
	return hStyle.Render(data)
}

func transformHH(data string) string {
	return hhStyle.Render(data)
}
func transformHHH(data string) string {
	return hhhStyle.Render(data)
}

var (
	bullet = "•"
	hollow = "◦"
	square = "▪"
)

func transformList(data string) string {
	// 2 spaces = bullet
	if strings.HasPrefix(data, "-") {
		return strings.Replace(data, "-", bullet, 1)
	}
	// 4 spaces = hollow
	if strings.HasPrefix(data, "  -") {
		return strings.Replace(data, "-", hollow, 1)
	}
	return strings.Replace(data, "-", square, 1)
}

var (
	codeStyle = lipgloss.NewStyle().
		Padding(1).Width(80).Background(lipgloss.Color("#1E1E1E"))
)

func transformCodeBlock(start int, lines []string) (string, int) {
	var out []string
	language := strings.Replace(strings.TrimSpace(lines[start]), "```", "", 1)
	i := start + 1
	for ; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "```" {
			break
		}
		out = append(out, line)
	}
	var buff bytes.Buffer
	err := quick.Highlight(&buff, strings.Join(out, "\n"), language, "terminal16m", "monokai")
	if err != nil {
		panic(err)
	}
	return codeStyle.Render(buff.String()), i
}
