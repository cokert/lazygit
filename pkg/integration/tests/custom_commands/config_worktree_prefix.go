package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ConfigWorktreePrefix = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Custom command can access Config.Git.Worktree.CreatePathPrefix",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("blah")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Worktree.CreatePathPrefix = "../wt-"
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     "X",
				Context: "files",
				Command: "printf '%s' '{{ .Config.Git.Worktree.CreatePathPrefix }}' > prefix.txt",
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press("X")

		t.FileSystem().FileContent("prefix.txt", Equals("../wt-"))
	},
})
