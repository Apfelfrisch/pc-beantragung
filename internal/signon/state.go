package signon

import (
	_ "github.com/mattn/go-sqlite3"
)

type ProcessingState string

const (
	StateProcessing ProcessingState = "processing"
	StateHandOver   ProcessingState = "handover"
	StateDone       ProcessingState = "done"
)

func (self ProcessingState) String() string {
	return string(self)
}
