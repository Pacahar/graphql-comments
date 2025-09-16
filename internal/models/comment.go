package models

import "time"

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	ParentID  *int64    `json:"parent_id,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
