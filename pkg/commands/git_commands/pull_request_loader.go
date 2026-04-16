package git_commands

import (
	"encoding/json"
	"os/exec"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"golang.org/x/sync/errgroup"
)

type PullRequestInfo struct {
	HeadRefName string `json:"headRefName"`
	Number      int    `json:"number"`
	URL         string `json:"url"`
	State       string `json:"state"`
}

type PullRequestLoader struct {
	*common.Common
	cmd oscommands.ICmdObjBuilder
}

func NewPullRequestLoader(
	cmn *common.Common,
	cmd oscommands.ICmdObjBuilder,
) *PullRequestLoader {
	return &PullRequestLoader{
		Common: cmn,
		cmd:    cmd,
	}
}

func (self *PullRequestLoader) ghIsInstalled() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

func (self *PullRequestLoader) fetchOpenPRs() ([]PullRequestInfo, error) {
	cmdObj := self.cmd.New([]string{
		"gh", "pr", "list",
		"--state", "open",
		"--limit", "100",
		"--json", "headRefName,number,url,state",
	}).DontLog()

	output, err := cmdObj.RunWithOutput()
	if err != nil {
		return nil, err
	}

	var prs []PullRequestInfo
	if err := json.Unmarshal([]byte(output), &prs); err != nil {
		return nil, err
	}

	return prs, nil
}

func (self *PullRequestLoader) fetchPRForBranch(branchName string) (*PullRequestInfo, error) {
	cmdObj := self.cmd.New([]string{
		"gh", "pr", "view", branchName,
		"--json", "headRefName,number,url,state",
	}).DontLog()

	output, err := cmdObj.RunWithOutput()
	if err != nil {
		if strings.Contains(err.Error(), "no pull requests found") ||
			strings.Contains(output, "no pull requests found") {
			return nil, nil
		}
		return nil, err
	}

	var pr PullRequestInfo
	if err := json.Unmarshal([]byte(output), &pr); err != nil {
		return nil, err
	}

	return &pr, nil
}

func (self *PullRequestLoader) SetPullRequestInfoOnBranches(
	branches []*models.Branch,
	renderFunc func(),
) error {
	if !self.ghIsInstalled() {
		return nil
	}

	t := time.Now()

	openPRs, err := self.fetchOpenPRs()
	if err != nil {
		self.Log.Warnf("gh pr list (open) failed: %v", err)
		return nil
	}

	prByBranch := make(map[string]PullRequestInfo, len(openPRs))
	for _, pr := range openPRs {
		prByBranch[pr.HeadRefName] = pr
	}

	var upstreamGoneBranches []*models.Branch
	for _, branch := range branches {
		if _, found := prByBranch[branch.Name]; found {
			continue
		}
		if branch.UpstreamGone {
			upstreamGoneBranches = append(upstreamGoneBranches, branch)
		}
	}

	if len(upstreamGoneBranches) > 0 {
		errg := errgroup.Group{}
		for _, branch := range upstreamGoneBranches {
			errg.Go(func() error {
				pr, err := self.fetchPRForBranch(branch.Name)
				if err != nil {
					self.Log.Warnf("gh pr view for %s failed: %v", branch.Name, err)
					return nil
				}
				if pr != nil {
					prByBranch[pr.HeadRefName] = *pr
				}
				return nil
			})
		}
		_ = errg.Wait()
	}

	for _, branch := range branches {
		if pr, ok := prByBranch[branch.Name]; ok {
			branch.PullRequestNumber.Store(int32(pr.Number))
			branch.PullRequestState.Store(pr.State)
			branch.PullRequestURL.Store(pr.URL)
		}
	}

	self.Log.Debugf("time to load pull request info: %s", time.Since(t))
	renderFunc()
	return nil
}
