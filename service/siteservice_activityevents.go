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

func (s *SiteService) EmitSubmissionCreatedEvent(pgdbs database.PGDBSession, userID, submissionID int64) error {
	ctx := pgdbs.Ctx()
	event := activityevents.BuildSubmissionCreatedEvent(userID, submissionID)

	err := s.pgdal.CreateEvent(pgdbs, event)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}
	return nil
}

func (s *SiteService) EmitSubmissionCommentEvent(pgdbs database.PGDBSession, userID, submissionID, commentID int64, action string, fileID *int64) error {
	ctx := pgdbs.Ctx()
	event := activityevents.BuildSubmissionCommentEvent(userID, submissionID, commentID, action, fileID)

	err := s.pgdal.CreateEvent(pgdbs, event)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}
	return nil
}

func (s *SiteService) EmitSubmissionOverrideEvent(pgdbs database.PGDBSession, userID, submissionID, commentID int64) error {
	ctx := pgdbs.Ctx()
	event := activityevents.BuildSubmissionCommentEvent(userID, submissionID, commentID, "approve-override", nil)

	err := s.pgdal.CreateEvent(pgdbs, event)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}
	return nil
}

func (s *SiteService) EmitSubmissionDeleteEvent(pgdbs database.PGDBSession, userID, submissionID int64, commentID, fileID *int64) error {
	ctx := pgdbs.Ctx()
	event := activityevents.BuildSubmissionDeleteEvent(userID, submissionID, commentID, fileID)

	err := s.pgdal.CreateEvent(pgdbs, event)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}
	return nil
}

func (s *SiteService) EmitSubmissionFreezeEvent(pgdbs database.PGDBSession, userID, submissionID int64, toFreeze bool) error {
	ctx := pgdbs.Ctx()
	event := activityevents.BuildSubmissionFreezeEvent(userID, submissionID, toFreeze)

	err := s.pgdal.CreateEvent(pgdbs, event)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}
	return nil
}