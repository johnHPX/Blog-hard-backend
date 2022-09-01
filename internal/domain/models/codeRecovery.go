package models

import "time"

type CodeRecovery struct {
	Code      string
	UserID    string
	ExpiredAt time.Time
}
