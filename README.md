# github-deploy - Track deployments on GitHub

[![Build Status](https://travis-ci.com/zimbatm/github-deploy.svg?branch=master)](https://travis-ci.com/zimbatm/github-deploy)
[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)

An opinionated command-line utility that integrates deployments with the github Deployment API.

## Description

This is a wrapper command that abstracts the deployment method through a set of scripts which interfaces are clearly defined.

The wrapper command tracks the deployment statuses by interacting with the github API.

## Assumptions

The command is being run in the checkout of the project that is about to be deployed, with the right
git commit checked out.

## Usage

```
NAME:
   github-deploy - A CLI that integrates deployments with github

USAGE:
   github-deploy [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR:
   zimbatm <zimbatm@zimbatm.com>

COMMANDS:
     please   Initiates a deployment
     cleanup  Removes old temporary deployments
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --git-commit value    git commit ID
   --git-origin value    URL of the repo
   --github-token value  Github Personal access token to interact with the Github API [$GITHUB_AUTH_TOKEN]
   --help, -h            show help
   --version, -v         print the version
```

## Scripts interface

### `DEPLOY_SCRIPT [PR]`

The deploy script takes an optional PR ID as an argument and returns the target URL on stdout.

Depending on the script exit status, the deployment is marked as a failure or success.

### `LIST_SCRIPT`

Returns the list of PR IDs of all the temporary deployed application on stdout.

### `UNDEPLOY_SCRIPT <PR>`

Deletes a deployment mapping to the given PR ID.

## Install

To install, use `go get`:

```bash
$ go get -d github.com/zimbatm/github-deploy
```

## Setup

### Create a token

Go to https://github.com/settings/tokens/new

Select `repo`

export GITHUB_AUTH_TOKEN=<new-token>

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
