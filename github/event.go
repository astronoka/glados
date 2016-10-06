package github

import (
	"encoding/json"
	"errors"
)

// Event is git hub repository event interface
type Event interface {
	Type() string
}

// PullRequestEvent is pull request event object
// https://developer.github.com/v3/activity/events/types/#pullrequestevent
type PullRequestEvent struct {
	Action      string         `json:"action"`
	Number      int64          `json:"number"`
	PullRequest GitPullRequest `json:"pull_request"`
	Repository  GitRepository  `json:"repository"`
}

const (
	// OPENED is one of pull request status
	OPENED = "opened"
	// REOPENED is one of pull request status
	REOPENED = "reopened"
	// CLOSED is one of pull request status
	CLOSED = "closed"
	// MERGED is one of pull request status
	MERGED = "merged"
)

// Type is return event type string
func (*PullRequestEvent) Type() string {
	return "pull_request"
}

// Status is detect pull request status
func (e *PullRequestEvent) Status() string {
	switch e.Action {
	case OPENED, REOPENED:
		return e.Action
	case CLOSED:
		if e.PullRequest.MergeCommitSha != "" {
			return MERGED
		}
		return CLOSED
	default:
		return ""
	}
}

// IssueCommentEvent is issue comment event object
// https://developer.github.com/v3/activity/events/types/#issuecommentevent
type IssueCommentEvent struct {
	Action     string        `json:"action"`
	Issue      GitIssue      `json:"issue"`
	Comment    GitComment    `json:"comment"`
	Repository GitRepository `json:"repository"`
	Sender     GitUser       `json:"sender"`
}

// Type is return event type string
func (*IssueCommentEvent) Type() string {
	return "issue_comment"
}

// PullRequestReviewCommentEvent is pull request review comment event object
// https://developer.github.com/v3/activity/events/types/#pullrequestreviewcommentevent
type PullRequestReviewCommentEvent struct {
	Action      string         `json:"action"`
	Comment     GitComment     `json:"comment"`
	PullRequest GitPullRequest `json:"pull_request"`
	Repository  GitRepository  `json:"repository"`
	Sender      GitUser        `json:"sender"`
}

// Type is return event type string
func (*PullRequestReviewCommentEvent) Type() string {
	return "pull_request_review_comment"
}

// BuildEvent is create event instance from bytes
func BuildEvent(typ string, bytes []byte) (Event, error) {
	var event Event
	switch typ {
	case "pull_request":
		event = &PullRequestEvent{}
	case "issue_comment":
		event = &IssueCommentEvent{}
	case "pull_request_review_comment":
		event = &PullRequestReviewCommentEvent{}
	default:
	}
	if event != nil {
		err := json.Unmarshal(bytes, &event)
		if err != nil {
			return nil, errors.New("github: unmarshal event failed. " + err.Error())
		}
		return event, nil
	}
	return nil, errors.New("github: unsupported event type " + typ)
}
