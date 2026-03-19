package db

import "time"

type Subscriber struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"unique;not null"`
	CreatedAt time.Time
}
type Post struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	Content   string
	CreatedAt time.Time
}

type Event struct {
	ID         uint `gorm:"primaryKey"`
	PostID     uint
	Subscriber string // email (simple for now)
	Type       string // "open" or "click"
	CreatedAt  time.Time
}
