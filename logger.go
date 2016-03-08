package telegram

import (
	"fmt"
	"log"
)

type Logger struct {
	prefix string
}

func (l Logger) Println(s ...interface{}) {
	log.Println(l.prefix + fmt.Sprintln(s...))
}
