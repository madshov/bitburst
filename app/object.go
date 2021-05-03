package app

import (
	"time"
)

type Object struct {
	ID        int       `json:"id"`
	Online    bool      `json:"online`
	Timestamp time.Time `json:"timestamp"`
}

type ObjectStorage interface {
	CreateObjects([]Object) error
	DeleteObjects() error
}
