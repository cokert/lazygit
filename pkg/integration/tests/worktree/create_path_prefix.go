package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CreatePathPrefix = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that the createPathPrefix config seeds the worktree path prompt",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Worktree.CreatePathPrefix = "../"
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Worktrees().
			Focus().
			Press(keys.Universal.New).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Worktree")).
					Select(Contains("Create worktree from ref").DoesNotContain("detached")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New worktree base ref")).
					InitialText(Equals("mybranch")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New worktree path")).
					InitialText(Equals("../")).
					Type("linked-worktree").
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New branch name (leave blank to checkout mybranch)")).
					Type("newbranch").
					Confirm()
			}).
			Lines(
				Contains("linked-worktree").IsSelected(),
				Contains("repo (main)"),
			)
	},
})
