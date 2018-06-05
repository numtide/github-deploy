package command

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/github"
	"gopkg.in/urfave/cli.v1"
)

const TaskName = "github-deploy"

func CmdPlease(c *cli.Context) (err error) {
	owner, repo := githubSlug(c)
	deployScript := c.String("deploy-script")
	ref := c.GlobalString("git-commit")
	pr := c.String("pull-request")
	logURL := c.String("build-url")
	environment := c.String("environment")
	if pr != "" {
		environment = fmt.Sprintf("review-%s", pr)
	}

	ctx := context.Background()
	gh := githubClient(ctx, c)

	log.Println("commit ID", ref)
	log.Println("origin", c.GlobalString("git-origin"))

	// First, declare the new deployment to GitHub

	// Look for an existing deployment
	deployments, _, err := gh.Repositories.ListDeployments(ctx, owner, repo, &github.DeploymentsListOptions{
		Ref:  ref,
		Task: TaskName,
	})
	if err != nil {
		return err
	}
	var deployment *github.Deployment
	if len(deployments) > 0 {
		deployment = deployments[0]
	} else {
		deployment, _, err = gh.Repositories.CreateDeployment(ctx, owner, repo, &github.DeploymentRequest{
			Ref:                   refString(ref),
			Task:                  refString(TaskName),
			AutoMerge:             refBool(false),
			RequiredContexts:      refStringList([]string{}),
			Payload:               refString("{}"),
			Environment:           refString(environment),
			Description:           refString(TaskName),
			TransientEnvironment:  refBool(pr != ""),
			ProductionEnvironment: refBool(pr == ""),
		})
		if err != nil {
			return err
		}

	}

	// Prepare deploy script
	var stdout strings.Builder
	cmd := exec.Command(deployScript, fmt.Sprintf("pr-%s", pr))
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	updateStatus := func(state string, environmentURL string) {
		_, _, err := gh.Repositories.CreateDeploymentStatus(ctx, owner, repo, *deployment.ID, &github.DeploymentStatusRequest{
			State:          refString(state),
			LogURL:         refString(logURL),
			Description:    refString(TaskName),
			EnvironmentURL: refString(environmentURL),
			//AutoInactive: refBool(true),
		})
		if err != nil {
			log.Println("updateStatus:", err)
		}
	}

	// Start deploy script
	err = cmd.Start()
	if err != nil {
		updateStatus(StateError, "")
		return err
	}

	// Record progress
	updateStatus(StatePending, "")

	// Wait on the deploy to finish
	err = cmd.Wait()
	if err != nil {
		updateStatus(StateFailure, "")
		return err
	}

	// Success!
	out := strings.SplitN(stdout.String(), "\n", 2)
	environmentURL := out[0]
	updateStatus(StateSuccess, environmentURL)

	return nil
}
