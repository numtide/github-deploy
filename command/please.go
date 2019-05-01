package command

import (
	"context"
	"errors"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"gopkg.in/urfave/cli.v1"
)

const TaskName = "github-deploy"

func CmdPlease(c *cli.Context) (err error) {
	owner, repo := githubSlug(c)
	deployScript := c.String("deploy-script")
	environment := c.String("environment")
	logURL := c.String("build-url")
	prStr := c.String("pull-request")
	commit := c.GlobalString("git-commit")
	branch := c.GlobalString("git-branch")
	commitRef := c.GlobalBool("git-ref-commit")

	ref := ""

	if commitRef {
		if commit == "" {
			return errors.New("Trying to use commit as ref but commit is not set")
		} else {
			ref = commit
		}
	} else {
		if branch == "" {
			return errors.New("Trying to use branch as ref but branch is not set")
		} else {
			ref = branch
		}
	}

	var pr int
	if prStr != "" && prStr != "false" {
		// prStr might be a URL, in that case pull the last component of the path
		strs := strings.Split(prStr, "/")
		prStr = strs[len(strs)-1]

		pr, err = strconv.Atoi(prStr)
		if err != nil {
			return err
		}
	}

	// If the environment is not set, set as follows:
	//   * branch is master: production
	//   * otherwise: pr-preview
	if environment == "" {
		if branch == "master" {
			environment = "production"
		} else {
			environment = "pr-preview"
		}
	}

	ctx := context.Background()
	gh := githubClient(ctx, c)

	log.Println("deploy ref", ref)
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
			TransientEnvironment:  refBool(pr > 0),
			ProductionEnvironment: refBool(pr == 0),
		})
		if err != nil {
			return err
		}

	}

	// Prepare deploy script
	var stdout strings.Builder
	cmd := exec.Command(deployScript, environment)
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
