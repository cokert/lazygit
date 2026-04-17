package git_commands

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
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

func branchAlias(idx int) string {
	return fmt.Sprintf("pr_%d", idx)
}

func escapeGraphQLString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

func (self *PullRequestLoader) buildGraphQLQuery(branchNames []string) string {
	var sb strings.Builder
	sb.WriteString("query($owner: String!, $repo: String!) {")
	for i, name := range branchNames {
		fmt.Fprintf(&sb,
			` %s: repository(owner: $owner, name: $repo) { pullRequests(headRefName: "%s", first: 1, orderBy: {field: CREATED_AT, direction: DESC}) { nodes { number url state headRefName } } }`,
			branchAlias(i), escapeGraphQLString(name))
	}
	sb.WriteString(" }")
	return sb.String()
}

type graphQLPRNode struct {
	Number      int    `json:"number"`
	URL         string `json:"url"`
	State       string `json:"state"`
	HeadRefName string `json:"headRefName"`
}

type graphQLRepoResult struct {
	PullRequests struct {
		Nodes []graphQLPRNode `json:"nodes"`
	} `json:"pullRequests"`
}

type graphQLResponse struct {
	Data   map[string]graphQLRepoResult `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (self *PullRequestLoader) getRepoOwnerAndName() (string, string, error) {
	cmdObj := self.cmd.New([]string{
		"gh", "repo", "view", "--json", "owner,name",
	}).DontLog()

	output, err := cmdObj.RunWithOutput()
	if err != nil {
		return "", "", err
	}

	var repo struct {
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(output), &repo); err != nil {
		return "", "", err
	}

	return repo.Owner.Login, repo.Name, nil
}

func (self *PullRequestLoader) fetchPRsGraphQL(owner, repoName string, branchNames []string) ([]PullRequestInfo, error) {
	query := self.buildGraphQLQuery(branchNames)

	cmdObj := self.cmd.New([]string{
		"gh", "api", "graphql",
		"-F", fmt.Sprintf("owner=%s", owner),
		"-F", fmt.Sprintf("repo=%s", repoName),
		"-f", fmt.Sprintf("query=%s", query),
	}).DontLog()

	output, err := cmdObj.RunWithOutput()
	if err != nil {
		return nil, fmt.Errorf("gh api graphql failed: %w: %s", err, output)
	}

	var resp graphQLResponse
	if err := json.Unmarshal([]byte(output), &resp); err != nil {
		return nil, fmt.Errorf("failed to parse graphql response: %w", err)
	}

	if len(resp.Errors) > 0 {
		return nil, fmt.Errorf("graphql errors: %s", resp.Errors[0].Message)
	}

	var results []PullRequestInfo
	for i := range branchNames {
		alias := branchAlias(i)
		if repoData, ok := resp.Data[alias]; ok {
			for _, node := range repoData.PullRequests.Nodes {
				results = append(results, PullRequestInfo{
					HeadRefName: node.HeadRefName,
					Number:      node.Number,
					URL:         node.URL,
					State:       node.State,
				})
			}
		}
	}

	return results, nil
}

const graphQLMaxAliases = 50

func (self *PullRequestLoader) SetPullRequestInfoOnBranches(
	branches []*models.Branch,
	renderFunc func(),
) error {
	if !self.ghIsInstalled() {
		return nil
	}

	t := time.Now()

	owner, repoName, err := self.getRepoOwnerAndName()
	if err != nil {
		self.Log.Warnf("failed to get repo owner/name from gh: %v", err)
		return nil
	}

	var branchNames []string
	for _, branch := range branches {
		if branch.Name != "" && !branch.DetachedHead {
			branchNames = append(branchNames, branch.Name)
		}
	}

	if len(branchNames) == 0 {
		return nil
	}

	prByBranch := make(map[string]PullRequestInfo)

	for i := 0; i < len(branchNames); i += graphQLMaxAliases {
		end := i + graphQLMaxAliases
		if end > len(branchNames) {
			end = len(branchNames)
		}

		prs, err := self.fetchPRsGraphQL(owner, repoName, branchNames[i:end])
		if err != nil {
			self.Log.Warnf("PR GraphQL query failed: %v", err)
			return nil
		}

		for _, pr := range prs {
			existing, exists := prByBranch[pr.HeadRefName]
			if !exists || prPriority(pr.State) > prPriority(existing.State) {
				prByBranch[pr.HeadRefName] = pr
			}
		}
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

func prPriority(state string) int {
	switch state {
	case "OPEN":
		return 3
	case "MERGED":
		return 2
	case "CLOSED":
		return 1
	default:
		return 0
	}
}
