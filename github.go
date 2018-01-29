package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func getClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: LoadedConfig.GithubToken})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func CreatePR(user string, repo string, title string, description string, headBranch string, baseBranch string) *github.PullRequest {
	client := getClient()
	pullRequestData := &github.NewPullRequest{
		Title: &title,
		Head:  &headBranch,
		Base:  &baseBranch,
		Body:  &title,
	}

	pullRequest, _, err := client.PullRequests.Create(context.Background(), user, repo, pullRequestData)
	if err != nil {
		fmt.Println("Error creating the pull request:", err)
		os.Exit(1)
	}

	return pullRequest
}
