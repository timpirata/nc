package main

import (
	"bytes"
	// "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"schnoddelbotz/nc/quiz"
)

type httpdConfig struct {
	address string
}

var (
	config httpdConfig
	//goatnoet:embed ../../output/templates/html/index.html
	indexHTML string
)

func thisOrThat(this, that string) string {
	if this == "" {
		return that
	}
	return this
}

func RunWebserver(config httpdConfig) {
	http.HandleFunc("/", indexHandler)

	log.Printf("NC Listening for requests on %s", config.address)
	err := http.ListenAndServe(config.address, nil)

	if err != nil {
		log.Fatalf("could not start listening or address %s %s", config.address, err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	quizDoc := urlParamsToQuizDocument(r)
	log.Printf("built new quiz index doc, haha!")
	w.Write(renderIndexTemplate(*quizDoc))
}

func renderIndexTemplate(card quiz.QuizDocument) []byte {
	buf := &bytes.Buffer{}
	tpl, err := template.New("index").Parse(indexHTML)
	if err != nil {
		log.Fatalf("Template parsing error: %v\n", err)
	}
	err = tpl.Execute(buf, card)
	if err != nil {
		log.Printf("Template execution error: %v\n", err)
	}
	return buf.Bytes()
}

func urlParamsToQuizDocument(r *http.Request) *quiz.QuizDocument {
	doc := &quiz.QuizDocument{
		Title:       r.URL.Query().Get(quiz.QUIZ_PARAM_Title),
		QueryString: r.URL.RawQuery,
	}

	quizList := r.URL.Query().Get("quizzes")
	if quizList != "" {
		fmt.Printf("TODO ... built quiz for %s\n", quizList)
		//append Quizzes:     r.URL.Query().Get(quiz.QUIZ_PARAM_List),
	}

	return doc
}
