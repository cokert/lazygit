package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// FormatTerminalTitle resolves the gui.terminalTitle template. Control
// characters are removed from the result, because the terminal would
// interpret them instead of showing them.
func FormatTerminalTitle(terminalTitle string, data config.TerminalTitleData) (string, error) {
	title, err := utils.ResolveTemplate(terminalTitle, data, nil)
	if err != nil {
		return "", err
	}

	return utils.SanitizeTerminalTitle(title), nil
}
