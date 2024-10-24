// Helper script designed to send a cohesive Slack notification about the result of a collection of GitHub workflows
// This is a simpler variant of notify-from-json.go, due to some issues we
// encountered around the accuracy of the notification: https://github.com/solo-io/solo-projects/issues/5191
//
// This works by:
// 	1. Read in the test summary
//  2. Send a Slack notification with the contents
//
// Example usage:
// 	 export SLACKBOT_BEARER=${{ secrets.SLACKBOT_BEARER }}
// 	 export SLACK_CHANNEL=C0314KESVNV
// 	go run .github/workflows/helpers/notify/slack_test_summary.go ./_test/test_log/go-test-summary

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	postMessageEndpoint = "https://slack.com/api/chat.postMessage"
)

type Payload struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func main() {
	summaryFile := os.Args[1]
	fmt.Printf("slack_test_summary.go invoked with: %v", summaryFile)

	b, err := os.ReadFile(summaryFile)
	if err != nil {
		panic(err)
	}

	fails, err := findFailures(b)
	mustSendSlackText(string(b))
}

func findFailures(b []byte) ([]string, error) {

	var allReasons, failureReasons []string

	allReasons = strings.FieldsFunc(string(b), func(r rune) bool { return r == '\n' })

	for _, reason := range allReasons {
		if strings.HasPrefix(reason, "FAIL") {
			failureReasons = append(failureReasons, reason)
		}
	}

	// See https://api.slack.com/reference/surfaces/formatting for slack hyperlink formatting
	text := fmt.Sprintf(":red_circle: <$PARENT_JOB_URL|$PREAMBLE> have failed some jobs:\n%s", strings.Join(failureReasons, "\n"))
	mustSendSlackText(text)
}
func mustSendSlackText(text string) {
	fmt.Printf("send slack message with text: %s", text)
	mustSendSlackMessage(Payload{
		Channel: os.ExpandEnv("$SLACK_CHANNEL"),
		Text:    os.ExpandEnv(text),
	})
}

func mustSendSlackMessage(data Payload) {
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, postMessageEndpoint, bytes.NewReader(payloadBytes))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", os.ExpandEnv("Bearer $SLACKBOT_BEARER"))

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := netClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

}
