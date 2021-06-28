package transport

import (
	"context"
	"fmt"
	"github.com/Dri0m/flashpoint-submission-system/constants"
	"github.com/Dri0m/flashpoint-submission-system/service"
	"github.com/Dri0m/flashpoint-submission-system/types"
	"github.com/Dri0m/flashpoint-submission-system/utils"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func (a *App) HandleCommentReceiverBatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := utils.UserID(ctx)

	params := mux.Vars(r)
	submissionIDs := strings.Split(params["submission-ids"], ",")
	sids := make([]int64, 0, len(submissionIDs))

	for _, submissionFileID := range submissionIDs {
		sid, err := strconv.ParseInt(submissionFileID, 10, 64)
		if err != nil {
			utils.LogCtx(ctx).Error(err)
			writeError(ctx, w, perr("invalid submission id", http.StatusBadRequest))
			return
		}
		sids = append(sids, sid)
	}

	if err := r.ParseForm(); err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to parse form", http.StatusBadRequest))
		return
	}

	// TODO use gorilla/schema
	formAction := r.FormValue("action")
	formMessage := r.FormValue("message")
	formIgnoreDupeActions := r.FormValue("ignore-duplicate-actions")

	if len([]rune(formMessage)) > 20000 {
		err := fmt.Errorf("message cannot be longer than 20000 characters")
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, constants.PublicError{Msg: err.Error(), Status: http.StatusBadRequest})
		return
	}

	if err := a.Service.ReceiveComments(ctx, uid, sids, formAction, formMessage, formIgnoreDupeActions); err != nil {
		writeError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *App) HandleDownloadSubmissionFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	submissionFileID := params[constants.ResourceKeyFileID]

	sfid, err := strconv.ParseInt(submissionFileID, 10, 64)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("invalid submission file id", http.StatusBadRequest))
		return
	}

	sfs, err := a.Service.GetSubmissionFiles(ctx, []int64{sfid})
	if err != nil {
		writeError(ctx, w, err)
		return
	}
	sf := sfs[0]

	f, err := os.Open(fmt.Sprintf("%s/%s", constants.SubmissionsDir, sf.CurrentFilename))

	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to read file", http.StatusInternalServerError))
		return
	}
	defer f.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", sf.CurrentFilename))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(w, r, sf.CurrentFilename, sf.UploadedAt, f)
}

func (a *App) HandleSoftDeleteSubmissionFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	submissionFileID := params[constants.ResourceKeyFileID]

	sfid, err := strconv.ParseInt(submissionFileID, 10, 64)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("invalid submission file id", http.StatusBadRequest))
		return
	}

	if err := r.ParseForm(); err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to parse form", http.StatusBadRequest))
		return
	}

	deleteReason := r.FormValue("reason")
	if len(deleteReason) < 3 {
		writeError(ctx, w, perr("reason must be at least 3 characters long", http.StatusBadRequest))
		return
	} else if len(deleteReason) > 255 {
		writeError(ctx, w, perr("reason cannot be longer than 255 characters", http.StatusBadRequest))
		return
	}

	if err := a.Service.SoftDeleteSubmissionFile(ctx, sfid, deleteReason); err != nil {
		writeError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *App) HandleSoftDeleteSubmission(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	submissionID := params[constants.ResourceKeySubmissionID]

	sid, err := strconv.ParseInt(submissionID, 10, 64)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("invalid submission id", http.StatusBadRequest))
		return
	}

	deleteReason := r.FormValue("reason")
	if len(deleteReason) < 3 {
		writeError(ctx, w, perr("reason must be at least 3 characters long", http.StatusBadRequest))
		return
	} else if len(deleteReason) > 255 {
		writeError(ctx, w, perr("reason cannot be longer than 255 characters", http.StatusBadRequest))
		return
	}

	if err := a.Service.SoftDeleteSubmission(ctx, sid, deleteReason); err != nil {
		writeError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *App) HandleSoftDeleteComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	commentID := params[constants.ResourceKeyCommentID]

	cid, err := strconv.ParseInt(commentID, 10, 64)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("invalid comment id", http.StatusBadRequest))
		return
	}

	deleteReason := r.FormValue("reason")
	if len(deleteReason) < 3 {
		writeError(ctx, w, perr("reason must be at least 3 characters long", http.StatusBadRequest))
		return
	} else if len(deleteReason) > 255 {
		writeError(ctx, w, perr("reason cannot be longer than 255 characters", http.StatusBadRequest))
		return
	}

	if err := a.Service.SoftDeleteComment(ctx, cid, deleteReason); err != nil {
		writeError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *App) HandleDownloadSubmissionBatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	submissionFileIDs := strings.Split(params[constants.ResourceKeyFileIDs], ",")
	sfids := make([]int64, 0, len(submissionFileIDs))

	for _, submissionFileID := range submissionFileIDs {
		sfid, err := strconv.ParseInt(submissionFileID, 10, 64)
		if err != nil {
			utils.LogCtx(ctx).Error(err)
			writeError(ctx, w, perr("invalid submission file id", http.StatusBadRequest))
			return
		}
		sfids = append(sfids, sfid)
	}

	sfs, err := a.Service.GetSubmissionFiles(ctx, sfids)
	if err != nil {
		writeError(ctx, w, err)
		return
	}

	filePaths := make([]string, 0, len(sfs))

	for _, sf := range sfs {
		filePaths = append(filePaths, fmt.Sprintf("%s/%s", constants.SubmissionsDir, sf.CurrentFilename))
	}

	filename := fmt.Sprintf("fpfss-batch-%dfiles-%s.tar", len(sfs), utils.NewRealRandomStringProvider().RandomString(16))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	if err := utils.WriteTarball(w, filePaths); err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to create tarball", http.StatusInternalServerError))
		return
	}
}

