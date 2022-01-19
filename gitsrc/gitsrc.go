// gitsrc satisfies the urfave/cli/altsrc.InputSourceContext interface
package gitsrc

import (
	"errors"
	"fmt"
	"time"

	git "github.com/go-git/go-git/v5"
	cli "gopkg.in/urfave/cli.v1"
	altsrc "gopkg.in/urfave/cli.v1/altsrc"
)

// ErrNotSupported is returned by all the functions that are not String().
var ErrNotSupported = errors.New("gitsrc: operation not supported")

// ErrNotFound is returned when the String(name) doesn't have a reference.
var ErrNotFound = errors.New("gitsrc: key not found")

// ErrNotBranch is returned when the git repo is not on a branch.
var ErrNotBranch = errors.New("gitsrc: ref is not a branch")

// FromCurrentDir tries to open $PWD as the git repo.
func FromCurrentDir(*cli.Context) (altsrc.InputSourceContext, error) { //nolint:ireturn
	r, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})

	return &gitSource{r}, err
}

type gitSource struct {
	*git.Repository
}

func (x *gitSource) String(name string) (string, error) {
	// These names need to be aligned with the top-level GlobalFlags
	switch name {
	case "git-origin":
		remote, err := x.Remote("origin")
		if err != nil {
			return "", fmt.Errorf("gitsrc: %w", err)
		}

		return remote.Config().URLs[0], nil

	case "git-commit":
		ref, err := x.Head()
		if err != nil {
			return "", fmt.Errorf("gitsrc: %w", err)
		}

		return ref.Hash().String(), nil

	case "git-branch":
		ref, err := x.Head()
		if err != nil {
			return "", fmt.Errorf("gitsrc: %w", err)
		}

		refName := ref.Name()
		if refName.IsBranch() {
			return refName.String(), nil
		}

		return "", ErrNotBranch
	}

	return "", ErrNotFound
}

// These are implemented to satisfy the altsrc.InputSourceContext interface

func (x *gitSource) Int(name string) (int, error) {
	return 0, ErrNotSupported
}

func (x *gitSource) Duration(name string) (time.Duration, error) {
	return 0, ErrNotSupported
}

func (x *gitSource) Float64(name string) (float64, error) {
	return 0, ErrNotSupported
}

func (x *gitSource) StringSlice(name string) ([]string, error) {
	return nil, ErrNotSupported
}

func (x *gitSource) IntSlice(name string) ([]int, error) {
	return nil, ErrNotSupported
}

func (x *gitSource) Generic(name string) (cli.Generic, error) { //nolint:ireturn
	return nil, ErrNotSupported
}

func (x *gitSource) Bool(name string) (bool, error) {
	return false, ErrNotSupported
}

func (x *gitSource) BoolT(name string) (bool, error) {
	return false, ErrNotSupported
}
