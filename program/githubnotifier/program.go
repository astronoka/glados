package githubnotifier

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/astronoka/glados"
	"github.com/google/go-github/github"
)

// GitHubUserNamePattern is regexp for detect github user name
var GitHubUserNamePattern = regexp.MustCompile(`@([a-zA-Z0-9_-]+)`)

// NewProgram is create test program
func NewProgram(nameTable map[string]string) glados.Program {
	return &program{
		nameTable: nameTable,
	}
}

type program struct {
	nameTable map[string]string
}

func (p *program) Initialize(c glados.Context) {
	secret := c.Env("GLADOS_GITHUB_NOTIFIER_SECRET", random())
	// destination -> channel_name
	c.Router().POST("/github/notify_events/:destination", NotifyEvent(c, p, secret))
	c.ChatAdapter().Respond(`(?i)ping$`, sayPong)
}

func sayPong(adapter glados.ChatAdapter, message *glados.ChatMessageEvent) {
	adapter.PostTextMessage(message.Channel, "@"+message.User+" pong")
}

func (p *program) ConvertGithubName2ChatName(githubName string) string {
	if name, exist := p.nameTable[githubName]; exist {
		return name
	}
	return githubName
}

func (p *program) ConvertGithubName2ChatNameInText(text string) string {
	var oldNew []string
	matched := GitHubUserNamePattern.FindAllStringSubmatch(text, -1)
	for _, m := range matched {
		githubName := m[1]
		chatName := p.ConvertGithubName2ChatName(githubName)
		if chatName == "" {
			continue
		}
		oldNew = append(oldNew, githubName, chatName)
	}
	if len(oldNew) <= 0 {
		return text
	}
	r := strings.NewReplacer(oldNew...)
	return r.Replace(text)
}

func (p *program) ConvertEventToChatMessage(event interface{}) *glados.ChatMessage {
	switch event := event.(type) {
	case *github.PullRequestEvent:
		return p.buildPullRequestEventMessage(event)
	case *github.IssueCommentEvent:
		return p.buildIssueCommentEventMessage(event)
	case *github.PullRequestReviewCommentEvent:
		return p.buildPullRequestReviewCommentEventMessage(event)
	}
	return nil
}

func (p *program) buildPullRequestEventMessage(event *github.PullRequestEvent) *glados.ChatMessage {
	author := glados.MessageAuthor{
		Name:    *event.PullRequest.User.Login,
		Subname: p.ConvertGithubName2ChatName(*event.PullRequest.User.Login),
		Link:    *event.PullRequest.User.HTMLURL,
		IconURL: *event.PullRequest.User.AvatarURL,
	}
	text := p.ConvertGithubName2ChatNameInText(*event.PullRequest.Body) +
		"\n" +
		fmt.Sprintf(`<%s|github>`, event.PullRequest.HTMLURL)
	return &glados.ChatMessage{
		Author: author,
		Title:  fmt.Sprintf("%s: pull reques %s", *event.Repo.FullName, *event.Action),
		Text:   text,
		Color:  "#000000",
	}
}

func (p *program) buildIssueCommentEventMessage(event *github.IssueCommentEvent) *glados.ChatMessage {
	if *event.Action != "created" {
		return nil
	}
	text := "issue owner: @" + p.ConvertGithubName2ChatName(*event.Issue.User.Login) +
		"\n" +
		p.ConvertGithubName2ChatNameInText(*event.Comment.Body) +
		"\n" +
		fmt.Sprintf(`<%s|github>`, event.Comment.HTMLURL)
	author := glados.MessageAuthor{
		Name:    *event.Sender.Login,
		Subname: p.ConvertGithubName2ChatName(*event.Sender.Login),
		Link:    *event.Sender.HTMLURL,
		IconURL: *event.Sender.AvatarURL,
	}
	return &glados.ChatMessage{
		Author: author,
		Title:  fmt.Sprintf("%s: issue comment created", *event.Repo.FullName),
		Text:   text,
		Color:  "#000000",
	}
}

func (p *program) buildPullRequestReviewCommentEventMessage(event *github.PullRequestReviewCommentEvent) *glados.ChatMessage {
	if *event.Action != "created" {
		return nil
	}
	author := glados.MessageAuthor{
		Name:    *event.Sender.Login,
		Subname: p.ConvertGithubName2ChatName(*event.Sender.Login),
		Link:    *event.Sender.HTMLURL,
		IconURL: *event.Sender.AvatarURL,
	}
	text := "pull request owner: @" + p.ConvertGithubName2ChatName(*event.PullRequest.User.Login) +
		"\n" +
		p.ConvertGithubName2ChatNameInText(*event.Comment.Body) +
		"\n" +
		fmt.Sprintf(`<%s|github>`, event.Comment.HTMLURL)
	return &glados.ChatMessage{
		Author: author,
		Title:  fmt.Sprintf("%s: review comment created", *event.Repo.FullName),
		Text:   text,
		Color:  "#000000",
	}
}
