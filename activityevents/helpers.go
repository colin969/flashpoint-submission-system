package activityevents

import (
	"strconv"
	"time"
)

func strptr(s string) *string {
	return &s
}

func BuildSubmissionCommentEvent(userID int64, submissionID, commentID int64, action string, fileID *int64) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    strconv.FormatInt(userID, 10),
		CreatedAt: time.Now(),
		Area:      aea.Submission(),
		Operation: aeo.Create(),
		Data: ActivityEventDataSubmission{
			Action:       &action,
			SubmissionID: &submissionID,
			CommentID:    &commentID,
			FileID:       fileID,
		},
	}
}

func BuildSubmissionDownloadEvent(userID int64, submissionID, fileID int64) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    strconv.FormatInt(userID, 10),
		CreatedAt: time.Now(),
		Area:      aea.Submission(),
		Operation: aeo.Read(),
		Data: ActivityEventDataSubmission{
			Action:       nil,
			SubmissionID: &submissionID,
			CommentID:    nil,
			FileID:       &fileID,
		},
	}
}
