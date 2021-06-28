package service

import (
	"fmt"
	"github.com/Dri0m/flashpoint-submission-system/constants"
	"github.com/Dri0m/flashpoint-submission-system/database"
	"github.com/Dri0m/flashpoint-submission-system/types"
	"github.com/Dri0m/flashpoint-submission-system/utils"
	"strings"
)

// createNotification formats and stores notification
func (s *SiteService) createNotification(dbs database.DBSession, authorID, sid int64, action string) error {
	validAction := false
	for _, a := range constants.GetActionsWithNotification() {
		if action == a {
			validAction = true
			break
		}
	}
	if !validAction {
		return nil
	}

	mentionUserIDs, err := s.dal.GetUsersForNotification(dbs, authorID, sid, action)
	if err != nil {
		utils.LogCtx(dbs.Ctx()).Error(err)
		return err
	}

	if len(mentionUserIDs) == 0 {
		return nil
	}

	var b strings.Builder
	b.WriteString("You've got mail!\n")
	b.WriteString(fmt.Sprintf("<https://fpfss.unstable.life/submission/%d>\n", sid))

	if action == constants.ActionComment {
		b.WriteString(fmt.Sprintf("There is a new comment on the submission."))
	} else if action == constants.ActionApprove {
		b.WriteString(fmt.Sprintf("The submission has been approved."))
	} else if action == constants.ActionRequestChanges {
		b.WriteString(fmt.Sprintf("User has requested changes on the submission."))
	} else if action == constants.ActionMarkAdded {
		b.WriteString(fmt.Sprintf("The submission has been marked as added to Flashpoint."))
	} else if action == constants.ActionUpload {
		b.WriteString(fmt.Sprintf("A new version has been uploaded by <@%d>", authorID))
	}
	b.WriteString("\n")

	for _, userID := range mentionUserIDs {
		b.WriteString(fmt.Sprintf(" <@%d>", userID))
	}

	b.WriteString("\n----------------------------------------------------------\n")
	msg := b.String()

	if err := s.dal.StoreNotification(dbs, msg, constants.NotificationDefault); err != nil {
		utils.LogCtx(dbs.Ctx()).Error(err)
		return dberr(err)
	}

	return nil
}

