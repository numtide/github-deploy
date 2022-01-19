package command

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	cli "gopkg.in/urfave/cli.v1"
)

func CmdCleanup(c *cli.Context) (err error) {
	owner, repo := githubSlug(c)
	listScript := c.String("list-script")
	selectedPullRequest := c.StringSlice("pull-request")
	if listScript == "" && len(selectedPullRequest) == 0 {
		return errors.New("`--list-script` or `--pull-request` missing. You have to define one or the other to select the PRs to cleanup !")
	}
	ctx := context.Background()
	gh := githubClient(ctx, c)

	var toUndeploy []string
	if len(selectedPullRequest) > 0 {
		// undeploy only selected pull requests
		toUndeploy = selectedPullRequest
	} else {
		// undeploy all closed pull requests
		var deployed []string
		deployed, err = listDeployedPullRequests(listScript)
		if err != nil {
			return err
		}

		// Get the list of open PRs
		prs, _, err := gh.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
			State: "open",
		})
		if err != nil {
			return err
		}
		openPRs := make([]string, len(prs))
		for i, pr := range prs {
			openPRs[i] = fmt.Sprintf("pr-%d", *pr.Number)
		}
		log.Println("open PRs:", openPRs)

		// Now get a list of all the deployed PRs that are not open
		for _, name := range deployed {
			if !contains(name, openPRs) {
				toUndeploy = append(toUndeploy, name)
			}
		}
	}
	log.Println("to undeploy:", toUndeploy)

	for _, name := range toUndeploy {
		log.Println("Undeploying", name)
		cmd := exec.Command(c.Args().Get(0), c.Args()[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Println("undeploy error: ", err)
			continue
		}
		pullRequestId, err := strconv.Atoi(name)
		if err != nil {
			log.Println("Unable to parse pull request id: ", name)
			continue
		}
		destroyGitHubDeployments(ctx, gh, owner, repo, pullRequestId)
	}

	return nil
}

func contains(item string, list []string) bool {
	for _, entry := range list {
		if item == entry {
			return true
		}
	}
	return false
}

// Get the list of deployed Pull request based on given script
func listDeployedPullRequests(listScript string) ([]string, error) {
	var stdout strings.Builder
	cmd := exec.Command(listScript)
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	var deployed []string
	for _, line := range strings.Split(stdout.String(), "\n") {
		if line == "" {
			continue
		}

		deployed = append(deployed, line)
	}
	log.Println("deployed:", deployed)
	return deployed, nil
}

// destroy deployments related to a PR by marking them
func destroyGitHubDeployments(ctx context.Context, gh *github.Client, owner string, repo string, pullRequestId int) {
	pr, _, err := gh.PullRequests.Get(ctx, owner, repo, pullRequestId)
	if err != nil {
		log.Fatalf("Unable to fetch from github pull request id: %d", pullRequestId)
	}

	// Look for an existing deployments
	// We filter deployments related to a PR based on the PR head branch name (as the `deploy` creates them)
	deployments, _, err := gh.Repositories.ListDeployments(ctx, owner, repo, &github.DeploymentsListOptions{
		Ref:  fmt.Sprintf("refs/heads/%s", *pr.Head.Ref),
		Task: TaskName,
	})
	if err != nil || len(deployments) == 0 {
		log.Println("unable to find deployments related to PR ", pr.ID)
	}
	for _, deployment := range deployments {
		_, _, err := gh.Repositories.CreateDeploymentStatus(ctx, owner, repo, *deployment.ID, &github.DeploymentStatusRequest{
			State: refString("inactive"),
		})
		if err != nil {
			log.Println("Error while inactivating deployment for PR:", pr.ID)
		}
	}
}
