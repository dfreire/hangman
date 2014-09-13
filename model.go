package hangman

type Game struct {
	Id          string `sql:"type:varchar(100);unique"`
	AppId       string
	Theme       string
	Clue        string
	Answer      string
	Url         string
	Approved    bool
	FlagCount   int
	LikeCount   int
	UnlikeCount int
}