func (a *App) HandleSubmissionReceiver(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	submissionID := params[constants.ResourceKeySubmissionID]

	var sid *int64

	if submissionID != "" {
		sidParsed, err := strconv.ParseInt(submissionID, 10, 64)
		if err != nil {
			utils.LogCtx(ctx).Error(err)
			writeError(ctx, w, perr("invalid submission id", http.StatusBadRequest))
			return
		}
		sid = &sidParsed
	}

	// limit RAM usage to 10MB
	if err := r.ParseMultipartForm(10 * 1000 * 1000); err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to parse form", http.StatusInternalServerError))
		return
	}

	fileHeaders := r.MultipartForm.File["files"]

	if len(fileHeaders) == 0 {
		err := fmt.Errorf("no files received")
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, constants.PublicError{Msg: err.Error(), Status: http.StatusBadRequest})
		return
	}

	fileWrappers := make([]service.MultipartFileProvider, 0, len(fileHeaders))
	for _, fileHeader := range fileHeaders {
		fileWrappers = append(fileWrappers, service.NewMutlipartFileWrapper(fileHeader))
	}

	if err := a.Service.ReceiveSubmissions(ctx, sid, fileWrappers); err != nil {
		writeError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *App) HandleRootPage(w http.ResponseWriter, r *http.Request) {
	uid, err := a.GetUserIDFromCookie(r)
	ctx := r.Context()
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, err)
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), utils.CtxKeys.UserID, uid))
	ctx = r.Context()

	pageData, err := a.Service.GetBasePageData(ctx)
	if err != nil {
		utils.UnsetCookie(w, utils.Cookies.Login)
	}

	a.RenderTemplates(ctx, w, r, pageData, "templates/root.gohtml")
}

func (a *App) HandleProfilePage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := utils.UserID(ctx)

	pageData, err := a.Service.GetProfilePageData(ctx, uid)
	if err != nil {
		writeError(ctx, w, err)
		return
	}

	a.RenderTemplates(ctx, w, r, pageData, "templates/profile.gohtml")
}

func (a *App) HandleSubmitPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pageData, err := a.Service.GetBasePageData(ctx)
	if err != nil {
		writeError(ctx, w, err)
		return
	}

	a.RenderTemplates(ctx, w, r, pageData, "templates/submit.gohtml")
}

func (a *App) HandleSubmissionsPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter := &types.SubmissionsFilter{}

	if err := a.decoder.Decode(filter, r.URL.Query()); err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to decode query params", http.StatusInternalServerError))
		return
	}

	if err := filter.Validate(); err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, err)
		return
	}

	pageData, err := a.Service.GetSubmissionsPageData(ctx, filter)
	if err != nil {
		writeError(ctx, w, err)
		return
	}

	pageData.FilterLayout = r.FormValue("filter-layout")

	a.RenderTemplates(ctx, w, r, pageData,
		"templates/submissions.gohtml",
		"templates/submission-filter.gohtml",
		"templates/submission-table.gohtml",
		"templates/submission-pagenav.gohtml",
		"templates/submission-filter-chunks.gohtml",
		"templates/comment-form.gohtml")
}

