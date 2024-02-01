package service

import (
	"context"

	"github.com/FlashpointProject/flashpoint-submission-system/activityevents"
	"github.com/FlashpointProject/flashpoint-submission-system/database"
	"github.com/FlashpointProject/flashpoint-submission-system/utils"
)

func (s *SiteService) EmitSubmissionDownloadEvent(ctx context.Context, userID, submissionID, fileID int64) error {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}
	defer dbs.Rollback()

	event := activityevents.BuildSubmissionDownloadEvent(userID, submissionID, fileID)

	err = s.pgdal.CreateEvent(dbs, event)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}
	return nil
}

func (s *SiteService) EmitSubmissionCommentEvent(pgdbs database.PGDBSession, userID int64, submissionID int64, commentID int64, action string, fileID *int64) error {
	ctx := pgdbs.Ctx()
	event := activityevents.BuildSubmissionCommentEvent(userID, submissionID, commentID, action, fileID)

	err := s.pgdal.CreateEvent(pgdbs, event)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}
	return nil
}
