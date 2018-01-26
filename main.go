package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alecthomas/kingpin"
)

var (
	app = kingpin.New(
		"feature-up", "A tool to facilitate feature development in SumUp")
	debug = app.Flag(
		"debug", "Enable debug mode").Bool()

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
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case create.FullCommand():
		runCreateCmd()
	case review.FullCommand():
		runRequestReviewCmd()
	}
}

func runCreateCmd() {
	fmt.Println("Creating feature", *createTicket, "...")
	branchName := strings.Join([]string{"features", *createTicket}, "/")
	GitCreateNewBranch(branchName)

	fmt.Println("Feature branch created, have fun!")
	fmt.Println("After you are done, please run `feature-up review` to create a PR")
}

func runRequestReviewCmd() {
	fmt.Println("Requesting a new code review ...")

	branchName := GitGetCurrentBranchName()
	ticket := strings.Split(branchName, "/")[1]
	user, repo := GitGetOriginData()
	remoteUpstream := strings.Join([]string{"origin", *reviewUpstream}, "/")
	newCommits := GitGetNewCommits(branchName, remoteUpstream)
	newCommitsCount := GitCountNewCommits(branchName, remoteUpstream)
	commitMessage := GitGetLastCommitMessage()

	GitFetch()

	if GitIsCurrentBranchClean() == false && *reviewForce == false {
		fmt.Println("You have stuff to commit in this branch, please do so before calling me.")
		os.Exit(1)
	}

	if newCommitsCount == 0 {
		fmt.Println("There is no difference between the feature and the upstream, please do some work")
		os.Exit(1)
	}

	fmt.Println("")
	fmt.Println("Running unit tests...")
	RunUnitTests()

	fmt.Println("")
	fmt.Println("New commits:")
	fmt.Println(newCommits)

	if newCommitsCount > 1 {
		rebaseStartPositionPieces := []string{"HEAD~", strconv.Itoa(newCommitsCount)}
		rebaseStartPosition := strings.Join(rebaseStartPositionPieces, "")

		fmt.Println("There is more than one commit to be applied on the upstream, please squash your stuff")
		fmt.Println("Hint: git rebase -i", rebaseStartPosition)
		os.Exit(1)
	}

	if commitMessage[1:len(ticket)+1] != ticket {
		fmt.Println("The commit message must start with [" + ticket + "], please change")
		fmt.Println("Hint: git commit --amend")
		os.Exit(1)
	}

	GitPushCurrentBranch(branchName)

	fmt.Println("")
	fmt.Println("Creating pull request...")
	pullRequest := CreatePR(
		user,
		repo,
		commitMessage,
		commitMessage,
		branchName,
		"master")

	fmt.Println("Pull request created, now just post this link in the #code-review channel: ", *pullRequest.HTMLURL)
}
