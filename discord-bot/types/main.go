package types

type CreateMSG struct {
	Title     string
	MSG       string
	Author    string
	CreatedAt string
}

type Gamer struct {
	User   string
	Budget int
	Win    int
}
