package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
)

var LoadedConfig *Config

var (
	app = kingpin.New(
		"feature-up", "A tool to facilitate feature development in SumUp")
	debug = app.Flag(
		"debug", "Enable debug mode").Bool()

	setup = app.Command("setup", "Setup feature-up")

	create = app.Command(
		"create", "Initialize a new feature")
	createTicket = create.Arg(
		"ticket", "The ID of the ticket in JIRA (e.g.: SSMR-987)").
		Required().
		String()

	review = app.Command(
		"review", "Request a code review in the current feature branch")
	reviewUpstream = review.Arg(
		"upstream", "The upstream to compare your feature with (theta, beta, ...)").
		Required().
		String()
	reviewForce = review.Flag(
		"force", "Forces the script to bypass errors when there are uncommited stuff").
		Bool()
)

func main() {
	config, err := ReadConfigFromFile()
	if err != nil {
		fmt.Println("Could not find config file, starting setup...")
		config = RunSetupCommand()
	}

	LoadedConfig = config

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case setup.FullCommand():
		RunSetupCommand()
	case create.FullCommand():
		RunCreateCommand()
	case review.FullCommand():
		RunReviewCommand()
	}
}
