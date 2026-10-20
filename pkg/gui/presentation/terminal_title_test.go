package presentation

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestFormatTerminalTitle(t *testing.T) {
	scenarios := []struct {
		testName      string
		terminalTitle string
		data          config.TerminalTitleData
		expected      string
	}{
		{
			testName:      "Default format",
			terminalTitle: config.GetDefaultConfig().Gui.TerminalTitle,
			data:          config.TerminalTitleData{RepoName: "myrepo", BranchName: "main"},
			expected:      "lazygit [myrepo]",
		},
		{
			testName:      "All fields",
			terminalTitle: "{{.RepoName}}:{{.WorktreeName}}:{{.BranchName}}",
			data:          config.TerminalTitleData{RepoName: "myrepo", BranchName: "main", WorktreeName: "wt"},
			expected:      "myrepo:wt:main",
		},
		{
			testName:      "Worktree name only when in a linked worktree",
			terminalTitle: "{{.RepoName}}{{if .WorktreeName}}({{.WorktreeName}}){{end}}",
			data:          config.TerminalTitleData{RepoName: "myrepo"},
			expected:      "myrepo",
		},
		{
			testName:      "Worktree name in a linked worktree",
			terminalTitle: "{{.RepoName}}{{if .WorktreeName}}({{.WorktreeName}}){{end}}",
			data:          config.TerminalTitleData{RepoName: "myrepo", WorktreeName: "wt"},
			expected:      "myrepo(wt)",
		},
		{
			testName:      "Control characters are removed",
			terminalTitle: "{{.RepoName}}",
			data:          config.TerminalTitleData{RepoName: "my\x1b]0;evil\x07repo"},
			expected:      "my]0;evilrepo",
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			title, err := FormatTerminalTitle(s.terminalTitle, s.data)
			assert.NoError(t, err)
			assert.Equal(t, s.expected, title)
		})
	}

	t.Run("Invalid template", func(t *testing.T) {
		_, err := FormatTerminalTitle("{{.RepoName", config.TerminalTitleData{})
		assert.Error(t, err)
	})
}
