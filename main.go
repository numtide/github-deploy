package main

import (
	"log"
	"os"

	"github.com/zimbatm/github-deploy/gitsrc"
	cli "gopkg.in/urfave/cli.v1"
	altsrc "gopkg.in/urfave/cli.v1/altsrc"
)

func main() {
	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "zimbatm"
	app.Email = "zimbatm@zimbatm.com"
	app.Usage = "A CLI that integrates deployments with github"
	app.EnableBashCompletion = true
	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	// Also load info from the current git repo
	app.Before = altsrc.InitInputSourceWithContext(GlobalFlags, gitsrc.FromCurrentDir)

	if err := app.Run(os.Args); err != nil {
		log.Fatal("ERROR:", err)
	}
}
