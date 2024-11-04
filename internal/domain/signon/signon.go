package signon

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SignOnRepo interface {
	GetById(id int) (*SignOn, error)
	GetAll() ([]SignOn, error)
	GetByState(state ProcessingState) ([]SignOn, error)
	SaveAll(signons []SignOn) error
	UpdateContext(signon *SignOn) error
}

type ProcessingState string

const (
	StateProcessing ProcessingState = "processing"
	StateHandOver   ProcessingState = "handover"
	StateDone       ProcessingState = "done"
)

func (self ProcessingState) String() string {
	return string(self)
}

func ProcessingStateFromString(s string) ProcessingState {
	switch s {
	case "processing":
		return StateProcessing
	case "handover":
		return StateHandOver
	case "done":
		return StateDone
	default:
		return StateProcessing
	}
}

type SignOn struct {
	Id                   int
	IdPc                 int
	Company              sql.NullString
	Firstname            sql.NullString
	Lastname             sql.NullString
	Zip                  sql.NullString
	City                 sql.NullString
	Street               sql.NullString
	HouseNo              sql.NullString
	PCState              sql.NullString
	DesiredDeliveryStart sql.NullString
	MeterNo              sql.NullString
	Malo                 sql.NullString
	Melo                 sql.NullString
	ConfigId             sql.NullString
	MyState              ProcessingState
	MyComment            string
}
