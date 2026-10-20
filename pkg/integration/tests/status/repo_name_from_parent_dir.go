package status

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RepoNameFromParentDir = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Names the repo after the directory around its main worktree if the main worktree's directory is listed in git.mainWorktreeDirNames",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		// The main worktree of an integration test's repo is a directory called
		// "repo", inside the test's "actual" directory.
		config.GetUserConfig().Git.MainWorktreeDirNames = []string{"repo"}
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().Content(Contains("actual → master"))
	},
})
