package activityevents

import "time"

type ActivityEvent struct {
	ID        int64                  `json:"id"`
	UserID    string                 `json:"user_id"`
	CreatedAt time.Time              `json:"created_at"`
	Area      ActivityEventArea      `json:"area"`
	Operation ActivityEventOperation `json:"operation"`
	Data      interface{}            `json:"data"`
}

type ActivityEventArea string

var aea *ActivityEventArea

func (*ActivityEventArea) Auth() ActivityEventArea {
	return "auth"
}

func (*ActivityEventArea) Admin() ActivityEventArea {
	return "admin"
}

func (*ActivityEventArea) Submission() ActivityEventArea {
	return "submission"
}

func (*ActivityEventArea) Tag() ActivityEventArea {
	return "tag"
}

func (*ActivityEventArea) Game() ActivityEventArea {
	return "game"
}

type ActivityEventOperation string

var aeo *ActivityEventOperation

func (*ActivityEventOperation) Create() ActivityEventOperation {
	return "create"
}

func (*ActivityEventOperation) Read() ActivityEventOperation {
	return "read"
}

func (*ActivityEventOperation) Update() ActivityEventOperation {
	return "update"
}

func (*ActivityEventOperation) Delete() ActivityEventOperation {
	return "delete"
}

func (*ActivityEventOperation) Restore() ActivityEventOperation {
	return "restore"
}

type ActivityEventDataSubmission struct {
	Action       *string `json:"action"`
	SubmissionID *int64  `json:"submission_id"`
	CommentID    *int64  `json:"comment_id"`
	FileID       *int64  `json:"file_id"`
}

type ActivityEventDataAuth struct {
	Operation string  `json:"operation"`
	SessionID *int64  `json:"session_id"`
	ClientID  *string `json:"client_id"`
}

type ActivityEventDataGame struct {
	GameUUID  string `json:"game_uuid"`
	Operation string `json:"operation"`
}
