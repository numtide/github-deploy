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
	ignoreMissing := c.Bool("ignore-missing")

	if listScript == "" && len(selectedPullRequest) == 0 {
		return errors.New("`--list-script` or `--pull-request` missing." +
			"You have to define one or the other to select the PRs to cleanup")
	}

	var toUndeploy []string

	ctx := contextWithHandler()
	ghCli := githubClient(ctx, c)
	owner, repo := githubSlug(c)

	if len(selectedPullRequest) > 0 {
		// undeploy only selected pull requests
		toUndeploy = selectedPullRequest
	} else {
		// undeploy all closed pull requests
		var deployed []string

		deployed, err = listDeployedPullRequests(ctx, listScript)
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

	var lastErr error

	for _, name := range toUndeploy {
		log.Println("Undeploying", name)

		pullRequestID, err := strconv.Atoi(name)
		if err != nil {
			log.Println("Unable to parse pull request id: ", name)

			lastErr = err

			continue
		}

		cmd := exec.Command(c.Args().Get(0), c.Args()[1:]...) //#nosec
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Start()
		if err != nil {
			log.Println("undeploy error: ", err)

			lastErr = err

			select {
			case <-ctx.Done():
				log.Println("undeploy cancelled: ", ctx.Err())

				return lastErr
			default:
				continue
			}
		}

		err = waitOrStop(ctx, cmd)
		if err != nil {
			log.Println("undeploy error: ", err)

			lastErr = err

			select {
			case <-ctx.Done():
				log.Println("undeploy cancelled: ", ctx.Err())

				return lastErr
			default:
				continue
			}
		}

		destroyGitHubDeployments(ctx, ghCli, owner, repo, pullRequestID, ignoreMissing)
	}

	return lastErr
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
func listDeployedPullRequests(ctx context.Context, listScript string) ([]string, error) {
	var (
		stdout strings.Builder
		err    error
	)

	cmd := exec.Command(listScript)
	cmd.Stdout = &stdout

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	err = waitOrStop(ctx, cmd)
	if err != nil {
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
func destroyGitHubDeployments(ctx context.Context, ghCli *github.Client, owner string, repo string, pullRequestID int,
	ignoreMissing bool,
) {
	// Look for existing deployments related to the pull request by filtering deployments
	// by the environment name that matches the pattern 'pr-{pullRequestID}' (as the 'deploy'
	// action creates deployments with such names).
	deployments, _, err := ghCli.Repositories.ListDeployments(ctx, owner, repo, &github.DeploymentsListOptions{
		Task:        TaskName,
		Environment: fmt.Sprintf("pr-%d", pullRequestID),
	})
	if err != nil {
		log.Fatalf("Error while listing deployments for PR %d", pullRequestID)
	}

	if len(deployments) == 0 {
		if ignoreMissing {
			log.Println("No deployments found for PR ", pullRequestID)
		} else {
			log.Fatalf("unable to find deployments related to PR %d", pullRequestID)
		}
	}

	for _, deployment := range deployments {
		_, _, err := ghCli.Repositories.CreateDeploymentStatus(ctx, owner, repo,
			*deployment.ID, &github.DeploymentStatusRequest{
				State: refString("inactive"),
			})
		if err != nil {
			log.Println("Error while inactivating deployment for PR ", pullRequestID)
		}
	}
}
