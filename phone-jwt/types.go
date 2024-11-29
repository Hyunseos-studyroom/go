package main

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
}

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}

type Todo struct {
	UserID uint64 `json:"user_id"`
	Title  string `json:"title"`
}
