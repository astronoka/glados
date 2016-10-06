package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/astronoka/glados"
	"github.com/astronoka/glados/chatadapter/slackbind"
	"github.com/astronoka/glados/program/githubnotifier"
	"github.com/astronoka/glados/router/ginbind"
	"github.com/astronoka/glados/storage/memory"
	"github.com/joho/godotenv"
)

func main() {
	logger := logrus.New()
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Level = logrus.DebugLevel

	err := godotenv.Load()
	if err != nil {
		logger.Errorln("Error loading .env file")
	}

	context := glados.BuildDefaultContextFromEnv()
	context.SetLogger(logger)
	context.SetStorage(memory.NewStorage(context))
	context.SetRouter(ginbind.NewRouter(context))
	context.SetChatAdapter(slackbind.NewChatAdapter(context))

	glados := glados.New(context)
	glados.Install(githubnotifier.NewProgram(map[string]string{
		"githubUserName": "slackUserName",
	}))
	glados.Boot()
}
