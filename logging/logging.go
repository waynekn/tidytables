package logging

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func OpenLogFile() (*os.File, error) {
	f, err := tea.LogToFile("debug.log", "debug")

	if err != nil {
		return nil, err
	}

	return f, nil
}
