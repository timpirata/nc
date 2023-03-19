package quiz

import (
	"fmt"
	"math/rand"
	urlpkg "net/url"
	"regexp"
	"strconv"
)

const (
	QUIZ_T_ADDITION int = iota
	QUIZ_T_SUBSTRACTION
	QUIZ_T_MULTIPLICATION
	QUIZ_T_DIVISION
	QUIZ_T_ROWS
	QUIZ_T_POWS
	QUIZ_T_SQRT
	QUIZ_T_NUMBER_RAY
	QUIZ_T_MATRIX
	QUIZ_T_FFT
	QUIZ_T_GRAPH
	QUIZ_T_IT
	QUIZ_T_NaturalLanguage
	QUIZ_T_ProgrammingLanguage
	QUIZ_T_History
	QUIZ_T_Person
	QUIZ_T_Fortune
	QUIZ_T_RANDOM
	QUIZ_PARAM_Title          = "title"
	QUIZ_PARAM_List           = "quiz_list"
	QUIZ_PARAM_Show           = "show"
	QUIZ_FORMFIELD_TaskType   = "tasktype"
	QUIZ_FORMFIELD_Count      = "count"
	QUIZ_FORMFIELD_Solutions  = "solutions"
	QUIZ_FORMFIELD_Parameters = "params"
	SolutionsIncluded         = true
	SolutionsExcluded         = false
)

type Quiz struct {
	Category        int
	parameters      string
	IncludeSolution bool
	HTML            string
	Solution        string
}

type QuizDocument struct {
	Title   string
	Quizzes []Quiz
	// debug, make private/lowercase...
	QueryString  string
	ShowForm     bool
	ShowQuiz     bool
	ShowExamples bool
}

func NewQuizDocument() *QuizDocument {
	return &QuizDocument{}
}

func NewQuiz(category int, parameters string, includeSolution bool) *Quiz {
	q := Quiz{Category: category, parameters: parameters, IncludeSolution: includeSolution}
	return &q
}

func (d *QuizDocument) AddQuiz(q *Quiz) {
	d.Quizzes = append(d.Quizzes, *q)
}

func (d *QuizDocument) BuildFromForm(formValues urlpkg.Values) {

	count := formValues[QUIZ_FORMFIELD_Count]
	types := formValues[QUIZ_FORMFIELD_TaskType]
	solutions := formValues[QUIZ_FORMFIELD_Solutions]
	params := formValues[QUIZ_FORMFIELD_Parameters]

	for index, cnt := range count {
		cntNum, err := strconv.Atoi(cnt)
		if err != nil {
			panic("Yikes")
		}
		typ := types[index]
		t, err := strconv.Atoi(typ)
		if err != nil {
			panic("Bad integer input")
		}
		solution := solutions[index]
		sol := false
		if solution == "Yes" {
			sol = true
		}
		param := params[index]
		fmt.Printf("idx %d cnt %s -> typ %d\n", index, cnt, t)
		for i := 0; i < cntNum; i++ {
			q := NewQuiz(t, param, sol)
			q.Produce()
			if solution == "yes" {
				q.IncludeSolution = true
			}
			d.AddQuiz(q)
		}
	}
}

func (quiz *Quiz) Produce() *Quiz {
	// generics? nicer? tbd ...
	switch quiz.Category {
	case QUIZ_T_ADDITION:
		return buildQuiz_Addition(quiz)
	case QUIZ_T_SUBSTRACTION:
		return buildQuiz_Substraction(quiz)
	}
	panic("Unkonw quiz type!")
	return nil // fmt.Sprintf("Unsupported QUIZ type/category %d", quiz.Category)
}

func buildQuiz_Addition(q *Quiz) *Quiz {
	// required parameters
	left_op_max := paramMustInt(q.parameters, "left_op_max")
	right_op_max := paramMustInt(q.parameters, "right_op_max")
	// defaults for optional parameters
	left_op_min := paramMayInt(q.parameters, "left_op_min", 0)
	right_op_min := paramMayInt(q.parameters, "right_op_min", 0)
	// randomise operands based on params
	left_op := rand.Intn(left_op_max-left_op_min) + left_op_min
	right_op := rand.Intn(right_op_max-right_op_min) + right_op_min
	/// after lunch
	q.HTML = fmt.Sprintf("%d + %d", left_op, right_op)
	// fmt.Printf("LOG %s\n", q.HTML)
	q.Solution = fmt.Sprintf("%d", left_op+right_op)
	return q
}

func buildQuiz_Substraction(q *Quiz) *Quiz {
	left_op_max := paramMustInt(q.parameters, "left_op_max")
	left_op_min := paramMayInt(q.parameters, "left_op_min", 0)
	left_op := rand.Intn(left_op_max-left_op_min) + left_op_min

	right_op_min := paramMayInt(q.parameters, "right_op_min", 0)
	right_op_max := 0
	right_op_max_v := paramMayString(q.parameters, "right_op_max")
	// right_op_max max point at (random) &left_op value.
	// otherwise, it must be provided as regular int
	if right_op_max_v == "&left_op" {
		right_op_max = left_op
	} else {
		right_op_max = paramMustInt(q.parameters, "right_op_max")
	}
	fmt.Printf("right_op_max_v: %s\n", right_op_max_v)
	right_op := rand.Intn(right_op_max-right_op_min) + right_op_min

	q.HTML = fmt.Sprintf("%d - %d", left_op, right_op)
	q.Solution = fmt.Sprintf("%d", left_op-right_op)
	return q
}

func paramMustInt(haystack string, needle string) int {
	r := regexp.MustCompile(needle + `=(\d+)`)
	matches := r.FindStringSubmatch(haystack)
	if len(matches) > 0 {
		fmt.Printf("MustIntMATCH %v\n", matches)
		v, err := strconv.Atoi(matches[1])
		if err == nil {
			return v
		}
	}
	panic("oops")
}

func paramMayInt(haystack string, needle string, defaultValue int) int {
	r := regexp.MustCompile(needle + `=(\d+)`)
	matches := r.FindStringSubmatch(haystack)
	if len(matches) > 0 {
		fmt.Printf("MayIntMATCH %v\n", matches)
		v, err := strconv.Atoi(matches[1])
		if err == nil {
			return v
		}
	}
	return defaultValue
}

func paramMayString(haystack string, needle string) string {
	r := regexp.MustCompile(needle + `=(\S+)`)
	matches := r.FindStringSubmatch(haystack)
	if len(matches) > 0 {
		fmt.Printf("MayStringMATCH(%s in %s) %v = %v\n", needle, haystack, r, matches)
		return matches[1]
	}
	return ""
}

func ExampleQuiz1() *QuizDocument {
	d := NewQuizDocument()
	d.Title = "Maths basic test, 3 random: +,  -,  *, /"
	// each q type => "random" params; parsed per quiz generator
	_addition_max_20 := "left_op_max=10, right_op_max=10"
	d.AddQuiz(NewQuiz(QUIZ_T_ADDITION, _addition_max_20, SolutionsIncluded))
	d.AddQuiz(NewQuiz(QUIZ_T_ADDITION, _addition_max_20, SolutionsExcluded))
	_substraction_positive_from_100 := "left_op_max=100, right_op_max=$left_op"
	d.AddQuiz(NewQuiz(QUIZ_T_SUBSTRACTION, _substraction_positive_from_100, SolutionsExcluded))
	return d
}
