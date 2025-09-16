package models

import "time"

type Post struct {
	ID               int64     `json:"id"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	CommentsDisabled bool      `json:"comments_disabled"`
	CreatedAt        time.Time `json:"created_at"`
}
