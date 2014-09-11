package hangman

type Game struct {
	Id          string `sql:"type:varchar(100);unique"`
	AppId       string
	Theme       string
	Clue        string
	Answer      string
	Url         string
	AuthorId    string
	Approved    bool
	FlagCount   int
	LikeCount   int
	UnlikeCount int
}

type Author struct {
	Id    string
	AppId string
	Email string
	Games []string
}
