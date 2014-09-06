package hangman

type Game struct {
	AppId       string
	Id          string
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
	AppId string
	Id    string
	Email string
	Games []string
}
