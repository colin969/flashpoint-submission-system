package activityevents

import (
	"time"
)

// BuildSubmissionCreatedEvent is used for submission creation
func BuildSubmissionCreatedEvent(userID, submissionID int64) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Submission(),
		Operation: aeo.Create(),
		Data: &ActivityEventDataSubmission{
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
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Submission(),
		Operation: aeo.Update(),
		Data: &ActivityEventDataSubmission{
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
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Submission(),
		Operation: aeo.Read(),
		Data: &ActivityEventDataSubmission{
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
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Submission(),
		Operation: aeo.Delete(),
		Data: &ActivityEventDataSubmission{
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
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Submission(),
		Operation: aeo.Update(),
		Data: &ActivityEventDataSubmission{
			Action:       &action,
			SubmissionID: &submissionID,
			CommentID:    nil,
			FileID:       nil,
		},
	}
}

func BuildAuthLoginEvent(userID int64) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Auth(),
		Operation: aeo.Create(),
		Data: &ActivityEventDataAuth{
			Operation: "login",
		},
	}
}

func BuildAuthLogoutEvent(userID int64) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Auth(),
		Operation: aeo.Delete(),
		Data: &ActivityEventDataAuth{
			Operation: "logout",
		},
	}
}

func BuildGameLogoUpdateEvent(userID int64, gameUUID string) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Game(),
		Operation: aeo.Update(),
		Data: &ActivityEventDataGame{
			GameUUID:  gameUUID,
			Operation: "logo-update",
		},
	}
}

func BuildGameScreenshotUpdateEvent(userID int64, gameUUID string) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Game(),
		Operation: aeo.Update(),
		Data: &ActivityEventDataGame{
			GameUUID:  gameUUID,
			Operation: "screenshot-update",
		},
	}
}

func BuildGameDeleteEvent(userID int64, gameUUID string) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Game(),
		Operation: aeo.Delete(),
		Data: &ActivityEventDataGame{
			GameUUID:  gameUUID,
			Operation: "delete",
		},
	}
}

func BuildGameRestoreEvent(userID int64, gameUUID string) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Game(),
		Operation: aeo.Restore(),
		Data: &ActivityEventDataGame{
			GameUUID:  gameUUID,
			Operation: "restore",
		},
	}
}

func BuildGameFreezeEvent(userID int64, gameUUID string) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Game(),
		Operation: aeo.Update(),
		Data: &ActivityEventDataGame{
			GameUUID:  gameUUID,
			Operation: "freeze",
		},
	}
}

func BuildGameUnfreezeEvent(userID int64, gameUUID string) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Game(),
		Operation: aeo.Update(),
		Data: &ActivityEventDataGame{
			GameUUID:  gameUUID,
			Operation: "unfreeze",
		},
	}
}

func BuildAuthRevokeSessionEvent(userID, sessionID int64) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Auth(),
		Operation: aeo.Delete(),
		Data: &ActivityEventDataAuth{
			Operation: "revoke-session",
			SessionID: &sessionID,
		},
	}
}

func BuildAuthSetClientSecretEvent(userID int64, clientID string) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Auth(),
		Operation: aeo.Update(),
		Data: &ActivityEventDataAuth{
			Operation: "set-client-secret",
			ClientID:  &clientID,
		},
	}
}

func BuildTagUpdateEvent(userID, tagID int64) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Tag(),
		Operation: aeo.Update(),
		Data: &ActivityEventDataTag{
			TagID: tagID,
		},
	}
}

func BuildGameSaveEvent(userID int64, gameUUID string) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Game(),
		Operation: aeo.Update(),
		Data: &ActivityEventDataGame{
			GameUUID:  gameUUID,
			Operation: "save",
		},
	}
}

func BuildGameSaveDataEvent(userID int64, gameUUID string) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Game(),
		Operation: aeo.Update(),
		Data: &ActivityEventDataGame{
			GameUUID:  gameUUID,
			Operation: "save-data",
		},
	}
}

func BuildAuthDeviceEvent(userID int64, clientID string, approved bool) *ActivityEvent {
	op := "device-approve"
	if !approved {
		op = "device-deny"
	}
	return &ActivityEvent{
		ID:        -1,
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Auth(),
		Operation: aeo.Update(),
		Data: &ActivityEventDataAuth{
			Operation: op,
			ClientID:  &clientID,
		},
	}
}

func BuildAuthNewTokenEvent(userID int64, clientID string) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Auth(),
		Operation: aeo.Create(),
		Data: &ActivityEventDataAuth{
			Operation: "new-token",
			ClientID:  &clientID,
		},
	}
}

func BuildAuthDeleteUserSessionsEvent(userID int64, targetID int64) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Auth(),
		Operation: aeo.Delete(),
		Data: &ActivityEventDataAuth{
			Operation:    "delete-user-sessions",
			TargetUserID: &targetID,
		},
	}
}

func BuildGameRedirectEvent(userID int64, fromGameUUID, toGameUUID string) *ActivityEvent {
	return &ActivityEvent{
		ID:        -1,
		UserID:    userID,
		CreatedAt: time.Now(),
		Area:      aea.Game(),
		Operation: aeo.Create(),
		Data: &ActivityEventDataGame{
			GameUUID:          fromGameUUID,
			Operation:         "redirect",
			SecondaryGameUUID: &toGameUUID,
		},
	}
}
