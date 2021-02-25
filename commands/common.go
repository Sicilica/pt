package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/sicilica/pt/types"
)

// Get returns the command with the given name.
func Get(name string) (func(c types.CommandContext) error, bool) {
	fn, ok := commandFns[name]
	return fn, ok
}

func register(example string, fn func(c types.CommandContext) error, category, desc string) {
	name := strings.SplitN(example, " ", 2)[0]
	commandFns[name] = fn
	commandExamples[name] = example
	commandDescriptions[name] = desc
	categoryCommands[category] = append(categoryCommands[category], name)
}

var commandFns = map[string]func(c types.CommandContext) error{}
var commandExamples = map[string]string{}
var commandDescriptions = map[string]string{}
var categoryCommands = map[string][]string{}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	h := int(d.Hours()) - days*24
	m := int(d.Minutes()) - (days*24+h)*60
	s := int(d.Seconds()) - ((days*24+h)*60+m)*60

	if days >= 2 {
		return fmt.Sprintf("%dd", days)
	}
	if days >= 1 {
		return fmt.Sprintf("%dd%dh", days, h)
	}
	if h >= 1 {
		return fmt.Sprintf("%dh%sm", h, pad2(m))
	}
	if m >= 1 {
		return fmt.Sprintf("%dm%ss", m, pad2(s))
	}
	return fmt.Sprintf("%ds", s)
}

func formatTags(tags []string) string {
	return fmt.Sprintf("[%s]", strings.Join(tags, " "))
}

func formatTime(t time.Time) string {
	t = t.Local()
	now := time.Now().Local()

	tYear, tMonth, tDay := t.Date()
	nowYear, nowMonth, nowDay := now.Date()
	isSameDay := tYear == nowYear && tMonth == nowMonth && tDay == nowDay

	if isSameDay {
		return t.Format("15:04")
	}

	return t.Format("Jan 2 15:04")
}

func pad2(i int) string {
	if i < 10 {
		return fmt.Sprintf("0%d", i)
	}
	return fmt.Sprint(i)
}
