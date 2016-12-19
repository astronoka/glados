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

// GitHubEventConverter is translate event to message
type GitHubEventConverter interface {
	BuildPullRequestEventMessage(event *github.PullRequestEvent) *glados.ChatMessage
	BuildIssueCommentEventMessage(event *github.IssueCommentEvent) *glados.ChatMessage
	BuildPullRequestReviewCommentEventMessage(event *github.PullRequestReviewCommentEvent) *glados.ChatMessage
}

// NotifyEvent is create request handler
func NotifyEvent(context glados.Context, converter GitHubEventConverter, secret string) glados.RequestHandler {
	return func(rc glados.RequestContext) {
		event, status, message := buildEventFromRequest(context, rc, secret)
		if event == nil {
			rc.JSON(status, glados.H{
				"message": message,
			})
			return
		}

		destination := rc.Param("destination")
		notifyEventToChatAdapter(context, destination, event, converter)
		rc.JSON(http.StatusOK, glados.H{
			"message": "ok",
		})
	}
}

func buildEventFromRequest(context glados.Context, rc glados.RequestContext, secret string) (github.Event, int, string) {
	signature := strings.TrimSpace(rc.Header("X-Hub-Signature"))
	maxSize := int64(1024 * 1024 * 5)
	reader := io.LimitReader(rc.RequestBody(), maxSize+1)
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, http.StatusBadRequest, "githubnotifier: " + err.Error()
	}
	if int64(len(body)) > maxSize {
		return nil, http.StatusBadRequest, "githubnotifier: post body too large"
	}

	if !verifySignature(secret, signature, body) {
		return nil, http.StatusUnauthorized, "githubnotifier: invalid signature"
	}

	eventType := rc.Header("X-GitHub-Event")
	if eventType == "" {
		return nil, http.StatusBadRequest, "githubnotifier: event type required"
	}

	if eventType == "ping" {
		return nil, http.StatusOK, "ok"
	}

	event, err := github.BuildEvent(eventType, body)
	if err != nil {
		context.Logger().Infoln("githubnotifier: " + err.Error())
		return nil, http.StatusAccepted, "githubnotifier: unsupported event"
	}
	return event, http.StatusOK, "ok"
}

func notifyEventToChatAdapter(context glados.Context, destination string, event github.Event, converter GitHubEventConverter) {
	message := buildNotifyMessage(converter, event)
	if message != nil {
		context.ChatAdapter().PostMessage(destination, message)
	}
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

func buildNotifyMessage(converter GitHubEventConverter, event github.Event) *glados.ChatMessage {
	switch e := event.(type) {
	case *github.PullRequestEvent:
		return converter.BuildPullRequestEventMessage(e)
	case *github.IssueCommentEvent:
		return converter.BuildIssueCommentEventMessage(e)
	case *github.PullRequestReviewCommentEvent:
		return converter.BuildPullRequestReviewCommentEventMessage(e)
	}
	return &glados.ChatMessage{
		Title: "Unsupported",
		Text:  fmt.Sprintf("event %s not supported yet", event.Type()),
	}
}
