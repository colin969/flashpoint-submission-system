package activityevents

import (
	"strconv"
	"time"
)

func strptr(s string) *string {
	return &s
}

// BuildSubmissionCreatedEvent is used for submission creation
func BuildSubmissionCreatedEvent(userID, submissionID int64) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    strconv.FormatInt(userID, 10),
		CreatedAt: time.Now(),
		Area:      aea.Submission(),
		Operation: aeo.Create(),
		Data: ActivityEventDataSubmission{
			Action:       nil,
			SubmissionID: &submissionID,
			CommentID:    nil,
			FileID:       nil,
		},
	}
}

// BuildSubmissionCommentEvent is used for any comment received on a submission
func BuildSubmissionCommentEvent(userID int64, submissionID, commentID int64, action string, fileID *int64) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    strconv.FormatInt(userID, 10),
		CreatedAt: time.Now(),
		Area:      aea.Submission(),
		Operation: aeo.Update(),
		Data: ActivityEventDataSubmission{
			Action:       &action,
			SubmissionID: &submissionID,
			CommentID:    &commentID,
			FileID:       fileID,
		},
	}
}

// BuildSubmissionDownloadEvent is used for file downloads
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

// BuildSubmissionDeleteEvent is used for any submission deletion event (submission, comment, file)
func BuildSubmissionDeleteEvent(userID int64, submissionID int64, commentID *int64, fileID *int64) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    strconv.FormatInt(userID, 10),
		CreatedAt: time.Now(),
		Area:      aea.Submission(),
		Operation: aeo.Delete(),
		Data: ActivityEventDataSubmission{
			Action:       nil,
			SubmissionID: &submissionID,
			CommentID:    commentID,
			FileID:       fileID,
		},
	}
}

// BuildSubmissionFreezeEvent is used for manual freeze/unfreeze
func BuildSubmissionFreezeEvent(userID int64, submissionID int64, toFreeze bool) *ActivityEvent {
	action := "freeze"
	if !toFreeze {
		action = "unfreeze"
	}

	return &ActivityEvent{
		ID:        -1,
		UserID:    strconv.FormatInt(userID, 10),
		CreatedAt: time.Now(),
		Area:      aea.Submission(),
		Operation: aeo.Update(),
		Data: ActivityEventDataSubmission{
			Action:       &action,
			SubmissionID: &submissionID,
			CommentID:    nil,
			FileID:       nil,
		},
	}
}
