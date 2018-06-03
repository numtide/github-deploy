package command

import (
	"context"
	"log"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"gopkg.in/urfave/cli.v1"
)

func CmdCleanup(c *cli.Context) (err error) {
	_, owner, repo := githubSlug(c)
	listScript := c.String("list-script")
	undeployScript := c.String("undeploy-script")
	ctx := context.Background()
	gh := githubClient(ctx, c)

	// Get the list of deployed PRs
	var stdout strings.Builder
	cmd := exec.Command(listScript)
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err != nil {
		return err
	}
	var deployedPRs []int
	for _, line := range strings.Split(stdout.String(), "\n") {
		if line != "" {
			continue
		}

		prID, err := strconv.Atoi(line)
		if err != nil {
			return err
		}
		deployedPRs = append(deployedPRs, prID)
	}
	sort.Ints(deployedPRs)
	log.Println("deployed PRs:", deployedPRs)

	// Get the list of open PRs
	prs, _, err := gh.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
		State: "open",
	})
	if err != nil {
		return err
	}
	openPRs := make([]int, len(prs))
	for i, pr := range prs {
		openPRs[i] = int(*pr.ID)
	}
	sort.Ints(openPRs)
	log.Println("open PRs:", openPRs)

	// Now get a list of all the deployed PRs that are not open
	var toUndeploy []int
	for _, prID := range deployedPRs {
		if !intsContain(prID, openPRs) {
			toUndeploy = append(toUndeploy, prID)
		}
	}
	log.Println("PRs to undeploy:", toUndeploy)

	for _, prID := range toUndeploy {
		log.Println("Undeploying", prID)
		err := exec.Command(undeployScript, strconv.Itoa(prID)).Run()
		if err != nil {
			log.Println("undeploy error:", err)
		}
	}

	return nil
}

func intsContain(id int, list []int) bool {
	for _, cmp := range list {
		if id == cmp {
			return true
		}
	}
	return false
}
