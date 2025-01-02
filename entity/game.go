package entity

type Game struct {
	ID           uint
	Category     string
	QuestionList []string
	Players      []uint
}
