package glados

import (
	"os"
	"sync"
)

// Context is glados system context
type Context interface {
	BotName() string
	BotNameAlias() string
	SetLogger(Logger)
	SetStorage(Storage)
	SetRouter(Router)
	SetChatAdapter(ChatAdapter)
	Logger() Logger
	Storage() Storage
	Router() Router
	ChatAdapter() ChatAdapter
	ListenPort() string
	Env(string, string) string
}

type contextImpl struct {
	mu           sync.Mutex
	botName      string
	botNameAlias string
	listenPort   string
	logger       Logger
	storage      Storage
	router       Router
	chatadapter  ChatAdapter
}

func (c *contextImpl) BotName() string {
	return c.botName
}

func (c *contextImpl) BotNameAlias() string {
	return c.botNameAlias
}

func (c *contextImpl) SetLogger(l Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger = l
}

func (c *contextImpl) SetStorage(s Storage) {
	c.storage = s
}

func (c *contextImpl) SetRouter(r Router) {
	c.router = r
}

func (c *contextImpl) SetChatAdapter(a ChatAdapter) {
	c.chatadapter = a
}

func (c *contextImpl) Logger() Logger {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.logger
}

func (c *contextImpl) Storage() Storage {
	return c.storage
}

func (c *contextImpl) Router() Router {
	return c.router
}

func (c *contextImpl) ChatAdapter() ChatAdapter {
	return c.chatadapter
}

func (c *contextImpl) ListenPort() string {
	return c.listenPort
}

func (c *contextImpl) Env(key, valueIfNotFound string) string {
	return env(key, valueIfNotFound)
}

// BuildDefaultContextFromEnv is create default contextImpl instance
func BuildDefaultContextFromEnv() Context {
	botName := env("BOT_NAME", "GLaDOS")
	return &contextImpl{
		botName:      botName,
		botNameAlias: env("BOT_NAME_ALIAS", botName),
		listenPort:   env("PORT", "7000"),
	}
}

func env(key, valueIfNotFound string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return valueIfNotFound
}
