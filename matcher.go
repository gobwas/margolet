package telegram

import "regexp"

type DataKey byte

const (
	SUB_MATCH DataKey = iota
)

type Slug struct {
	Key, Value string
}

type Match struct {
	Text  string
	Slugs []Slug
	Data  map[interface{}]interface{}
}

type Matcher interface {
	Match(text string) (*Match, bool)
}

type Equal struct {
	pattern string
}

func (self Equal) Match(text string) (match *Match, ok bool) {
	if self.pattern == text {
		match = &Match{Text: text}
		ok = true
	}

	return
}

type RegExp struct {
	pattern *regexp.Regexp
}

func (self RegExp) Match(text string) (match *Match, ok bool) {
	if sub := self.pattern.FindStringSubmatch(text); len(sub) == 1+self.pattern.NumSubexp() {
		match = &Match{Text: text, Data: make(map[interface{}]interface{})}
		match.Data[SUB_MATCH] = sub
		ok = true
	}

	return
}
