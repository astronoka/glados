package glados

// ChatBotMessageHandler is chat bot callback function
type ChatBotMessageHandler func(adapter ChatAdapter, message *ChatMessageEvent)

// ChatAdapter is glados chat interface
type ChatAdapter interface {
	PostTextMessage(channel, text string)
	PostMessage(channel string, message *ChatMessage)

	Here(pattern string, handler ChatBotMessageHandler)
	Respond(pattern string, handler ChatBotMessageHandler)
}

// MessageAuthor is message author
type MessageAuthor struct {
	Name    string
	Subname string
	Link    string
	IconURL string
}

// ChatMessage is chat message
type ChatMessage struct {
	Author       MessageAuthor
	Title        string
	TitleLinkURL string
	Text         string
	Color        string
}

// ChatMessageEvent is message from chat system
type ChatMessageEvent struct {
	Channel string
	User    string
	Text    string
	Matches [][]string
}
