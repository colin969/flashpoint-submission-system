package types

import "time"

type ActivityEvent struct {
	ID             int64
	UID            string
	CreatedAt      time.Time
	EventArea      ActivityEventArea
	EventOperation string
	EventData      interface{}
}

type ActivityEventSubmission struct {
	SID    int64
	Action string
}

type ActivityEventArea string

func (a *ActivityEventArea) Auth() string {
	return "auth"
}

func (a *ActivityEventArea) Submission() string {
	return "submission"
}

type ActivityEventOperation string

func (a *ActivityEventOperation) Create() string {
	return "create"
}

func (a *ActivityEventOperation) Update() string {
	return "update"
}

func (a *ActivityEventOperation) Delete() string {
	return "delete"
}
