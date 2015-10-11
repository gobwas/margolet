package telegram

import "regexp"

type Match map[string]interface{}

type Matcher interface {
	Match(text string) (*Match, bool)
}

type Equal struct {
	pattern string
}

func (self Equal) Match(text string) (match *Match, ok bool) {
	if self.pattern == text {
		match = &Match{"text": text}
		ok = true
	}

	return
}

// todo
//type RegExp struct {
//	pattern regexp.Regexp
//}
//
//func (self RegExp) Match(text string) (match *Match, ok bool) {
//	if self.pattern.Match([]byte(text)) {
//		match = &Match{"text": text}
//		ok = true
//	}
//
//	return
//}
