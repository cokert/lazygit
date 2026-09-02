package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NavigateCommitFiles = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Navigates to adjacent commits' files from within the commit files panel, refusing to land on a merge commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file.txt", "one\n")
		shell.Commit("one")

		shell.NewBranch("feature")
		shell.CreateFileAndAdd("feature.txt", "feature\n")
		shell.Commit("feature work")

		shell.Checkout("master")
		shell.CreateFileAndAdd("master.txt", "master\n")
		shell.Commit("master work")

		shell.Merge("feature")

		shell.CreateFileAndAdd("after.txt", "after\n")
		shell.Commit("after merge")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("after merge").IsSelected(),
				Contains("Merge branch"),
				Contains("feature work"),
				Contains("master work"),
				Contains("one"),
			).
			NavigateToLine(Contains("feature work")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Title(Contains("feature work")).
			Lines(
				Equals("A feature.txt"),
			)

		// The next (newer) commit from here is the merge commit: refused.
		t.Views().CommitFiles().Press(keys.CommitFiles.NextCommit)
		t.ExpectToast(Equals("Disabled: Can't navigate to a merge commit's files this way, because it has more than one parent"))
		t.Views().CommitFiles().
			Title(Contains("feature work")).
			Lines(
				Equals("A feature.txt"),
			)

		// The previous (older) commit from here is "master work", which
		// isn't a merge, so navigating onto it is fine.
		t.Views().CommitFiles().Press(keys.CommitFiles.PrevCommit)
		t.Views().CommitFiles().
			Title(Contains("master work")).
			Lines(
				Equals("A master.txt"),
			)

		// And back again.
		t.Views().CommitFiles().Press(keys.CommitFiles.NextCommit)
		t.Views().CommitFiles().
			Title(Contains("feature work")).
			Lines(
				Equals("A feature.txt"),
			)
	},
})
