package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Snippet struct {
	ID      int
	Title   string
	Content string
	MV1     string
	MV2     string
	MV3     string
	Created time.Time
	Expires time.Time
}
