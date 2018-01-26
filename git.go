package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GitCountNewCommits(baseBranch string, upstream string) int {
	comparisionStr := strings.Join([]string{upstream, baseBranch}, "..")
	output, err := RunCommand("git", []string{"rev-list", "--count", comparisionStr}, false)
	if err != nil {
		fmt.Println("Could not get the difference between the feature and the upstram, are you in a git repo?")
		os.Exit(1)
	}

	i, _ := strconv.Atoi(strings.Trim(output, " \r\n"))
	return i
}

func GitCreateNewBranch(branchName string) {
	RunCommand("git", []string{"fetch"}, true)
	RunCommand("git", []string{"checkout", "master"}, true)
	RunCommand("git", []string{"reset", "--hard", "origin/master"}, true)
	RunCommand("git", []string{"checkout", "-b", branchName}, true)
}

func GitFetch() {
	RunCommand("git", []string{"fetch"}, true)
}

func GitGetCurrentBranchName() string {
	output, err := RunCommand("git", []string{"rev-parse", "--abbrev-ref", "HEAD"}, false)
	if err != nil {
		fmt.Println("Could not get the current branch name, are you in a git repo?")
		os.Exit(1)
	}

	return strings.Trim(output, " \n\r")
}

func GitGetLastCommitMessage() string {
	output, _ := RunCommand("git", []string{"log", "-1", "--pretty=%B"}, true)
	return strings.Trim(output, " \n\r")
}

func GitGetNewCommits(_baseBranch string, upstream string) string {
	output, err := RunCommand("git", []string{"cherry", upstream}, false)
	if err != nil {
		fmt.Println("Could not get the difference between the feature and the upstram, are you in a git repo?")
		os.Exit(1)
	}

	return output
}

func GitGetOriginData() (string, string) {
	output, _ := RunCommand("git", []string{"remote", "get-url", "origin", "--push"}, true)
	parts := strings.Split(strings.Split(output, ":")[1], "/")
	user := parts[0]
	repo := parts[1][0 : len(parts[1])-5]

	return user, repo
}

func GitGetRepositoryRoot() string {
	output, _ := RunCommand("git", []string{"rev-parse", "--show-toplevel"}, true)
	return strings.Trim(output, " \r\n")
}

func GitIsCurrentBranchClean() bool {
	output, err := RunCommand("git", []string{"status", "--porcelain"}, false)
	if err != nil {
		fmt.Println("Could not get branch status, are you in a git repo?")
	}

	if strings.Trim(output, " ") == "" {
		return true
	}

	return false
}

func GitPushCurrentBranch(branchName string) {
	RunCommand("git", []string{"push", "-u", "origin", branchName}, true)
}
