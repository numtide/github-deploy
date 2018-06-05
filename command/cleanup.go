package command

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/google/go-github/github"
	"gopkg.in/urfave/cli.v1"
)

func CmdCleanup(c *cli.Context) (err error) {
	owner, repo := githubSlug(c)
	listScript := c.String("list-script")
	undeployScript := c.String("undeploy-script")
	ctx := context.Background()
	gh := githubClient(ctx, c)

	// Get the list of temporary deployments
	var stdout strings.Builder
	cmd := exec.Command(listScript)
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err != nil {
		return err
	}
	var deployed []string
	for _, line := range strings.Split(stdout.String(), "\n") {
		if line == "" {
			continue
		}

		deployed = append(deployed, line)
	}
	log.Println("deployed:", deployed)

	// Get the list of open PRs
	prs, _, err := gh.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
		State: "open",
	})
	if err != nil {
		return err
	}
	openPRs := make([]string, len(prs))
	for i, pr := range prs {
		openPRs[i] = fmt.Sprintf("pr-%d", *pr.ID)
	}
	log.Println("open PRs:", openPRs)

	// Now get a list of all the deployed PRs that are not open
	var toUndeploy []string
	for _, name := range deployed {
		if !contains(name, openPRs) {
			toUndeploy = append(toUndeploy, name)
		}
	}
	log.Println("to undeploy:", toUndeploy)

	for _, name := range toUndeploy {
		log.Println("Undeploying", name)
		err := exec.Command(undeployScript, name).Run()
		if err != nil {
			log.Println("undeploy error:", err)
		}
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
