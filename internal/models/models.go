package models

type AvaOfUser struct {
	Login   string
	Picture []byte
}

type StoriesResponse struct {
	Id         int
	Story      []byte
	Like       int
	Dislike    int
	EmojiHaha  int
	EmojiZzz   int
	EmojiWho   int
	EmojiRvota int
	Text       string
	Y          string
	X          string
	Color      string
	Style      string
}
