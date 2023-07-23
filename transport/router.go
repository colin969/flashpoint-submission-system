package transport

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Dri0m/flashpoint-submission-system/constants"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (a *App) handleRequests(l *logrus.Entry, srv *http.Server, router *mux.Router) {
	isStaff := func(r *http.Request, uid int64) (bool, error) {
		return a.UserHasAnyRole(r, uid, constants.StaffRoles())
	}
	isTrialEditor := func(r *http.Request, uid int64) (bool, error) {
		return a.UserHasAnyRole(r, uid, constants.TrialEditorRoles())
	}
	isTrialCurator := func(r *http.Request, uid int64) (bool, error) {
		return a.UserHasAnyRole(r, uid, constants.TrialCuratorRoles())
	}
	isDeleter := func(r *http.Request, uid int64) (bool, error) {
		return a.UserHasAnyRole(r, uid, constants.DeleterRoles())
	}
	isFreezer := func(r *http.Request, uid int64) (bool, error) {
		return a.UserHasAnyRole(r, uid, constants.FreezerRoles())
	}
	isInAudit := func(r *http.Request, uid int64) (bool, error) {
		s, err := a.UserHasAnyRole(r, uid, constants.StaffRoles())
		if err != nil {
			return false, err
		}
		t, err := a.UserHasAnyRole(r, uid, constants.TrialCuratorRoles())
		if err != nil {
			return false, err
		}
		return !(s || t), nil
	}
	isColin := func(r *http.Request, uid int64) (bool, error) {
		return uid == 689080719460663414, nil
	}
	isGod := func(r *http.Request, uid int64) (bool, error) {
		s, err := isColin(r, uid)
		if err != nil || s == true {
			return s, err
		} else {
			return a.UserHasAnyRole(r, uid, constants.GodRoles())
		}
	}
	userOwnsSubmission := func(r *http.Request, uid int64) (bool, error) {
		return a.UserOwnsResource(r, uid, constants.ResourceKeySubmissionID)
	}
	userOwnsAllSubmissions := func(r *http.Request, uid int64) (bool, error) {
		return a.UserOwnsResource(r, uid, constants.ResourceKeySubmissionIDs)
	}
	userHasNoSubmissions := func(r *http.Request, uid int64) (bool, error) {
		return a.IsUserWithinResourceLimit(r, uid, constants.ResourceKeySubmissionID, 1)
	}
	isSubmissionFrozen := func(r *http.Request, uid int64) (bool, error) {
		return a.IsResourceFrozen(r, constants.ResourceKeySubmissionID)
	}
	isAnySubmissionFrozen := func(r *http.Request, uid int64) (bool, error) {
		return a.IsResourceFrozen(r, constants.ResourceKeySubmissionIDs)
	}
	isAnySubmissionFileFrozen := func(r *http.Request, uid int64) (bool, error) {
		return a.IsResourceFrozen(r, constants.ResourceKeyFileIDs)
	}
	isSubmissionMarkedAsAdded := func(r *http.Request, uid int64) (bool, error) {
		return a.IsResourceMarkedAsAdded(r, constants.ResourceKeySubmissionID)
	}
	isAnySubmissionMarkedAsAdded := func(r *http.Request, uid int64) (bool, error) {
		return a.IsResourceMarkedAsAdded(r, constants.ResourceKeySubmissionIDs)
	}
	isActionMarkAsAddedForMultipleSubmissions := func(r *http.Request, uid int64) (bool, error) {
		params := mux.Vars(r)
		submissionIDs := strings.Split(params[constants.ResourceKeySubmissionIDs], ",")
		if len(submissionIDs) == 1 {
			return false, nil
		}
		return a.IsAction(r, constants.ActionMarkAdded)
	}

	// static file server
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", NoCache(http.FileServer(http.Dir("./static/")))))

	// auth
	router.Handle(
		"/auth",
		http.HandlerFunc(a.RequestWeb(a.HandleDiscordAuth))).
		Methods("GET")
	router.Handle(
		"/auth/callback",
		http.HandlerFunc(a.RequestWeb(a.HandleDiscordCallback))).
		Methods("GET")
	router.Handle(
		"/api/logout",
		http.HandlerFunc(a.RequestJSON(a.HandleLogout))).
		Methods("GET")

	// device flow
	router.Handle(
		"/auth/token",
		http.HandlerFunc(a.RequestJSON(a.HandleNewDeviceToken))).
		Methods("GET", "POST")
	router.Handle(
		"/auth/device",
		http.HandlerFunc(a.UserAuthMux(a.RequestWeb(a.HandleApproveDevice), muxAny(isStaff, isTrialCurator, isInAudit)))).
		Methods("GET", "POST")

	// pages
	router.Handle(
		"/",
		http.HandlerFunc(a.RequestWeb(a.HandleRootPage))).
		Methods("GET")

	router.Handle(
		"/web",
		http.HandlerFunc(a.RequestWeb(a.HandleRootPage))).
		Methods("GET")

	router.Handle(
		"/web/submit",
		http.HandlerFunc(a.RequestWeb(a.UserAuthMux(
			a.HandleSubmitPage, muxAny(isStaff, isTrialCurator, isInAudit))))).
		Methods("GET")

	router.Handle(
		"/web/help",
		http.HandlerFunc(a.RequestWeb(a.HandleHelpPage))).
		Methods("GET")

	router.Handle(
		"/web/flashfreeze/submit",
		http.HandlerFunc(a.RequestWeb(a.UserAuthMux(
			a.HandleFlashfreezeSubmitPage, muxAny(isStaff, isTrialCurator, isInAudit))))).
		Methods("GET")

	////////////////////////

	f := a.UserAuthMux(a.HandleProfilePage)

	router.Handle(
		"/web/profile",
		http.HandlerFunc(a.RequestWeb(f))).
		Methods("GET")

	router.Handle(
		"/api/profile",
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET")

	////////////////////////

	f = a.UserAuthMux(
		a.HandleSubmissionsPage, muxAny(isStaff, isTrialCurator, isInAudit))

	router.Handle(
		"/web/submissions",
		http.HandlerFunc(a.RequestWeb(f))).
		Methods("GET")

	router.Handle(
		"/api/submissions",
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET")

	/////////////////////////

	f = a.HandleTagsPage

	router.Handle(
		"/web/tags",
		http.HandlerFunc(a.RequestWeb(f))).
		Methods("GET")

	router.Handle(
		"/api/tags",
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET")

	f = a.UserAuthMux(
		a.HandlePostTag, muxAny(isDeleter))

	router.Handle(
		fmt.Sprintf("/api/tag/{%s}", constants.ResourceKeyTagID),
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("POST")

	f = a.UserAuthMux(
		a.HandleTagPage, muxAny(isStaff, isTrialCurator, isInAudit))

	router.Handle(
		fmt.Sprintf("/web/tag/{%s}", constants.ResourceKeyTagID),
		http.HandlerFunc(a.RequestWeb(f))).
		Methods("GET")

	router.Handle(
		fmt.Sprintf("/api/tag/{%s}", constants.ResourceKeyTagID),
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET")

	f = a.UserAuthMux(
		a.HandleTagEditPage, muxAny(isDeleter))

	router.Handle(
		fmt.Sprintf("/web/tag/{%s}/edit", constants.ResourceKeyTagID),
		http.HandlerFunc(a.RequestWeb(f))).
		Methods("GET")

	////////////////////////

	f = a.HandleMetadataStats

	router.Handle(
		"/web/metadata-stats",
		http.HandlerFunc(a.RequestWeb(f))).
		Methods("GET")

	////////////////////////

	f = a.HandleMinLauncherVersion

	router.Handle(
		"/api/min-launcher",
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET")

	////////////////////////

	f = a.HandleGameCountSinceDate

	router.Handle(
		"/api/games/updates",
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET")

	f = a.HandleGamesPage

	router.Handle(
		"/api/games",
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET")

	f = a.HandleDeletedGames

	router.Handle(
		"/api/games/deleted",
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET")

	////////////////////////

	f = a.HandlePlatformsPage

	router.Handle(
		"/web/platforms",
		http.HandlerFunc(a.RequestWeb(f))).
		Methods("GET")

	router.Handle(
		"/api/platforms",
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET")

	////////////////////////

	f = a.UserAuthMux(
		a.HandleGamePage, muxAny(isStaff, isTrialEditor))

	router.Handle(
		fmt.Sprintf("/web/game/{%s}", constants.ResourceKeyGameID),
		http.HandlerFunc(a.RequestWeb(f))).
		Methods("GET")

	router.Handle(
		fmt.Sprintf("/web/game/{%s}/revision/{%s}", constants.ResourceKeyGameID, constants.ResourceKeyGameRevision),
		http.HandlerFunc(a.RequestWeb(f))).
		Methods("GET")

	router.Handle(
		fmt.Sprintf("/api/game/{%s}", constants.ResourceKeyGameID),
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET", "POST")

	f = a.UserAuthMux(
		a.HandleGameDataIndexPage, muxAny(isTrialCurator, isStaff))

	router.Handle(
		fmt.Sprintf("/web/game/{%s}/data/{%s}/index", constants.ResourceKeyGameID, constants.ResourceKeyGameDataDate),
		http.HandlerFunc(a.RequestWeb(f))).
		Methods("GET")

	router.Handle(
		fmt.Sprintf("/api/game/{%s}/data/{%s}/index", constants.ResourceKeyGameID, constants.ResourceKeyGameDataDate),
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET")

	f = a.UserAuthMux(
		a.HandleGameDataEditPage, muxAny(isStaff))

	router.Handle(
		fmt.Sprintf("/web/game/{%s}/data/{%s}/edit", constants.ResourceKeyGameID, constants.ResourceKeyGameDataDate),
		http.HandlerFunc(a.RequestWeb(f))).
		Methods("GET")

	router.Handle(
		fmt.Sprintf("/web/game/{%s}/data/{%s}/edit", constants.ResourceKeyGameID, constants.ResourceKeyGameDataDate),
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("POST")

	f = a.UserAuthMux(
		a.HandleDeleteGame, muxAny(isDeleter))

	router.Handle(
		fmt.Sprintf("/api/game/{%s}", constants.ResourceKeyGameID),
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("DELETE")

	f = a.UserAuthMux(
		a.HandleRestoreGame, muxAny(isDeleter))

	router.Handle(
		fmt.Sprintf("/api/game/{%s}/restore", constants.ResourceKeyGameID),
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET")

	f = a.UserAuthMux(
		a.HandleGameLogo, muxAny(isStaff))

	router.Handle(
		fmt.Sprintf("/api/game/{%s}/logo", constants.ResourceKeyGameID),
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("POST")

	f = a.UserAuthMux(
		a.HandleGameScreenshot, muxAny(isStaff))

	router.Handle(
		fmt.Sprintf("/api/game/{%s}/screenshot", constants.ResourceKeyGameID),
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("POST")

	f = a.UserAuthMux(
		a.HandleFreezeGame, muxAny(isFreezer))

	router.Handle(
		fmt.Sprintf("/api/game/{%s}/freeze", constants.ResourceKeyGameID),
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("POST")

	f = a.UserAuthMux(
		a.HandleUnfreezeGame, muxAny(isFreezer))

	router.Handle(
		fmt.Sprintf("/api/game/{%s}/unfreeze", constants.ResourceKeyGameID),
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("POST")

	////////////////////////

	f = a.UserAuthMux(
		a.HandleMatchingIndexHash, muxAny(isStaff, isTrialCurator))

	router.Handle(
		fmt.Sprintf("/api/index/hash/{%s}", constants.ResourceKeyHash),
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("POST")

	////////////////////////

	f = a.UserAuthMux(
		a.HandleMySubmissionsPage, muxAny(isStaff, isTrialCurator, isInAudit))

	router.Handle(
		"/web/my-submissions",
		http.HandlerFunc(a.RequestWeb(f))).
		Methods("GET")

	router.Handle(
		"/api/my-submissions",
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET")

	////////////////////////

	f = a.UserAuthMux(
		a.HandleViewSubmissionPage,
		muxAny(isStaff, isTrialCurator, isInAudit))

	router.Handle(
		fmt.Sprintf("/web/submission/{%s}", constants.ResourceKeySubmissionID),
		http.HandlerFunc(a.RequestWeb(f))).
		Methods("GET")

	router.Handle(
		fmt.Sprintf("/api/submission/{%s}", constants.ResourceKeySubmissionID),
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET")

	f = a.UserAuthMux(
		a.HandleApplyContentPatchPage,
		isDeleter)

	router.Handle(
		fmt.Sprintf("/web/submission/{%s}/apply", constants.ResourceKeySubmissionID),
		http.HandlerFunc(a.RequestWeb(f))).
		Methods("GET")

	////////////////////////

	f = a.UserAuthMux(
		a.HandleViewSubmissionFilesPage,
		muxAny(isStaff, isTrialCurator, isInAudit))

	router.Handle(
		fmt.Sprintf("/web/submission/{%s}/files", constants.ResourceKeySubmissionID),
		http.HandlerFunc(a.RequestWeb(f))).
		Methods("GET")

	router.Handle(
		fmt.Sprintf("/api/submission/{%s}/files", constants.ResourceKeySubmissionID),
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET")

	////////////////////////

	f = a.UserAuthMux(
		a.HandleSearchFlashfreezePage,
		muxAny(isStaff, isTrialCurator, isInAudit))

	router.Handle(
		"/web/flashfreeze/files",
		http.HandlerFunc(a.RequestWeb(f))).
		Methods("GET")

	router.Handle(
		"/api/flashfreeze/files",
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET")

	////////////////////////

	f = a.UserAuthMux(a.HandleStatisticsPage)

	router.Handle(
		"/web/statistics",
		http.HandlerFunc(a.RequestWeb(f))).
		Methods("GET")

	router.Handle(
		"/api/statistics",
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET")

	////////////////////////

	f = a.UserAuthMux(a.HandleUserStatisticsPage)

	router.Handle(
		"/web/user-statistics",
		http.HandlerFunc(a.RequestWeb(f))).
		Methods("GET")

	router.Handle(
		"/api/user-statistics",
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET")

	////////////////////////

	// receivers

	////////////////////////

	router.Handle(
		"/api/submission-receiver-resumable",
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(
			a.HandleSubmissionReceiverResumable, muxAny(
				isStaff,
				isTrialCurator,
				muxAll(isInAudit, userHasNoSubmissions)))))).
		Methods("POST")

	router.Handle(
		fmt.Sprintf("/api/submission-receiver-resumable/{%s}", constants.ResourceKeySubmissionID),
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(
			a.HandleSubmissionReceiverResumable, muxAny(
				isStaff,
				muxAll(isTrialCurator, userOwnsSubmission),
				muxAll(isInAudit, userOwnsSubmission)))))).
		Methods("POST")

	router.Handle(
		"/api/submission-receiver-resumable",
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(
			a.HandleReceiverResumableTestChunk, muxAny(
				isStaff,
				isTrialCurator,
				muxAll(isInAudit, userHasNoSubmissions)))))).
		Methods("GET")

	router.Handle(
		fmt.Sprintf("/api/submission-receiver-resumable/{%s}", constants.ResourceKeySubmissionID),
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(
			a.HandleReceiverResumableTestChunk, muxAny(
				isStaff,
				muxAll(isTrialCurator, userOwnsSubmission),
				muxAll(isInAudit, userOwnsSubmission)))))).
		Methods("GET")

	////////////////////////

	router.Handle(
		"/api/flashfreeze-receiver-resumable",
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(
			a.HandleFlashfreezeReceiverResumable,
			muxAny(isStaff, isTrialCurator, isInAudit))))).
		Methods("POST")

	router.Handle(
		"/api/flashfreeze-receiver-resumable",
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(
			a.HandleReceiverResumableTestChunk,
			muxAny(isStaff, isTrialCurator, isInAudit))))).
		Methods("GET")

	////////////////////////

	router.Handle(
		fmt.Sprintf("/api/submission-batch/{%s}/comment", constants.ResourceKeySubmissionIDs),
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(
			a.HandleCommentReceiverBatch,
			muxAll(
				muxNot(isActionMarkAsAddedForMultipleSubmissions),
				muxNot(isAnySubmissionMarkedAsAdded),
				muxAny(
					muxAll(muxNot(isAnySubmissionFrozen), isStaff, a.UserCanCommentAction), // TODO plural!
					muxAll(muxNot(isAnySubmissionFrozen), isTrialCurator, userOwnsAllSubmissions, a.UserCanCommentAction),
					muxAll(muxNot(isAnySubmissionFrozen), isInAudit, userOwnsAllSubmissions, a.UserCanCommentAction),
					muxAll(isAnySubmissionFrozen, isFreezer, a.UserCanCommentAction),
				)))))).
		Methods("POST")

	router.Handle("/api/notification-settings",
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(
			a.HandleUpdateNotificationSettings, muxAny(isStaff, isTrialCurator, isInAudit))))).
		Methods("PUT")

	router.Handle(
		fmt.Sprintf("/api/submission/{%s}/subscription-settings", constants.ResourceKeySubmissionID),
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(
			a.HandleUpdateSubscriptionSettings, muxAny(isStaff, isTrialCurator, isInAudit))))).
		Methods("PUT")

	////////////////////////

	router.Handle(
		"/web/developer",
		http.HandlerFunc(a.RequestWeb(a.UserAuthMux(
			a.HandleDeveloperPage, muxAny(isGod, isColin))))).
		Methods("GET")

	router.Handle(
		"/api/developer/submit_dump",
		http.HandlerFunc(a.RequestWeb(a.UserAuthMux(
			a.HandleDeveloperDumpUpload, muxAny(isGod, isColin))))).
		Methods("POST")

	router.Handle(
		"/api/developer/tag_desc_from_validator",
		http.HandlerFunc(a.RequestWeb(a.UserAuthMux(
			a.HandleDeveloperTagDescFromValidator, muxAny(isGod, isColin))))).
		Methods("GET")

	////////////////////////

	// providers

	////////////////////////

	router.Handle(
		fmt.Sprintf("/data/submission/{%s}/file/{%s}", constants.ResourceKeySubmissionID, constants.ResourceKeyFileID),
		http.HandlerFunc(a.RequestData(a.UserAuthMux(
			a.HandleDownloadSubmissionFile,
			muxAny(
				muxAll(muxNot(isSubmissionFrozen), isStaff),
				muxAll(muxNot(isSubmissionFrozen), isTrialCurator),
				muxAll(muxNot(isSubmissionFrozen), isInAudit),
				muxAll(isSubmissionFrozen, isFreezer, muxAny(isStaff, isTrialCurator, isInAudit)),
			))))).
		Methods("GET")

	router.Handle(
		fmt.Sprintf("/data/submission-file-batch/{%s}", constants.ResourceKeyFileIDs),
		http.HandlerFunc(a.RequestData(a.UserAuthMux(
			a.HandleDownloadSubmissionBatch, muxAny(
				muxAll(muxNot(isAnySubmissionFileFrozen), isStaff),
				muxAll(muxNot(isAnySubmissionFileFrozen), isTrialCurator),
				muxAll(muxNot(isAnySubmissionFileFrozen), isInAudit),
				muxAll(isAnySubmissionFileFrozen, isFreezer, muxAny(isStaff, isTrialCurator, isInAudit)),
			))))).
		Methods("GET")

	router.Handle(
		fmt.Sprintf("/data/submission/{%s}/curation-image/{%s}.png", constants.ResourceKeySubmissionID, constants.ResourceKeyCurationImageID),
		http.HandlerFunc(a.RequestData(a.UserAuthMux(
			a.HandleDownloadCurationImage,
			muxAny(
				muxAll(muxNot(isSubmissionFrozen), isStaff),
				muxAll(muxNot(isSubmissionFrozen), isTrialCurator),
				muxAll(muxNot(isSubmissionFrozen), isInAudit),
				muxAll(isSubmissionFrozen, isFreezer, muxAny(isStaff, isTrialCurator, isInAudit)),
			))))).
		Methods("GET")

	router.Handle(
		fmt.Sprintf("/data/flashfreeze/file/{%s}", constants.ResourceKeyFlashfreezeRootFileID),
		http.HandlerFunc(a.RequestData(a.UserAuthMux(
			a.HandleDownloadFlashfreezeRootFile,
			muxAny(isStaff, isTrialCurator, isInAudit))))).
		Methods("GET")

	// soft delete
	router.Handle(
		fmt.Sprintf("/api/submission/{%s}/file/{%s}", constants.ResourceKeySubmissionID, constants.ResourceKeyFileID),
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(
			a.HandleSoftDeleteSubmissionFile, muxAll(muxNot(isSubmissionMarkedAsAdded), muxNot(isSubmissionFrozen), isDeleter))))).
		Methods("DELETE")

	router.Handle(
		fmt.Sprintf("/api/submission/{%s}", constants.ResourceKeySubmissionID),
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(
			a.HandleSoftDeleteSubmission, muxAll(muxNot(isSubmissionMarkedAsAdded), muxNot(isSubmissionFrozen), isDeleter))))).
		Methods("DELETE")

	router.Handle(
		fmt.Sprintf("/api/submission/{%s}/comment/{%s}", constants.ResourceKeySubmissionID, constants.ResourceKeyCommentID),
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(
			a.HandleSoftDeleteComment, muxAll(muxNot(isSubmissionMarkedAsAdded), muxNot(isSubmissionFrozen), isDeleter))))).
		Methods("DELETE")

	// bot override

	router.Handle(
		fmt.Sprintf("/api/submission/{%s}/override", constants.ResourceKeySubmissionID),
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(
			a.HandleOverrideBot, muxAny(isDeleter, isStaff))))).
		Methods("POST")

	// freeze

	router.Handle(
		fmt.Sprintf("/api/submission/{%s}/freeze", constants.ResourceKeySubmissionID),
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(
			a.HandleFreezeSubmission, muxAll(muxNot(isSubmissionMarkedAsAdded), isFreezer))))).
		Methods("POST")

	router.Handle(
		fmt.Sprintf("/api/submission/{%s}/unfreeze", constants.ResourceKeySubmissionID),
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(
			a.HandleUnfreezeSubmission, muxAll(muxNot(isSubmissionMarkedAsAdded), isFreezer))))).
		Methods("POST")

	// user statistics

	router.Handle(
		"/api/users",
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(a.HandleGetUsers, muxAny(isStaff, isTrialCurator, isInAudit))))).
		Methods("GET")

	router.Handle(
		fmt.Sprintf("/api/user-statistics/{%s}", constants.ResourceKeyUserID),
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(a.HandleGetUserStatistics, muxAny(isStaff, isTrialCurator, isInAudit))))).
		Methods("GET")

	// upload status
	router.Handle(
		fmt.Sprintf("/api/upload-status/{%s}", constants.ResourceKeyTempName),
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(a.HandleGetUploadProgress, muxAny(isStaff, isTrialCurator, isInAudit))))).
		Methods("GET")

	////////////////////////

	// god tools

	router.Handle("/web/internal",
		http.HandlerFunc(a.RequestWeb(a.UserAuthMux(a.HandleInternalPage, isGod)))).
		Methods("GET")

	router.Handle("/api/internal/update-master-db",
		http.HandlerFunc(a.RequestWeb(a.UserAuthMux(a.HandleUpdateMasterDB, isGod)))).
		Methods("GET")

	router.Handle("/api/internal/flashfreeze/ingest",
		http.HandlerFunc(a.RequestWeb(a.UserAuthMux(a.HandleIngestFlashfreeze, isGod)))).
		Methods("GET")

	router.Handle("/api/internal/recompute-submission-cache-all",
		http.HandlerFunc(a.RequestWeb(a.UserAuthMux(a.HandleRecomputeSubmissionCacheAll, isGod)))).
		Methods("GET")

	router.Handle("/api/internal/flashfreeze/ingest-unknown-files",
		http.HandlerFunc(a.RequestWeb(a.UserAuthMux(a.HandleIngestUnknownFlashfreeze, isGod)))).
		Methods("GET")

	router.Handle("/api/internal/flashfreeze/index-unindexed-files",
		http.HandlerFunc(a.RequestWeb(a.UserAuthMux(a.HandleIndexUnindexedFlashfreeze, isGod)))).
		Methods("GET")

	router.Handle("/api/internal/delete-user-sessions",
		http.HandlerFunc(a.RequestWeb(a.UserAuthMux(a.HandleDeleteUserSessions, isGod)))).
		Methods("POST")

	router.Handle("/api/internal/send-reminders-about-requested-changes",
		http.HandlerFunc(a.RequestWeb(a.UserAuthMux(a.HandleSendRemindersAboutRequestedChanges, isGod)))).
		Methods("GET")

	router.Handle("/api/internal/nuke-session-table",
		http.HandlerFunc(a.RequestWeb(a.UserAuthMux(a.HandleNukeSessionTable, isGod)))).
		Methods("GET")

	err := srv.ListenAndServe()
	if err != nil {
		l.Fatal(err)
	}
}
