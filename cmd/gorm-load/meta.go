package main

import (
	"html/template"
	"strings"

	"github.com/sniperkit/jargon"
	"github.com/sniperkit/jargon/stackexchange"
)

var (
	jargonLem = jargon.NewLemmatizer(stackexchange.Dictionary, 3)
	lemma     = template.Must(template.New("lemma").Parse(`<span class="lemma">{{ . }}</span>`))
	plain     = template.Must(template.New("plain").Parse(`{{ . }}`))
)

func ExtractStack(content string, slug string) (topics []string) {
	// text := `Let’s talk about Ruby on Rails and ASPNET MVC.`
	r := strings.NewReader(content)
	tokens := jargon.Tokenize(r)

	// Or! Pass tokens on to the lemmatizer
	lemmatized := jargonLem.Lemmatize(tokens)
	for {
		t := lemmatized.Next()
		if t == nil {
			break
		}
		if t.IsLemma() {
			// remove duplciate entries
			topics = append(topics, t.String())
			// fmt.Println("slug: ", slug, "lemma: ", t)
		} // else {
		//  fmt.Println("slug: ", slug, "plain: ", t)
		//}
	}
	return
}
