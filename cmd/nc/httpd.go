package main

// The webserver portion relies on resources.
// I rely on this hack as I don't know how to to it correctly using embed.
//go:generate esc -prefix ../../output/templates -pkg output -o ../../output/assets.go -private ../../output/templates

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"schnoddelbotz/nc/output"
	"schnoddelbotz/nc/quiz"
)

type httpdConfig struct {
	address string
	enabled bool
}

var (
	config    httpdConfig
	indexHTML string
)

func thisOrThat(this, that string) string {
	if this == "" {
		return that
	}
	return this
}

func RunWebserver(config httpdConfig) {
	if !config.enabled {
		return
	}
	log.Printf("Webeserver enabled %v, using address %s\n", *serveHTTP, config.address)

	http.HandleFunc("/", indexHandler)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(output.MustFs())))

	err := http.ListenAndServe(config.address, nil)
	if err != nil {
		log.Fatalf("could not start listening or address %s %s", config.address, err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	quizDoc := urlParamsToQuizDocument(r)
	w.Write(renderIndexTemplate(*quizDoc))
}

func renderIndexTemplate(quizDocument quiz.QuizDocument) []byte {
	buf := &bytes.Buffer{}

	indexHTML := output.MustGetTemplate("/html/index.html")
	tpl, err := template.New("index").Parse(indexHTML)
	if err != nil {
		log.Fatalf("Template parsing error: %v\n", err)
	}
	err = tpl.Execute(buf, quizDocument)
	if err != nil {
		log.Printf("Template execution error: %v\n", err)
	}
	return buf.Bytes()
}

func urlParamsToQuizDocument(r *http.Request) *quiz.QuizDocument {
	testQuiz := quiz.ExampleQuiz1() // DEV WIP -> test. tdd?
	testQuiz.Title = r.URL.Query().Get(quiz.QUIZ_PARAM_Title)
	testQuiz.QueryString = r.URL.RawQuery

	quizList := r.URL.Query().Get("quizzes")
	if quizList != "" {
		fmt.Printf("TODO ... built quiz for %s\n", quizList)
		//append Quizzes:     r.URL.Query().Get(quiz.QUIZ_PARAM_List),
	}

	return testQuiz
}
