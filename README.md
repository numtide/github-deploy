# github-deploy - Track deployments on GitHub

[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)

An opinionated command-line utility that integrates deployments with the github Deployment API.

## Description

This is a wrapper command that abstracts the deployment method through a set of scripts which interfaces are clearly defined.

The wrapper command tracks the deployment statuses by interacting with the github API.

## Assumptions

The command is being run in the checkout of the project that is about to be deployed, with the right
git commit checked out.

## Usage

`$ ./github-deploy --help`
```
NAME:
   github-deploy - A CLI that integrates deployments with github

USAGE:
   github-deploy [global options] command [command options] [arguments...]

VERSION:
   0.6.0

AUTHOR:
   zimbatm <zimbatm@zimbatm.com>

COMMANDS:
     please   Initiates a deployment
     cleanup  Removes deployments
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --git-commit value    git commit ID [$GITHUB_SHA, $BUILDKITE_COMMIT, $CIRCLE_SHA1, $TRAVIS_PULL_REQUEST_SHA]
   --git-branch value    git branch [$GITHUB_REF, $BUILDKITE_BRANCH, $CIRCLE_BRANCH, $TRAVIS_BRANCH]
   --git-origin value    URL of the repo [$BUILDKITE_REPO, $CIRCLE_REPOSITORY_URL]
   --git-ref-commit      use the commit as deployment reference instead of branch
   --github-token value  Github Personal access token to interact with the Github API (default: <secret:github-token>) [$GITHUB_TOKEN]
   --help, -h            show help
   --version, -v         print the version
```
## Scripts interface

### `DEPLOY_SCRIPT <TARGET>`

The deploy script takes an optional deployment name an argument and returns the target URL on stdout.

Depending on the script exit status, the deployment is marked as a failure or success.

### `LIST_SCRIPT`

Returns the list of all the temporary deployments on stdout.

### `UNDEPLOY_SCRIPT <TARGET>`

Deletes a deployment named `<TARGET>`. Should not undeploy production.

## Install

To install, use `go get`:

```bash
$ go get -d github.com/zimbatm/github-deploy
```

## Setup

### Create a token

Go to https://github.com/settings/tokens/new

Select `repo`

export GITHUB_TOKEN=<new-token>

### Create the wrapper scripts

TODO example

## Contribution

1. Fork ([https://github.com/zimbatm/github-deploy/fork](https://github.com/zimbatm/github-deploy/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[zimbatm](https://github.com/zimbatm)

## License

MIT
