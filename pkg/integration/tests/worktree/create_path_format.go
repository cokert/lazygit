package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CreatePathFormat = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that the createPathFormat config transforms the branch name in the worktree path prompt",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Worktree.CreatePathPrefix = "../"
		cfg.GetUserConfig().Git.Worktree.CreatePathFormat = "replace"
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("feature/my-thing")
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
					InitialText(Equals("feature/my-thing")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New worktree path")).
					InitialText(Equals("../feature-my-thing")).
					Clear().
					Type("../my-worktree").
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New branch name (leave blank to checkout feature/my-thing)")).
					Type("newbranch").
					Confirm()
			}).
			Lines(
				Contains("my-worktree").IsSelected(),
				Contains("repo (main)"),
			)
	},
})
