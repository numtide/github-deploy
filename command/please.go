package command

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	cli "gopkg.in/urfave/cli.v1"
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
	origin := c.GlobalString("git-origin")
	environmentURL := c.String("environment-url")

	if deployScript != "" {
		return fmt.Errorf("--deploy-script is deprecated, use a positional argument instead")
	}
	if c.NArg() == 0 {
		return fmt.Errorf("Missing the deploy script as a positional argument")
	}

	// Compose the Git originl URL in the case of GitHub Actions
	if origin == "" && os.Getenv("GITHUB_SERVER_URL") != "" {
		origin = fmt.Sprintf(
			"%s/%s.git",
			os.Getenv("GITHUB_SERVER_URL"),
			os.Getenv("GITHUB_REPOSITORY"),
		)
	}

	// Compose the log URL in the case of GitHub Actions
	if logURL == "" && os.Getenv("GITHUB_SERVER_URL") != "" {
		logURL = fmt.Sprintf(
			"%s/%s/actions/runs/%s",
			os.Getenv("GITHUB_SERVER_URL"),
			os.Getenv("GITHUB_REPOSITORY"),
			os.Getenv("GITHUB_RUN_ID"),
		)
	}

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

	// Override the deployment target on pull-request
	if pr > 0 {
		environment = fmt.Sprintf("pr-%d", pr)
	}

	ctx := context.Background()
	gh := githubClient(ctx, c)

	log.Println("deploy ref", ref)
	log.Println("origin", origin)

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

	deployScriptSubStrings := strings.Fields(deployScript)
	if len(deployScriptSubStrings) == 1 {
		deployScriptSubStrings = append(deployScriptSubStrings, environment)
	}

	// Prepare deploy script
	var stdout bytes.Buffer
	cmd := exec.Command(c.Args().Get(0), c.Args()[1:]...)
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = os.Stderr

	updateStatus := func(state string, environmentURL string) error {
		_, _, err := gh.Repositories.CreateDeploymentStatus(ctx, owner, repo, *deployment.ID, &github.DeploymentStatusRequest{
			State:          refString(state),
			LogURL:         refString(logURL),
			Description:    refString(TaskName),
			EnvironmentURL: refString(environmentURL),
			//AutoInactive: refBool(true),
		})
		return err
	}

	// Start deploy script
	err = cmd.Start()
	if err != nil {
		err2 := updateStatus(StateError, "")
		if err2 != nil {
			log.Println("updateStatus:", err)
		}
		return err
	}

	// Record progress
	err = updateStatus(StatePending, "")
	if err != nil {
		return err
	}

	// Wait on the deploy to finish
	err = cmd.Wait()
	if err != nil {
		err2 := updateStatus(StateFailure, "")
		if err2 != nil {
			log.Println("updateStatus:", err)
		}
		return err
	}

	// Success!
	out := strings.SplitN(stdout.String(), "\n", 2)
	if environmentURL == "" {
		environmentURL = strings.TrimSpace(out[0])
	}
	err = updateStatus(StateSuccess, environmentURL)
	return err
}
