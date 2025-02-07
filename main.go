package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v69/github"
	"os"
	"regexp"
	"strconv"
)

func printInfo(notification *github.Notification, pr *github.PullRequest, prNumberStr string, url string) {

	fmt.Println("Number: " + prNumberStr)
	fmt.Println("Url: " + url)
	fmt.Println("Subject: " + *notification.Subject.Title)
	fmt.Println("Reason: " + *notification.Reason)
	fmt.Println("state: " + *pr.State)
}

func main() {
	client := github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))
	ctx := context.Background()
	opt := &github.NotificationListOptions{}
	notifications, response, err := client.Activity.ListNotifications(ctx, opt)
	if err != nil {
		fmt.Println("error " + err.Error() + " " + response.Status)
		os.Exit(1)
	}
	for _, notification := range notifications {
		url := *notification.Subject.URL
		urlBytes := []byte(url)
		matched, err := regexp.Match(`pulls/[0-9]*`, []byte(url))
		if err != nil {
			continue
		}

		issueMatched, err := regexp.Match(`issues/[0-9]*`, []byte(url))
		if err != nil {
			continue
		}
		if matched {
			handlePR(urlBytes, err, client, ctx, notification, url)
		}
		if issueMatched {
			handleIssue(urlBytes, err, client, ctx, notification, url)
		}

	}
}

func handleIssue(urlBytes []byte, err error, client *github.Client, ctx context.Context, notification *github.Notification, url string) {
	issuesNumRe := regexp.MustCompile(`issues/[0-9]+`)
	foundB := issuesNumRe.Find(urlBytes)
	numRe := regexp.MustCompile(`[0-9]+`)
	issueNumber := numRe.Find(foundB)
	issueNumberStr := string(issueNumber)
	issueNumberInt, err := strconv.Atoi(issueNumberStr)
	if err != nil {
		fmt.Println("error " + err.Error())
		os.Exit(1)
	}
	issue, _, err := client.Issues.Get(ctx, *notification.Repository.Owner.Login, *notification.Repository.Name, issueNumberInt)
	if err != nil {
		fmt.Println("error " + err.Error())
		os.Exit(1)
	}
	if *issue.State == "closed" {
		printIssueInfo(notification, issue, issueNumberStr, url)
		idInt, err := strconv.Atoi(*notification.ID)
		if err != nil {
			fmt.Println("error " + err.Error())
			os.Exit(1)
		}
		resp, err := client.Activity.MarkThreadDone(ctx, int64(idInt))
		if err != nil {
			fmt.Println("error " + err.Error() + " " + resp.Status)
			os.Exit(1)
		}
		fmt.Println("Marked as done")
		fmt.Println()
	}
}

func printIssueInfo(notification *github.Notification, issue *github.Issue, issueNumberStr string, url string) {
	fmt.Println("Number: " + issueNumberStr)
	fmt.Println("Url: " + url)
	fmt.Println("Subject: " + *notification.Subject.Title)
	fmt.Println("Reason: " + *notification.Reason)
	fmt.Println("state: " + *issue.Title)
}

func handlePR(urlBytes []byte, err error, client *github.Client, ctx context.Context, notification *github.Notification, url string) {
	pullsNumRe := regexp.MustCompile(`pulls/[0-9]+`)
	foundB := pullsNumRe.Find(urlBytes)
	numRe := regexp.MustCompile(`[0-9]+`)
	prNumber := numRe.Find(foundB)
	prNumberStr := string(prNumber)
	prNumberInt, err := strconv.Atoi(prNumberStr)
	if err != nil {
		fmt.Println("error " + err.Error())
		os.Exit(1)
	}

	pr, _, err := client.PullRequests.Get(ctx, *notification.Repository.Owner.Login, *notification.Repository.Name, prNumberInt)
	if err != nil {
		fmt.Println("error " + err.Error())
		os.Exit(1)
	}
	if *pr.State == "closed" {
		printInfo(notification, pr, prNumberStr, url)
		idInt, err := strconv.Atoi(*notification.ID)
		if err != nil {
			fmt.Println("error " + err.Error())
			os.Exit(1)
		}
		resp, err := client.Activity.MarkThreadDone(ctx, int64(idInt))
		if err != nil {
			fmt.Println("error " + err.Error() + " " + resp.Status)
			os.Exit(1)
		}
		fmt.Println("Marked as done")
		fmt.Println()
	}
}
