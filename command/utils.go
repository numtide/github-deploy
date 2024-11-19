package command

import (
	"context"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/google/go-github/github"
	secretvalue "github.com/zimbatm/go-secretvalue"
	"golang.org/x/oauth2"
	cli "gopkg.in/urfave/cli.v1"
)

func githubClient(ctx context.Context, c *cli.Context) *github.Client {
	token := c.GlobalGeneric("github-token").(*secretvalue.StringFlag).SecretValue
	if !token.IsSet() {
		log.Fatal("missing github auth token")
	}

	// log.Println("github auth token", token)
	// TODO: determine the right API based on c.GlobalString("git-origin")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.GetString()},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

var ReSlug = regexp.MustCompile(`[:/]([\w-]+)/([\w-]+)`)

func githubSlug(c *cli.Context) (string, string) {
	origin := c.GlobalString("git-origin")
	matches := ReSlug.FindStringSubmatch(origin)

	if len(matches) < 3 {
		return "", ""
	}

	return matches[1], matches[2]
}

const (
	StateError    = "error"
	StateFailure  = "failure"
	StateInactive = "inactive"
	StatePending  = "pending"
	StateSuccess  = "success"
)

// Move things to the heap

func refBool(b bool) *bool {
	return &b
}

func refString(str string) *string {
	return &str
}

func refStringList(l []string) *[]string {
	return &l
}

func propagateSignalsTo(cmd *exec.Cmd, signalChannel chan os.Signal) {
	for sig := range signalChannel {
		if cmd.Process != nil {
			err := cmd.Process.Signal(sig)
			if err != nil {
				log.Printf("error sending signal to child process (%d): %s\n", cmd.Process.Pid, err)
			}
		} else {
			// TODO: is this always the right thing to do if we're not running a subprocess?
			os.Exit(1)
		}
	}
}
