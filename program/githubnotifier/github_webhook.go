package githubnotifier

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/astronoka/glados"
	"github.com/astronoka/glados/github"
)

func notifyEvent(p *program, c glados.Context, secret string) glados.RequestHandler {
	return func(rc glados.RequestContext) {
		status, message := notifyEventToChatAdapter(p, c, rc, secret)
		rc.JSON(status, glados.H{
			"message": message,
		})
	}
}

func notifyEventToChatAdapter(p *program, c glados.Context, rc glados.RequestContext, secret string) (int, string) {
	signature := strings.TrimSpace(rc.Header("X-Hub-Signature"))
	maxSize := int64(1024 * 1024 * 5)
	reader := io.LimitReader(rc.RequestBody(), maxSize+1)
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}
	if int64(len(body)) > maxSize {
		return http.StatusBadRequest, "githubnotifier: post body too large"
	}

	if !verifySignature(secret, signature, body) {
		return http.StatusUnauthorized, "githubnotifier: invalid signature"
	}

	eventType := rc.Header("X-GitHub-Event")
	if eventType == "" {
		return http.StatusBadRequest, "githubnotifier: event type required"
	}

	if eventType == "ping" {
		return http.StatusOK, "ok"
	}

	destination := rc.Param("destination")
	event, err := github.BuildEvent(eventType, body)
	if err != nil {
		c.Logger().Infoln("githubnotifier: " + err.Error())
		return http.StatusAccepted, "githubnotifier: " + err.Error()
	}
	message := buildNotifyMessage(p, c, event)
	if message != nil {
		c.ChatAdapter().PostMessage(destination, message)
	}
	return http.StatusOK, "ok"
}

func verifySignature(secret string, signature string, body []byte) bool {
	if signature == "" {
		return false
	}
	signatures := strings.Split(signature, "=")
	if len(signatures) != 2 {
		return false
	}
	hashType := signatures[0]
	hashString := signatures[1]
	hmac := hmac.New(newHash(hashType), []byte(secret))
	hmac.Write(body)
	hmacString := hex.EncodeToString(hmac.Sum(nil))
	return hashString == hmacString
}

func newHash(hashType string) func() hash.Hash {
	if hashType == "sha1" {
		return sha1.New
	}
	return sha256.New
}

func random() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, 36)
}

func buildNotifyMessage(p *program, c glados.Context, event github.Event) *glados.ChatMessage {
	switch e := event.(type) {
	case *github.PullRequestEvent:
		return buildPullrequestEventMessage(p, e)
	case *github.IssueCommentEvent:
		return buildIssueCommentEventMessage(p, e)
	case *github.PullRequestReviewCommentEvent:
		return buildPullRequestReviewCommentEventMessage(p, e)
	}
	return &glados.ChatMessage{
		Title: "Unsupported",
		Text:  fmt.Sprintf("event %s not supported yet", event.Type()),
	}
}

func buildPullrequestEventMessage(p *program, event *github.PullRequestEvent) *glados.ChatMessage {
	author := glados.MessageAuthor{
		Name:    event.PullRequest.User.Login,
		Subname: p.convertGithubName2ChatName(event.PullRequest.User.Login),
		Link:    event.PullRequest.User.HTMLURL,
		IconURL: event.PullRequest.User.AvatarURL,
	}
	text := p.convertGithubName2ChatNameInText(event.PullRequest.Body) +
		"\n" +
		fmt.Sprintf(`<%s|github>`, event.PullRequest.HTMLURL)
	return &glados.ChatMessage{
		Author: author,
		Title:  fmt.Sprintf("%s: pull reques %s", event.Repository.FullName, event.Status()),
		Text:   text,
		Color:  "#000000",
	}
}

func buildIssueCommentEventMessage(p *program, event *github.IssueCommentEvent) *glados.ChatMessage {
	if event.Action != "created" {
		return nil
	}
	text := "issue owner: @" + p.convertGithubName2ChatName(event.Issue.User.Login) +
		"\n" +
		p.convertGithubName2ChatNameInText(event.Comment.Body) +
		"\n" +
		fmt.Sprintf(`<%s|github>`, event.Comment.HTMLURL)
	author := glados.MessageAuthor{
		Name:    event.Sender.Login,
		Subname: p.convertGithubName2ChatName(event.Sender.Login),
		Link:    event.Sender.HTMLURL,
		IconURL: event.Sender.AvatarURL,
	}
	return &glados.ChatMessage{
		Author: author,
		Title:  fmt.Sprintf("%s: issue comment created", event.Repository.FullName),
		Text:   text,
		Color:  "#000000",
	}
}

func buildPullRequestReviewCommentEventMessage(p *program, event *github.PullRequestReviewCommentEvent) *glados.ChatMessage {
	if event.Action != "created" {
		return nil
	}
	author := glados.MessageAuthor{
		Name:    event.Sender.Login,
		Subname: p.convertGithubName2ChatName(event.Sender.Login),
		Link:    event.Sender.HTMLURL,
		IconURL: event.Sender.AvatarURL,
	}
	text := "pull request owner: @" + p.convertGithubName2ChatName(event.PullRequest.User.Login) +
		"\n" +
		p.convertGithubName2ChatNameInText(event.Comment.Body) +
		"\n" +
		fmt.Sprintf(`<%s|github>`, event.Comment.HTMLURL)
	return &glados.ChatMessage{
		Author: author,
		Title:  fmt.Sprintf("%s: review comment created", event.Repository.FullName),
		Text:   text,
		Color:  "#000000",
	}
}
