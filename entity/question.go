package entity

type Question struct {
	ID            uint
	Question      string
	AnswerList    []string
	CorrectAnswer string
	Difficulty    string
	Category      string
}
