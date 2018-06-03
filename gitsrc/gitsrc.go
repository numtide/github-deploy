// gitsrc satisfies the urfave/cli/altsrc.InputSourceContext interface
package gitsrc

import (
	"fmt"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

// ErrNotSupported is returned by all the functions that are not String()
var ErrNotSupported = fmt.Errorf("operation not supported")

// ErrNotFound is returns when the String(name) doesn't have a reference
var ErrNotFound = fmt.Errorf("key not found")

// FromCurrentDir tries to open $PWD as the git repo
func FromCurrentDir(*cli.Context) (altsrc.InputSourceContext, error) {
	r, err := git.PlainOpen(".")
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
			return "", err
		}
		return remote.Config().URLs[0], nil
	case "git-commit":
		ref, err := x.Head()
		if err != nil {
			return "", err
		}
		return ref.Hash().String(), nil
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

func (x *gitSource) Generic(name string) (cli.Generic, error) {
	return nil, ErrNotSupported
}

func (x *gitSource) Bool(name string) (bool, error) {
	return false, ErrNotSupported
}

func (x *gitSource) BoolT(name string) (bool, error) {
	return false, ErrNotSupported
}