// createCurationFeedMessage formats and stores message for the curation feed
func (s *SiteService) createCurationFeedMessage(dbs database.DBSession, authorID, sid int64, isSubmissionNew, isCurationValid bool, meta *types.CurationMeta) error {
	var b strings.Builder

	if isSubmissionNew {
		b.WriteString(fmt.Sprintf("A new submission has been uploaded by <@%d>\n", authorID))
	} else {
		b.WriteString(fmt.Sprintf("A submission update has been uploaded by <@%d>\n", authorID))
	}
	b.WriteString(fmt.Sprintf("<https://fpfss.unstable.life/submission/%d>\n", sid))

	if !isCurationValid {
		b.WriteString("Unfortunately, it does not quite reach the quality required to satisfy the cool crab.\n")
	}

	if meta.Library != nil && meta.Platform != nil && meta.Title != nil && meta.Extreme != nil {
		llib := strings.ToLower(*meta.Library)
		if strings.Contains(llib, "arcade") {
			b.WriteString("🎮")
		} else if strings.Contains(llib, "theatre") {
			b.WriteString("🎞️")
		} else {
			b.WriteString("❓")
		}

		b.WriteString(" ")

		lplat := strings.ToLower(*meta.Platform)
		if strings.Contains(lplat, "3d groove fx") {
			b.WriteString("<:3DGroove:569691574276063242>")
		} else if strings.Contains(lplat, "3dvia player") {
			b.WriteString("<:3DVIA_Player:496151464784166946")
		} else if strings.Contains(lplat, "axel player") {
			b.WriteString("<:AXEL_Player:813079894267265094>")
		} else if strings.Contains(lplat, "activex") {
			b.WriteString("<:ActiveX:699093212949643365>")
		} else if strings.Contains(lplat, "atmosphere") {
			b.WriteString("<:Atmosphere:781105689002901524>")
		} else if strings.Contains(lplat, "authorware") {
			b.WriteString("<:Authorware:582105144410243073>")
		} else if strings.Contains(lplat, "burster") {
			b.WriteString("<:Burster:743995494736461854>")
		} else if strings.Contains(lplat, "cult3d") {
			b.WriteString("<:Cult3D:806277196473040896>")
		} else if strings.Contains(lplat, "deepv") {
			b.WriteString("<:DeepV:812079774843142255>")
		} else if strings.Contains(lplat, "flash") {
			b.WriteString("<:Flash:750823911326875648>")
		} else if strings.Contains(lplat, "gobit") {
			b.WriteString("<:GoBit:629511736608686080>")
		} else if strings.Contains(lplat, "html5") {
			b.WriteString("<:HTML5:701930562746712142>")
		} else if strings.Contains(lplat, "hyper-g") {
			b.WriteString("<:HyperG:817543962088570880>")
		} else if strings.Contains(lplat, "hypercosm") {
			b.WriteString("<:Hypercosm:814623525038063697>")
		} else if strings.Contains(lplat, "java") {
			b.WriteString("<:Java:482697866377297920>")
		} else if strings.Contains(lplat, "livemath") {
			b.WriteString("<:LiveMath_Plugin:808999958043951104>")
		} else if strings.Contains(lplat, "octree view") {
			b.WriteString("<:Octree_View:809147835927756831>")
		} else if strings.Contains(lplat, "play3d") {
			b.WriteString("<:Play3D:812079775152734209>")
		} else if strings.Contains(lplat, "popcap plugin") {
			b.WriteString("<:PopCap:604433459179552798>")
		} else if strings.Contains(lplat, "protoplay") {
			b.WriteString("<:ProtoPlay:806614012829761587>")
		} else if strings.Contains(lplat, "pulse") {
			b.WriteString("<:Pulse:720682372982505472>")
		} else if strings.Contains(lplat, "rebol") {
			b.WriteString("<:REBOL:806995243085987862>")
		} else if strings.Contains(lplat, "shiva3d") {
			b.WriteString("<:ShiVa3d:643124144812326934>")
		} else if strings.Contains(lplat, "shockwave") {
			b.WriteString("<:Shockwave:727436274625019965>")
		} else if strings.Contains(lplat, "silverlight") {
			b.WriteString("<:Silverlight:492112373625257994>")
		} else if strings.Contains(lplat, "tcl") {
			b.WriteString("<:Tcl:737419431067779144>")
		} else if strings.Contains(lplat, "unity") {
			b.WriteString("<:Unity:600478910169481216>")
		} else if strings.Contains(lplat, "vrml") {
			b.WriteString("<:VRML:737049432817664070>")
		} else if strings.Contains(lplat, "viscape") {
			b.WriteString("<:Viscape:814623877039652886>")
		} else if strings.Contains(lplat, "vitalize") {
			b.WriteString("<:Vitalize:700924839912800332>")
		} else if strings.Contains(lplat, "xara plugin") {
			b.WriteString("<:Xara_Plugin:807439131768258561>")
		} else if strings.Contains(lplat, "alambik") {
			b.WriteString("<:Alambik:814621713350262856>")
		} else if strings.Contains(lplat, "animaflex") {
			b.WriteString("<:AnimaFlex:807016001618968596>")
		} else {
			b.WriteString("❓")
		}

		b.WriteString(" ")

		if *meta.Extreme == "Yes" {
			b.WriteString("<:extreme:778145279714918400>")
		}

		b.WriteString(" ")

		b.WriteString(*meta.Title)
		b.WriteString("\n")
	}

	b.WriteString("----------------------------------------------------------\n")
	msg := b.String()

	if err := s.dal.StoreNotification(dbs, msg, constants.NotificationCurationFeed); err != nil {
		return err
	}

	return nil
}

// createDeletionNotification formats and stores deletion notification
func (s *SiteService) createDeletionNotification(dbs database.DBSession, authorID, deleterID int64, sid, cid, fid *int64, reason string) error {
	if sid == nil {
		utils.LogCtx(dbs.Ctx()).Fatal("submission id cannot be nil")
	}
	if cid != nil && fid != nil {
		utils.LogCtx(dbs.Ctx()).Fatal("both cid and fid provided - not valid")
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("You've got mail! <@%d>\n", authorID))
	b.WriteString(fmt.Sprintf("<https://fpfss.unstable.life/submission/%d>\n", *sid))
	if cid != nil {
		b.WriteString(fmt.Sprintf("Your comment #%d was deleted by <@%d>\n", *cid, deleterID))
	} else if fid != nil {
		b.WriteString(fmt.Sprintf("Your file #%d was deleted by <@%d>\n", *fid, deleterID))
	} else {
		b.WriteString(fmt.Sprintf("Your submission #%d was deleted by <@%d>\n", *sid, deleterID))
	}
	b.WriteString(fmt.Sprintf("Reason: %s", reason))
	b.WriteString("\n----------------------------------------------------------\n")
	msg := b.String()

	if err := s.dal.StoreNotification(dbs, msg, constants.NotificationDefault); err != nil {
		utils.LogCtx(dbs.Ctx()).Error(err)
		return dberr(err)
	}

	return nil
}