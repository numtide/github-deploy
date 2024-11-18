package command

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"syscall"
	"time"

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

var DefaultKillDelay = 5 * time.Minute

// waitOrStop waits for the already-started command cmd by calling its Wait method.
//
// If cmd does not return before ctx is done, waitOrStop sends it the given interrupt signal.
// waitOrStop waits DefaultKillDelay for Wait to return before sending os.Kill.
//
// This function is copied from the one added to x/playground/internal in
// http://golang.org/cl/228438.
func waitOrStop(ctx context.Context, cmd *exec.Cmd) error {
	if cmd.Process == nil {
		panic("waitOrStop called with a nil cmd.Process — missing Start call?")
	}

	errc := make(chan error)
	go func() {
		select {
		case errc <- nil:
			return
		case <-ctx.Done():
		}

		err := cmd.Process.Signal(os.Interrupt)
		if err == nil {
			err = ctx.Err() // Report ctx.Err() as the reason we interrupted.
		} else if err.Error() == "os: process already finished" {
			errc <- nil

			return
		}

		if DefaultKillDelay > 0 {
			timer := time.NewTimer(DefaultKillDelay)
			select {
			// Report ctx.Err() as the reason we interrupted the process...
			case errc <- ctx.Err():
				timer.Stop()

				return
			// ...but after killDelay has elapsed, fall back to a stronger signal.
			case <-timer.C:
			}

			// Wait still hasn't returned.
			// Kill the process harder to make sure that it exits.
			//
			// Ignore any error: if cmd.Process has already terminated, we still
			// want to send ctx.Err() (or the error from the Interrupt call)
			// to properly attribute the signal that may have terminated it.
			_ = cmd.Process.Kill()
		}

		errc <- err
	}()

	waitErr := cmd.Wait()

	interruptErr := <-errc
	if interruptErr != nil {
		return interruptErr
	}

	return waitErr
}

// contextWithHandler returns a context that is canceled when the program receives a SIGINT or SIGTERM.
//
// !! Only call this function once per program.
func contextWithHandler() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		log.Printf("Received signal %s, stopping", sig)
		cancel()
	}()

	return ctx
}
