package text

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"testing"
)

func TestMessage(t *testing.T) {
	p := message.NewPrinter(language.BritishEnglish)
	printf, err := p.Printf("There are %v flowers in our garden.\n", 1500)
	t.Log(printf, err)
}
