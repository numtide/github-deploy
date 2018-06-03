package main

import (
	"fmt"
	"os"

	"github.com/zimbatm/github-deploy/command"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

var GlobalFlags = []cli.Flag{
	// This is only really needed for the "please" command
	altsrc.NewStringFlag(cli.StringFlag{
		Name:  "git-commit",
		Usage: "git commit ID",
		//EnvVar: "TRAVIS_COMMIT,BUILDKITE_COMMIT",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:  "git-origin",
		Usage: "URL of the repo",
		//EnvVar: "BUILDKITE_REPO", // Travis doesn't have an equivalent
	}),
	cli.StringFlag{
		Name:   "github-token",
		Usage:  "Github Personal access token to interact with the Github API",
		EnvVar: "GITHUB_AUTH_TOKEN",
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
				Usage: "Script that deploys the given PR",
			},
			cli.StringFlag{
				Name:   "pr, pull-request",
				Usage:  "Creates a temporary deployment for the give pull-request ID",
				EnvVar: "TRAVIS_PULL_REQUEST,BUILDKITE_PULL_REQUEST",
			},
			cli.StringFlag{
				Name:  "environment",
				Value: "production",
				Usage: "Sets the target environment, ignored if pull-request is passed",
			},
			cli.StringFlag{
				Name:   "build-url",
				Usage:  "URL to follow the build progress",
				EnvVar: "BUILDKITE_BUILD_URL", // NOTE: Travis doesn't have an equivalent
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
