package telegram

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