func (a *App) HandleMySubmissionsPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := utils.UserID(ctx)

	filter := &types.SubmissionsFilter{}

	if err := a.decoder.Decode(filter, r.URL.Query()); err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to decode query params", http.StatusInternalServerError))
		return
	}

	if err := filter.Validate(); err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, err)
		return
	}

	filter.SubmitterID = &uid

	pageData, err := a.Service.GetSubmissionsPageData(ctx, filter)
	if err != nil {
		writeError(ctx, w, err)
		return
	}

	pageData.FilterLayout = r.FormValue("filter-layout")

	a.RenderTemplates(ctx, w, r, pageData,
		"templates/my-submissions.gohtml",
		"templates/submission-filter.gohtml",
		"templates/submission-table.gohtml",
		"templates/submission-pagenav.gohtml",
		"templates/submission-filter-chunks.gohtml",
		"templates/comment-form.gohtml")
}

func (a *App) HandleViewSubmissionPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := utils.UserID(ctx)
	params := mux.Vars(r)
	submissionID := params[constants.ResourceKeySubmissionID]

	sid, err := strconv.ParseInt(submissionID, 10, 64)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("invalid submission id", http.StatusBadRequest))
		return
	}

	pageData, err := a.Service.GetViewSubmissionPageData(ctx, uid, sid)
	if err != nil {
		writeError(ctx, w, err)
		return
	}

	a.RenderTemplates(ctx, w, r, pageData,
		"templates/submission.gohtml",
		"templates/submission-table.gohtml",
		"templates/comment-form.gohtml",
		"templates/view-submission-nav.gohtml")
}

func (a *App) HandleViewSubmissionFilesPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	submissionID := params[constants.ResourceKeySubmissionID]

	sid, err := strconv.ParseInt(submissionID, 10, 64)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("invalid submission id", http.StatusBadRequest))
		return
	}

	pageData, err := a.Service.GetSubmissionsFilesPageData(ctx, sid)
	if err != nil {
		writeError(ctx, w, err)
		return
	}

	a.RenderTemplates(ctx, w, r, pageData, "templates/submission-files.gohtml", "templates/submission-files-table.gohtml")
}

func (a *App) HandleUpdateNotificationSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := utils.UserID(ctx)

	notificationSettings := &types.UpdateNotificationSettings{}

	if err := a.decoder.Decode(notificationSettings, r.URL.Query()); err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to decode query params", http.StatusInternalServerError))
		return
	}

	err := a.Service.UpdateNotificationSettings(ctx, uid, notificationSettings.NotificationActions)
	if err != nil {
		writeError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *App) HandleUpdateSubscriptionSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := utils.UserID(ctx)
	params := mux.Vars(r)
	submissionID := params[constants.ResourceKeySubmissionID]

	sid, err := strconv.ParseInt(submissionID, 10, 64)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("invalid submission id", http.StatusBadRequest))
		return
	}

	subscriptionSettings := &types.UpdateSubscriptionSettings{}

	if err := a.decoder.Decode(subscriptionSettings, r.URL.Query()); err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to decode query params", http.StatusInternalServerError))
		return
	}

	err = a.Service.UpdateSubscriptionSettings(ctx, uid, sid, subscriptionSettings.Subscribe)
	if err != nil {
		writeError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *App) HandleDownloadCurationImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	curationImageID := params[constants.ResourceKeyCurationImageID]

	ciid, err := strconv.ParseInt(curationImageID, 10, 64)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("invalid curation image id", http.StatusBadRequest))
		return
	}

	ci, err := a.Service.GetCurationImage(ctx, ciid)
	if err != nil {
		writeError(ctx, w, err)
		return
	}

	f, err := os.Open(fmt.Sprintf("%s/%s", constants.SubmissionImagesDir, ci.Filename))
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to read file", http.StatusInternalServerError))
		return
	}

	fi, err := f.Stat()
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to read file", http.StatusInternalServerError))
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "image")
	http.ServeContent(w, r, ci.Filename, fi.ModTime(), f)
}
