package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func RunCreateCommand() {
	fmt.Println("Creating feature", *createTicket, "...")
	branchName := strings.Join([]string{"features", *createTicket}, "/")
	GitCreateNewBranch(branchName)

	fmt.Println("Feature branch created, have fun!")
	fmt.Println("After you are done, please run `feature-up review` to create a PR")
}

func RunReviewCommand() {
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

func RunSetupCommand() *Config {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("In order to use this tool, you must create a Github personal access token.")
	fmt.Println("You can do so by accessing https://github.com/settings/tokens")
	fmt.Println("Please note that you must include the `repo` scope")
	fmt.Println("")
	fmt.Print("Github personal access token:")
	token, _ := reader.ReadString('\n')

	config := &Config{
		GithubToken: strings.Trim(token, " \r\n"),
	}

	config.Save()
	return config
}
