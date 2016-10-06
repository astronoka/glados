package githubnotifier

import "github.com/astronoka/glados"

func sayPong(adapter glados.ChatAdapter, message *glados.ChatMessageEvent) {
	adapter.PostTextMessage(message.Channel, "@"+message.User+" pong")
}
