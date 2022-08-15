package model

import "time"

type Access struct {
	Token     string
	UserID    string
	ExpiredAt time.Time
	IsBlocked bool
}
