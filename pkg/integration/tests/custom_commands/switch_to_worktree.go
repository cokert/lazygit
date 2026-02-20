package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SwitchToWorktree = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using a custom command to create a worktree and switch to it via the after.switchToWorktree hook",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     "W",
				Context: "localBranches",
				Command: `git worktree add -b {{.Form.Name}} ../{{.Form.Name}} {{.SelectedLocalBranch.Name | quote}}`,
				Prompts: []config.CustomCommandPrompt{
					{
						Key:   "Name",
						Type:  "input",
						Title: "New branch/worktree name",
					},
				},
				After: &config.CustomCommandAfterHook{
					SwitchToWorktree: "../{{.Form.Name}}",
				},
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("mybranch"),
			).
			Press("W")

		t.ExpectPopup().Prompt().
			Title(Equals("New branch/worktree name")).
			Type("newbranch").
			Confirm()

		// After the command completes and the switch happens,
		// verify we're now in the new worktree
		t.Views().Status().
			Lines(
				Contains("repo(newbranch) → newbranch"),
			)

		t.Views().Branches().
			Lines(
				Contains("newbranch"),
				Contains("mybranch (worktree)"),
			)
	},
})
