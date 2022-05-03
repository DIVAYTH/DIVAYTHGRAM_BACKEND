package models

import (
	"database/sql"
	"sync"
)

var ai autoInc

type User struct {
	Login    string `gorm:"primary_key"`
	Password string `gorm:"not null"`
}

type UserAva struct {
	Login string `gorm:"primary_key"`
	Ava   string `gorm:"not null"`
}

type Stories struct {
	Id         int    `gorm:"primary_key"`
	Login      string `gorm:"not null"`
	Story      string `gorm:"not null"`
	Like       int    `gorm:"not null"`
	Dislike    int    `gorm:"not null"`
	EmojiHaha  int    `gorm:"not null"`
	EmojiZzz   int    `gorm:"not null"`
	EmojiWho   int    `gorm:"not null"`
	EmojiRvota int    `gorm:"not null"`
	Text       sql.NullString
	Y          sql.NullString
	X          sql.NullString
	Color      sql.NullString
	Style      sql.NullString
}

type StoriesLikeDislike struct {
	Id      int    `gorm:"primary_key"`
	Login   string `gorm:"primary_key"`
	Like    bool   `gorm:"not null"`
	Dislike bool   `gorm:"not null"`
}

type StoriesEmoji struct {
	Id         int    `gorm:"primary_key"`
	Login      string `gorm:"primary_key"`
	EmojiHaha  bool   `gorm:"not null"`
	EmojiZzz   bool   `gorm:"not null"`
	EmojiWho   bool   `gorm:"not null"`
	EmojiRvota bool   `gorm:"not null"`
}

func NewStories() *Stories {
	return &Stories{
		Id: ai.id,
	}
}

type autoInc struct {
	sync.Mutex
	id int
}

func (a *autoInc) ID() (id int) {
	a.Lock()
	defer a.Unlock()
	id = a.id
	a.id++
	return
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
