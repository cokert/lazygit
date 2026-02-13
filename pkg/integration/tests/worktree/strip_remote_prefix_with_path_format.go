package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StripRemotePrefixWithPathFormat = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that the remote prefix is stripped before applying createPathFormat to the worktree path",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Worktree.CreatePathFormat = "replace"
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("feature/branch1")
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
					Type("origin/feature/branch1").
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New worktree path")).
					InitialText(Equals("feature-branch1")).
					Clear().
					Type("../linked-worktree").
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New branch name (leave blank to checkout origin/feature/branch1)")).
					InitialText(Equals("feature/branch1")).
					Confirm()
			}).
			Lines(
				Contains("linked-worktree").IsSelected(),
				Contains("repo (main)"),
			)
	},
})
