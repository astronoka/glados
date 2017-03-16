package githubnotifier

import (
	"crypto/rand"
	"encoding/binary"
	"net/http"
	"strconv"

	"github.com/astronoka/glados"
	"github.com/google/go-github/github"
)

// GitHubEventConverter is translate event to message
type GitHubEventConverter interface {
	ConvertEventToChatMessage(event interface{}) *glados.ChatMessage
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

func buildEventFromRequest(context glados.Context, rc glados.RequestContext, secret string) (interface{}, int, string) {
	payload, err := github.ValidatePayload(rc.Request(), []byte(secret))
	if err != nil {
		return nil, http.StatusBadRequest, "githubnotifier: validate payload failed:" + err.Error()
	}

	eventType := github.WebHookType(rc.Request())
	event, err := github.ParseWebHook(eventType, payload)
	if err != nil {
		context.Logger().Infoln("githubnotifier: " + err.Error())
		return nil, http.StatusAccepted, "githubnotifier: unsupported event"
	}

	if eventType == "ping" {
		return nil, http.StatusOK, "ok"
	}
	return event, http.StatusOK, "ok"
}

func notifyEventToChatAdapter(context glados.Context, destination string, event interface{}, converter GitHubEventConverter) {
	message := converter.ConvertEventToChatMessage(event)
	if message != nil {
		context.ChatAdapter().PostMessage(destination, message)
	}
}

func random() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, 36)
}
