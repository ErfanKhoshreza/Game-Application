package entity

type Question struct {
	ID              uint
	Question        string
	PossibleAnswers []PossibleAnswer
	CorrectAnswerID uint
	Difficulty      QuestionDifficulty
	CategoryID      uint
}
type PossibleAnswer struct {
	ID      uint
	Content string
	Choice  PossibleAnswerChoice
}
type Answer struct {
	ID       uint
	PlayerID uint
	Question uint
}
type PossibleAnswerChoice uint8

func (p PossibleAnswerChoice) IsValid() bool {
	if p >= PossibleAnswerA && p <= PossibleAnswerD {
		return true
	}
	return false
}

const (
	PossibleAnswerA PossibleAnswerChoice = iota + 1
	PossibleAnswerB
	PossibleAnswerC
	PossibleAnswerD
)

type QuestionDifficulty string

const (
	easyDifficulty   QuestionDifficulty = "easy"
	mediumDifficulty QuestionDifficulty = "medium"
	hardDifficulty   QuestionDifficulty = "hard"
)

func (d QuestionDifficulty) isValid() bool {
	if d == easyDifficulty || d == mediumDifficulty || d == hardDifficulty {
		return true
	}
	return false
}
