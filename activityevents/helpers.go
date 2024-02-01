package activityevents

import (
	"strconv"
	"time"
)

func strptr(s string) *string {
	return &s
}

func BuildSubmissionCommentEvent(userID int64, commentID int64, action string, fileID *int64) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    strconv.FormatInt(userID, 10),
		CreatedAt: time.Now(),
		Area:      aea.Submission(),
		Operation: aeo.Create(),
		Data: ActivityEventDataSubmission{
			Action:    &action,
			CommentID: &commentID,
			FileID:    fileID,
		},
	}
}

func BuildSubmissionDownloadEvent(userID int64, fileID int64) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    strconv.FormatInt(userID, 10),
		CreatedAt: time.Now(),
		Area:      aea.Submission(),
		Operation: aeo.Read(),
		Data: ActivityEventDataSubmission{
			Action:    nil,
			CommentID: nil,
			FileID:    &fileID,
		},
	}
}
