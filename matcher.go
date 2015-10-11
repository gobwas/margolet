package telegram

type Match map[string]interface{}

type Matcher interface {
	Match(text string) (*Match, bool)
}
