package githubnotifier

import (
	"regexp"
	"strings"

	"github.com/astronoka/glados"
)

var githubUserNamePattern = regexp.MustCompile(`@([a-zA-Z0-9_-]+)`)

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
	c.Router().POST("/github/notify_events/:destination", notifyEvent(p, c, secret))
	c.ChatAdapter().Respond(`(?i)ping$`, sayPong)
}

func (p *program) convertGithubName2ChatName(githubName string) string {
	if name, exist := p.nameTable[githubName]; exist {
		return name
	}
	return "unknown"
}

func (p *program) convertGithubName2ChatNameInText(text string) string {
	var oldNew []string
	matched := githubUserNamePattern.FindAllStringSubmatch(text, -1)
	for _, m := range matched {
		githubName := m[1]
		chatName := p.convertGithubName2ChatName(githubName)
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
