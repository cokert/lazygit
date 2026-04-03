package presentation

import (
	"fmt"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
)

func FormatStatus(
	repoName string,
	currentBranch *models.Branch,
	itemOperation types.ItemOperation,
	linkedWorktreeName string,
	runtreeName string,
	workingTreeState models.WorkingTreeState,
	tr *i18n.TranslationSet,
	userConfig *config.UserConfig,
) string {
	status := ""

	if currentBranch.IsRealBranch() {
		status += BranchStatus(currentBranch, itemOperation, tr, time.Now(), userConfig)
		if status != "" {
			status += " "
		}
	}

	if workingTreeState.Any() {
		status += style.FgYellow.Sprintf("(%s) ", workingTreeState.LowerCaseTitle(tr))
	}

	name := GetBranchTextStyle(currentBranch.Name).Sprint(currentBranch.Name)
	// If the user is in a linked worktree (i.e. not the main worktree) we'll display that,
	// but skip it if the worktree name matches the branch name (redundant)
	if linkedWorktreeName != "" && linkedWorktreeName != currentBranch.Name {
		icon := ""
		if icons.IsIconEnabled() {
			icon = icons.LINKED_WORKTREE_ICON + " "
		}
		repoName = fmt.Sprintf("%s(%s%s)", repoName, icon, style.FgCyan.Sprint(linkedWorktreeName))
	}
	status += fmt.Sprintf("%s → %s", repoName, name)

	if runtreeName != "" {
		if runtreeName == currentBranch.Name {
			status += " [🎄]"
		} else {
			status += fmt.Sprintf(" [🎄 %s]", style.FgYellow.Sprint(runtreeName))
		}
	}

	return status
}
