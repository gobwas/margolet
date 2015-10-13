package telegram

import "regexp"

type Slug struct {
	Key, Value string
}

type Match struct {
	Text  string
	Slugs []Slug
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

type RegExpMatch struct {
	Match
	Submatch []string
}

type RegExp struct {
	pattern *regexp.Regexp
}

func (self RegExp) Match(text string) (match *Match, ok bool) {
	if self.pattern.MatchString(text) {
		match = &Match{Text: text}
		for _, sub := range self.pattern.FindStringSubmatch(text)[1:] {
			match.Slugs = append(match.Slugs, Slug{Value: sub})
		}
		ok = true
	}

	return
}
