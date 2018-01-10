package models

import (
	"time"
)

type News struct {
	Id          int
	TimeCreated time.Time
	Header      string
	Text        string
	Image       string // path to image
}

// database mock
var (
	news = []News{
		News{
			Id:          1,
			TimeCreated: time.Now(),
			Header:      "naslov 1",
			Text:        "ovo je text ultra kul stvari",
			Image:       "/dev/null",
		},
		News{
			Id:          2,
			TimeCreated: time.Now(),
			Header:      "naslov 2",
			Text:        "ovo je text ultra kul 2",
			Image:       "/dev/null1",
		},
	}
	nextNewsId = 3
)

func GetAllNews() []News {
	return news
}
