package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StripRemotePrefix = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that the remote prefix is stripped from the worktree path and branch name prompts",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.CloneIntoRemote("origin")
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
					Clear().
					Type("origin/mybranch").
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New worktree path")).
					InitialText(Equals("mybranch")).
					Clear().
					Type("../linked-worktree").
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New branch name (leave blank to checkout origin/mybranch)")).
					InitialText(Equals("mybranch")).
					Confirm()
			}).
			Lines(
				Contains("linked-worktree").IsSelected(),
				Contains("repo (main)"),
			)
	},
})
