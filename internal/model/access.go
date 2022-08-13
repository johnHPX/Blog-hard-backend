package model

import "time"

type Access struct {
	Token     string
	UserId    string
	ExpiredAt time.Time
	IsBlocked time.Time
}
