package quiz

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
	QUIZ_PARAM_Title = "title"
	QUIZ_PARAM_List  = "quiz_list"
)

type Quiz struct {
	Category        int
	parameters      string
	includeSolution bool
}

type QuizDocument struct {
	Title   string
	Quizzes []Quiz
	// debug, make private/lowercase...
	QueryString string
}

func NewQuizDocument() *QuizDocument {
	return &QuizDocument{}
}

func (d *QuizDocument) AddQuiz(category int, parameters string, includeSolution bool) {
	q := Quiz{Category: category, parameters: parameters, includeSolution: includeSolution}
	d.Quizzes = append(d.Quizzes, q)
}
