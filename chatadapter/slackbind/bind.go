package slackbind

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/astronoka/glados"
	"github.com/nlopes/slack"
)

// NewChatAdapter is create slack chatadapter implement
func NewChatAdapter(c glados.Context) glados.ChatAdapter {
	slackClient := slack.New(c.Env("GLADOS_SLACK_BOT_UAER_TOKEN", ""))
	rtm := slackClient.NewRTM()
	adapter := &slackChatAdapter{
		users:    make(map[string]*slack.User),
		channels: make(map[string]*slack.Channel),
		context:  c,
		client:   slackClient,
		rtm:      rtm,
	}
	go rtm.ManageConnection()
	go adapter.handleRTMEvent()
	return adapter
}

var slackUserIDPattern = regexp.MustCompile(`<@([a-zA-Z0-9_-]+)>`)

type hereHandler struct {
	regexp *regexp.Regexp
	handle glados.ChatBotMessageHandler
}

type slackChatAdapter struct {
	mu           sync.Mutex
	users        map[string]*slack.User
	channels     map[string]*slack.Channel
	context      glados.Context
	client       *slack.Client
	rtm          *slack.RTM
	hereHandlers []hereHandler
}

func (s *slackChatAdapter) handleRTMEvent() {
	for s.recieveEventLoop() {
	}
}

func (s *slackChatAdapter) recieveEventLoop() bool {
	logger := s.context.Logger()
	rtm := s.rtm
	select {
	case msg, ok := <-rtm.IncomingEvents:
		if !ok {
			logger.Debugln("slackbind: Event channel closed")
			return false
		}
		logger.Debugln("slackbind: Event Received: " + msg.Type)
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello
		case *slack.MessageEvent:
			s.onMessageEvent(ev)
		case *slack.PresenceChangeEvent:
			logger.Debugf("slackbind: Presence Change: %v\n", ev)
		case *slack.LatencyReport:
			logger.Debugf("slackbind: Current latency: %v\n", ev.Value)
		case *slack.RTMError:
			logger.Debugf("slackbind: Error: %s\n", ev.Error())
		case *slack.InvalidAuthEvent:
			logger.Debugf("slackbind: Invalid credentials")
			return false
		default:
			// Ignore other events..
			//b, _ := json.Marshal(msg.Data)
			//logger.Debugf("slackbind: Unsupported: %s - %s\n", msg.Type, string(b))
		}
		//case <-time.After(time.Second * 1):
		//	logger.Debugf("slackbind: timeout")
	}
	return true
}

func (s *slackChatAdapter) PostTextMessage(channel, text string) {
	params := slack.PostMessageParameters{
		Username:  s.context.BotName(),
		AsUser:    true,
		LinkNames: 1,
	}
	_, _, err := s.client.PostMessage(channel, text, params)
	if err != nil {
		s.context.Logger().Warnln("slackbind: post text message error. " + err.Error())
	}
}

func (s *slackChatAdapter) PostMessage(channel string, message *glados.ChatMessage) {
	attachment := slack.Attachment{
		AuthorName:    message.Author.Name,
		AuthorSubname: message.Author.Subname,
		AuthorLink:    message.Author.Link,
		AuthorIcon:    message.Author.IconURL,
		Title:         message.Title,
		TitleLink:     message.TitleLinkURL,
		Text:          message.Text,
		Color:         message.Color,
	}
	params := slack.PostMessageParameters{
		Username:    s.context.BotName(),
		AsUser:      true,
		LinkNames:   1,
		Attachments: []slack.Attachment{attachment},
	}
	_, _, err := s.client.PostMessage(channel, "", params)
	if err != nil {
		s.context.Logger().Warnln("slackbind: post message error. " + err.Error())
	}
}

func (s *slackChatAdapter) Here(pattern string, handler glados.ChatBotMessageHandler) {
	s.hereHandlers = append(s.hereHandlers, hereHandler{
		regexp: regexp.MustCompile(pattern),
		handle: handler,
	})
}

func (s *slackChatAdapter) Respond(pattern string, handler glados.ChatBotMessageHandler) {
	respondRegexpPattern := fmt.Sprintf(`^(?:@?(?:%s|%s)[:,]?)\s+(?:%s)`,
		s.context.BotName(), s.context.BotNameAlias(), pattern)
	s.hereHandlers = append(s.hereHandlers, hereHandler{
		regexp: regexp.MustCompile(respondRegexpPattern),
		handle: handler,
	})
}

func (s *slackChatAdapter) onMessageEvent(event *slack.MessageEvent) {
	if isBotMessage(event) || event.Hidden {
		return
	}
	text := s.convertSlackUserID2Name(event.Text)
	for _, handler := range s.hereHandlers {
		matches := handler.regexp.FindAllStringSubmatch(text, -1)
		if len(matches) <= 0 {
			continue
		}
		handler.handle(s, &glados.ChatMessageEvent{
			Channel: s.getChannelName(event.Channel),
			User:    s.getUserName(event.User),
			Text:    text,
			Matches: matches,
		})
	}
}

func (s *slackChatAdapter) getChannelName(channelID string) string {
	if channel, exist := s.channels[channelID]; exist {
		return channel.Name
	}
	channel, err := s.client.GetChannelInfo(channelID)
	if err != nil {
		s.context.Logger().Warnln("slackbind: get channel info failed. " + err.Error())
		return channelID
	}
	s.channels[channelID] = channel
	return channel.Name
}

func (s *slackChatAdapter) getUserName(userID string) string {
	if user, exist := s.users[userID]; exist {
		return user.Name
	}
	user, err := s.client.GetUserInfo(userID)
	if err != nil {
		s.context.Logger().Warnln("slackbind: get user info failed. " + err.Error())
		return userID
	}
	s.users[userID] = user
	return user.Name
}

func (s *slackChatAdapter) convertSlackUserID2Name(text string) string {
	var oldNew []string
	uniqMap := map[string]struct{}{}
	matched := slackUserIDPattern.FindAllStringSubmatch(text, -1)
	for _, m := range matched {
		userID := m[1]
		if _, exist := uniqMap[userID]; exist {
			continue
		}
		userName := s.getUserName(userID)
		if userName == "" {
			continue
		}
		oldNew = append(oldNew, "<@"+userID+">", "@"+userName)
		uniqMap[userID] = struct{}{}
	}
	if len(oldNew) <= 0 {
		return text
	}
	r := strings.NewReplacer(oldNew...)
	return r.Replace(text)
}

func isBotMessage(event *slack.MessageEvent) bool {
	return event.BotID != ""
}
