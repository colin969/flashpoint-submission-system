package service

import (
	"context"

	"github.com/FlashpointProject/flashpoint-submission-system/activityevents"
	"github.com/FlashpointProject/flashpoint-submission-system/database"
	"github.com/FlashpointProject/flashpoint-submission-system/utils"
)

func (s *SiteService) EmitSubmissionDownloadEvent(ctx context.Context, userID, fileID int64) error {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}
	defer dbs.Rollback()

	event := activityevents.BuildSubmissionDownloadEvent(userID, fileID)

	err = s.pgdal.CreateEvent(dbs, event)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}
	return nil
}

func (s *SiteService) EmitSubmissionCommentEvent(dbs database.PGDBSession, userID int64, commentID int64, action string, fileID *int64) error {
	ctx := dbs.Ctx()
	event := activityevents.BuildSubmissionCommentEvent(userID, commentID, action, fileID)

	err := s.pgdal.CreateEvent(dbs, event)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}
	return nil
}
