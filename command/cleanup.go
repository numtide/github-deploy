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
	listScript := c.String("list-script")
	selectedPullRequest := c.StringSlice("pull-request")

	if listScript == "" && len(selectedPullRequest) == 0 {
		return errors.New("`--list-script` or `--pull-request` missing." +
			"You have to define one or the other to select the PRs to cleanup")
	}

	var toUndeploy []string

	ctx := context.Background()
	ghCli := githubClient(ctx, c)
	owner, repo := githubSlug(c)

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
		prs, _, err := ghCli.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
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

		cmd := exec.Command(c.Args().Get(0), c.Args()[1:]...) //#nosec
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			log.Println("undeploy error: ", err)

			continue
		}

		pullRequestID, err := strconv.Atoi(name)
		if err != nil {
			log.Println("Unable to parse pull request id: ", name)

			continue
		}

		destroyGitHubDeployments(ctx, ghCli, owner, repo, pullRequestID)
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

// Get the list of deployed Pull request based on given script.
func listDeployedPullRequests(listScript string) ([]string, error) {
	var stdout strings.Builder

	cmd := exec.Command(listScript)
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	lines := strings.Split(stdout.String(), "\n")
	deployed := make([]string, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		deployed = append(deployed, line)
	}

	log.Println("deployed:", deployed)

	return deployed, nil
}

// Destroy deployments related to a PR by marking them.
func destroyGitHubDeployments(ctx context.Context, ghCli *github.Client, owner string, repo string, pullRequestID int) {
	pr, _, err := ghCli.PullRequests.Get(ctx, owner, repo, pullRequestID)
	if err != nil {
		log.Fatalf("Unable to fetch from github pull request id: %d", pullRequestID)
	}

	// Look for existing deployments related to the pull request by filtering deployments
	// by the environment name that matches the pattern 'pr-{pullRequestID}' (as the 'deploy'
	// action creates deployments with such names).
	deployments, _, err := ghCli.Repositories.ListDeployments(ctx, owner, repo, &github.DeploymentsListOptions{
		Task:        TaskName,
		Environment: fmt.Sprintf("pr-%d", pullRequestID),
	})
	if err != nil {
		log.Fatalf("Error while listing deployments for PR %d", pr.Number)
	}

	if len(deployments) == 0 {
		log.Fatalf("unable to find deployments related to PR %d", pr.Number)
	}

	for _, deployment := range deployments {
		_, _, err := ghCli.Repositories.CreateDeploymentStatus(ctx, owner, repo,
			*deployment.ID, &github.DeploymentStatusRequest{
				State: refString("inactive"),
			})
		if err != nil {
			log.Println("Error while inactivating deployment for PR:", pr.ID)
		}
	}
}
