package main

import (
	"fmt"
	"os"

	"github.com/zimbatm/github-deploy/command"
	"github.com/zimbatm/go-secretvalue"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

var GlobalFlags = []cli.Flag{
	// This is only really needed for the "please" command
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "git-commit",
		Usage:  "git commit ID",
		EnvVar: "GITHUB_SHA,BUILDKITE_COMMIT,CIRCLE_SHA1,TRAVIS_PULL_REQUEST_SHA",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "git-branch",
		Usage:  "git branch",
		EnvVar: "GITHUB_REF,BUILDKITE_BRANCH,CIRCLE_BRANCH,TRAVIS_BRANCH",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:  "git-origin",
		Usage: "URL of the repo",
		// NOTE: In the case of GitHub Actions, there is no env var that provides
		//       this directly.
		EnvVar: "BUILDKITE_REPO,CIRCLE_REPOSITORY_URL", // Travis doesn't have an equivalent
	}),
	cli.BoolFlag{
		Name:  "git-ref-commit",
		Usage: "use the commit as deployment reference instead of branch",
	},
	cli.GenericFlag{
		Name:   "github-token",
		Usage:  "Github Personal access token to interact with the Github API",
		EnvVar: "GITHUB_TOKEN",
		Value: &secretvalue.StringFlag{
			SecretValue: secretvalue.New("github-token"),
		},
	},
}

var Commands = []cli.Command{
	{
		Name:   "please",
		Usage:  "Initiates a deployment",
		Action: command.CmdPlease,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "deploy-script",
				Usage: "DEPRECATED. Use a positional argument instead.",
			},
			cli.StringFlag{
				Name:  "pr, pull-request",
				Usage: "Creates a temporary deployment for the give pull-request ID",
				// NOTE: GitHub Actions doesn't have an env var like that and the
				//       argument must be passed explicitly.
				EnvVar: "BUILDKITE_PULL_REQUEST,CIRCLE_PULL_REQUEST,TRAVIS_PULL_REQUEST",
			},
			cli.StringFlag{
				Name:  "environment",
				Value: "production",
				Usage: "Sets the target environment, ignored if pull-request is passed",
			},
			cli.StringFlag{
				Name:  "build-url",
				Usage: "URL to follow the build progress",
				// NOTE: Travis doesn't have an equivalent
				// NOTE: For GitHub Actions, the URL is composed later in the command
				//       if empty.
				EnvVar: "BUILDKITE_BUILD_URL,CIRCLE_BUILD_URL",
			},
		},
	},
	{
		Name:   "cleanup",
		Usage:  "Removes old temporary deployments",
		Action: command.CmdCleanup,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "list-script",
				Usage: "Script that lists the deployed PRs",
			},
			cli.StringFlag{
				Name:  "undeploy-script",
				Usage: "Script that deleted a deployment given a specific PR",
			},
		},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
