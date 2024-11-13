package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v65/github"
	"github.com/pterm/pterm"
	// "github.com/pterm/pterm/putils"
)

func main() {
	startPage := 31
	// maxPages := 3 + startPage
	maxPages := 500

	donorDetails := []string{"solo-io", "gloo"}
	destDetails := []string{"solo-io", "gloo_2"}

	logger := pterm.DefaultLogger.
		WithLevel(pterm.LogLevelDebug)

	logger.Info("Starting Issue grabber",
		logger.Args("source", donorDetails, "destination", destDetails))

	ghToken := os.Getenv("GITHUB_TOKEN")
	if ghToken == "" {
		logger.Error("No Github Token found. That really hoses the whole thing as you get 60 requests an hour. Set GITHUB_TOKEN")
		return
	}

	ctx := context.TODO()

	logger.Trace("Setting up client")
	// httpClient := &http.Client{
	// 	Timeout: time.Duration(20) * time.Second,
	// }

	print := func(context *github_ratelimit.CallbackContext) {
		logger.Warn(fmt.Sprintf("Secondary rate limit reached! Sleeping for %.2f seconds [%v --> %v]",
			time.Until(*context.SleepUntil).Seconds(), time.Now(), *context.SleepUntil))
	}

	rateLimiter, err := github_ratelimit.NewRateLimitWaiterClient(nil, github_ratelimit.WithLimitDetectedCallback(print))
	if err != nil {
		panic(err)
	}

	g := github.NewClient(rateLimiter).WithAuthToken(ghToken)

	for i := startPage; i <= maxPages; i++ {

		opt := &github.IssueListByRepoOptions{
			State:       "all",
			ListOptions: github.ListOptions{Page: i, PerPage: 10},
			Sort:        "created",
			Direction:   "asc",
		}
		issues, _, err := g.Issues.ListByRepo(ctx, donorDetails[0], donorDetails[1], opt)
		if err != nil {
			logger.Error(err.Error())
		}

		logger.Info("Issues scraped", logger.Args("count", len(issues), "issuePage", i))
		for idx, issue := range issues {
			if idx%5 == 0 && idx != 0 {
				logger.Info("uploaded ", logger.Args("idx", idx, "issuePage", i))
				time.Sleep(2 * time.Second)
			}

			logger.Trace("Issue found", logger.Args("issue", fmt.Sprintf("%v", issue)))
			candidateBody := fmt.Sprintf("Issue has been moved to https://github.com/k8sgateway/k8sgateway/issues/%v", issue.GetNumber())
			// candidateBody := fmt.Sprintf("Issue has been moved to https://github.com/%s/%s/issues/%v", donorDetails[0], donorDetails[1], issue.GetNumber())
			newTitle := fmt.Sprintf("[Migrated] %s", *issue.Title)
			var err error

			_, _, err = g.Issues.Create(ctx, destDetails[0], destDetails[1], &github.IssueRequest{
				Title: &newTitle,
				Body:  &candidateBody,
				State: issue.State,
			},
			)
			if err != nil {
				logger.Error(err.Error())
			}

			// place := "placeholder"
			// _, _, err := g.Issues.Create(ctx, destDetails[0], destDetails[1], &github.IssueRequest{
			// 	Title: &place,
			// },
			// )
			// if err != nil {
			// 	logger.Error(err.Error())
			// }

			// Pause between mutative requests. If you are making a large number of POST, PATCH, PUT, or DELETE requests, wait at least one second between each request.
			time.Sleep(1 * time.Second)
			// _, _, err = g.Issues.Edit(ctx, destDetails[0], destDetails[1], issue.GetNumber(), &github.IssueRequest{
			// 	Title: &newTitle,
			// 	Body:  &candidateBody,
			// 	State: issue.State,
			// })
			// if err != nil {
			// 	logger.Error(err.Error())
			// }
			// Pause between mutative requests. If you are making a large number of POST, PATCH, PUT, or DELETE requests, wait at least one second between each request.
			time.Sleep(1 * time.Second)

		}

	}
	logger.Info("Issue grabber completed")
}
