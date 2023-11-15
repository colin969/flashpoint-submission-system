package transport

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/FlashpointProject/flashpoint-submission-system/logging"
	"github.com/FlashpointProject/flashpoint-submission-system/service"
	"github.com/FlashpointProject/flashpoint-submission-system/types"
	"github.com/FlashpointProject/flashpoint-submission-system/utils"
	"github.com/gofrs/uuid"
)

type discordUserResponse struct {
	ID            string  `json:"id"`
	Username      string  `json:"username"`
	Avatar        string  `json:"avatar"`
	Discriminator string  `json:"discriminator"`
	PublicFlags   int64   `json:"public_flags"`
	Flags         int64   `json:"flags"`
	Locale        string  `json:"locale"`
	MFAEnabled    bool    `json:"mfa_enabled"`
	GlobalName    *string `json:"global_name,omitempty"`
}

type StateKeeper struct {
	sync.Mutex
	states            map[string]time.Time
	expirationSeconds uint64
}

type State struct {
	Nonce string `json:"nonce"`
	Dest  string `json:"dest"`
}

var clientApps = []types.ClientApplication{
	{
		ClientId: "flashpoint-launcher",
		Name:     "Flashpoint Launcher",
		Scopes:   []string{types.AuthScopeIdentity, types.AuthScopeGameRead, types.AuthScopeGameDataEdit},
	},
}

var authScopes = []types.AuthScope{
	{
		Name:        types.AuthScopeIdentity,
		Description: "Read your username, avatar, Flashpoint discord server roles and FPFSS notification settings",
	},
	{
		Name:        types.AuthScopeSubmissionRead,
		Description: "Read basic submission information (comments, metadata)",
	},
	{
		Name:        types.AuthScopeSubmissionReadFiles,
		Description: "Read and download submission files",
	},
	{
		Name:        types.AuthScopeSubmissionEdit,
		Description: "Edit submission information",
	},
	{
		Name:        types.AuthScopeSubmissionUpload,
		Description: "Upload new submission files",
	},
	{
		Name:        types.AuthScopeFlashfreezeRead,
		Description: "Read basic flashfreeze information (archive metadata)",
	},
	{
		Name:        types.AuthScopeFlashfreezeReadFiles,
		Description: "Read and download flashfreeze files / archive directories",
	},
	{
		Name:        types.AuthScopeFlashfreezeUpload,
		Description: "Upload new flashfreeze files",
	},
	{
		Name:        types.AuthScopeGameDataRead,
		Description: "Read Game Data info / file indexes",
	},
	{
		Name:        types.AuthScopeGameDataEdit,
		Description: "Edit Game Data info",
	},
	{
		Name:        types.AuthScopeGameRead,
		Description: "Read Game metadata",
	},
	{
		Name:        types.AuthScopeGameEdit,
		Description: "Edit Game metadata",
	},
}

// Generate generates state and returns base64-encoded form
func (sk *StateKeeper) Generate(dest string) (string, error) {
	sk.Clean()
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	s := &State{
		Nonce: u.String(),
		Dest:  dest,
	}
	sk.Lock()
	sk.states[s.Nonce] = time.Now()
	sk.Unlock()

	j, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	b := base64.URLEncoding.EncodeToString(j)

	return b, nil
}

// Consume consumes base64-encoded state and returns destination URL
func (sk *StateKeeper) Consume(b string) (string, bool) {
	sk.Clean()
	sk.Lock()
	defer sk.Unlock()

	j, err := base64.URLEncoding.DecodeString(b)
	if err != nil {
		return "", false
	}

	s := &State{}

	err = json.Unmarshal(j, s)
	if err != nil {
		return "", false
	}

	_, ok := sk.states[s.Nonce]
	if ok {
		delete(sk.states, s.Nonce)
	}
	return s.Dest, ok
}

func (sk *StateKeeper) Clean() {
	sk.Lock()
	defer sk.Unlock()
	for k, v := range sk.states {
		if v.After(v.Add(time.Duration(sk.expirationSeconds))) {
			delete(sk.states, k)
		}
	}
}

var stateKeeper = StateKeeper{
	states:            make(map[string]time.Time),
	expirationSeconds: 30,
}

func (a *App) HandleDiscordAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dest := r.FormValue("dest")

	state, err := stateKeeper.Generate(dest)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to generate state", http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, a.Conf.OauthConf.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (a *App) HandleDiscordCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// verify state

	dest, ok := stateKeeper.Consume(r.FormValue("state"))
	if !ok {
		writeError(ctx, w, perr("state does not match", http.StatusBadRequest))
		return
	}

	// obtain token
	token, err := a.Conf.OauthConf.Exchange(context.Background(), r.FormValue("code"))

	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to obtain discord auth token", http.StatusInternalServerError))
		return
	}

	// obtain user data
	resp, err := a.Conf.OauthConf.Client(context.Background(), token).Get("https://discordapp.com/api/users/@me")

	if err != nil || resp.StatusCode != 200 {
		writeError(ctx, w, perr("failed to obtain discord user data", http.StatusInternalServerError))
		return
	}
	defer resp.Body.Close()

	var discordUserResp discordUserResponse
	err = json.NewDecoder(resp.Body).Decode(&discordUserResp)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to parse discord response", http.StatusInternalServerError))
		return
	}

	uid, err := strconv.ParseInt(discordUserResp.ID, 10, 64)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to parse discord response", http.StatusInternalServerError))
		return
	}
	username := discordUserResp.Username
	if discordUserResp.GlobalName != nil && *discordUserResp.GlobalName != "" {
		username = *discordUserResp.GlobalName
	}

	discordUser := &types.DiscordUser{
		ID:            uid,
		Username:      username,
		Avatar:        discordUserResp.Avatar,
		Discriminator: discordUserResp.Discriminator,
		PublicFlags:   discordUserResp.PublicFlags,
		Flags:         discordUserResp.Flags,
		Locale:        discordUserResp.Locale,
		MFAEnabled:    discordUserResp.MFAEnabled,
	}

	ipAddr := logging.RequestGetRemoteAddress(r)
	authToken, err := a.Service.SaveUser(ctx, discordUser, ipAddr)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, dberr(err))
		return
	}

	if err := a.CC.SetSecureCookie(w, utils.Cookies.Login, service.MapAuthToken(authToken), (int)(a.Conf.SessionExpirationSeconds)); err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to set cookie", http.StatusInternalServerError))
		return
	}

	if len(dest) == 0 || !isReturnURLValid(dest) {
		http.Redirect(w, r, "/web/profile", http.StatusFound)
		return
	}

	http.Redirect(w, r, dest, http.StatusFound)
}

func (a *App) HandleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const msg = "unable to log out, please clear your cookies"
	cookieMap, err := a.CC.GetSecureCookie(r, utils.Cookies.Login)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr(msg, http.StatusInternalServerError))
		return
	}

	token, err := service.ParseAuthToken(cookieMap) // TODO move this into the Logout method
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr(msg, http.StatusInternalServerError))
		return
	}

	if err := a.Service.Logout(ctx, token.Secret); err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr(msg, http.StatusInternalServerError))
		return
	}

	utils.UnsetCookie(w, utils.Cookies.Login)
	http.Redirect(w, r, "/web", http.StatusFound)
}

func (a *App) HandlePollDeviceAuth(w http.ResponseWriter, ctx context.Context, deviceCode string) {
	// Get device auth token from storage
	dfToken := a.DFStorage.GetUserAuthTokenByDevice(deviceCode)
	if dfToken == nil {
		writeError(ctx, w, perr("no tokens found", http.StatusBadRequest))
		return
	}

	switch dfToken.FlowState {
	case types.DeviceFlowComplete:
		if dfToken.AuthToken == nil {
			writeError(ctx, w, perr("device auth complete but no token found.", http.StatusInternalServerError))
			return
		}
		// Encode the auth token
		authJson, err := json.Marshal(dfToken.AuthToken)
		if err != nil {
			writeError(ctx, w, perr("failure marshalling token", http.StatusInternalServerError))
			return
		}
		encodedData := base64.StdEncoding.EncodeToString(authJson)
		jsonData := types.DeviceFlowPollResponse{
			Token: encodedData,
		}
		writeResponse(ctx, w, jsonData, http.StatusOK)
		return
	case types.DeviceFlowPending:
		jsonData := types.DeviceFlowPollResponse{
			Error: "authorization_pending",
		}
		writeResponse(ctx, w, jsonData, http.StatusOK)
		return
	case types.DeviceFlowErrorDenied:
		jsonData := types.DeviceFlowPollResponse{
			Error: "access_denied",
		}
		writeResponse(ctx, w, jsonData, http.StatusOK)
		return
	case types.DeviceFlowErrorExpired:
		jsonData := types.DeviceFlowPollResponse{
			Error: "expired_token",
		}
		writeResponse(ctx, w, jsonData, http.StatusOK)
		return
	}

}

func (a *App) HandleOauthToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to parse form", http.StatusBadRequest))
		return
	}

	// Validate client application
	client_id := r.Form.Get("client_id")
	if client_id == "" {
		writeError(ctx, w, perr("missing client_id", http.StatusBadRequest))
		return
	}
	var client *types.ClientApplication
	for _, app := range clientApps {
		if app.ClientId == client_id {
			client = &app
			break
		}
	}
	if client == nil {
		writeError(ctx, w, perr("invalid client_id", http.StatusBadRequest))
		return
	}

	deviceCode := r.Form.Get("device_code")
	if deviceCode == "" {
		// No device code given, must be requesting a new one
		// Read list of scopes, handle default case
		scope := r.Form.Get("scope")
		if scope != "" {
			// Filter out invalid scopes
			var validScopes []string
			for _, scopeStr := range strings.Split(scope, " ") {
				for _, allowedClientScope := range client.Scopes {
					if scopeStr == allowedClientScope {
						validScopes = append(validScopes, scopeStr)
					}
				}
			}
			scope = strings.Join(validScopes, " ")
			if scope == "" {
				// No valid scopes found, but scope given, give advice
				scopeNames := make([]string, len(authScopes))
				for i, authScope := range authScopes {
					scopeNames[i] = authScope.Name
				}
				writeError(ctx, w, perr("invalid scope: Must be of ["+strings.Join(scopeNames, ", ")+"]", http.StatusBadRequest))
				return
			}
		} else {
			scope = types.AuthScopeIdentity
		}
		token, err := a.DFStorage.NewToken(scope, client)
		if err != nil {
			utils.LogCtx(ctx).Error(err)
			writeError(ctx, w, perr("failed to create token", http.StatusInternalServerError))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		writeResponse(ctx, w, token, http.StatusOK)
	} else {
		// Device code given, must be polling
		a.HandlePollDeviceAuth(w, ctx, deviceCode)
		return
	}
}

func (a *App) HandleApproveDevice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()
	code := query.Get("user_code")

	// Get device auth token from storage
	token, err := a.DFStorage.Get(code)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr(err.Error(), http.StatusBadRequest))
		return
	}

	if r.Method == http.MethodPost {
		// POST User has responded
		action := query.Get("action")
		if action == "approve" {
			// Create a new auth token
			uid := utils.UserID(ctx)
			ipAddr := logging.RequestGetRemoteAddress(r)
			authToken, err := a.Service.GenAuthToken(ctx, uid, token.Scope, token.ClientApplication.ClientId, ipAddr)
			if err != nil {
				utils.LogCtx(ctx).Error(err)
				writeError(ctx, w, perr("failed to create new auth token", http.StatusInternalServerError))
				return
			}

			// Save inside device auth
			token.FlowState = types.DeviceFlowComplete
			token.AuthToken = authToken
			err = a.DFStorage.Save(token)
			if err != nil {
				utils.LogCtx(ctx).Error(err)
				writeError(ctx, w, perr("failed to save device token", http.StatusInternalServerError))
				return
			}
		} else if action == "deny" {
			token.FlowState = types.DeviceFlowErrorDenied
			err := a.DFStorage.Save(token)
			if err != nil {
				utils.LogCtx(ctx).Error(err)
				writeError(ctx, w, perr("failed to save token", http.StatusInternalServerError))
				return
			}
		} else {
			writeError(ctx, w, perr("invalid action, must be 'approve' or 'deny'", http.StatusBadRequest))
			return
		}
		// POST Action complete continue to show result same as GET
	}

	// GET Ask for user response
	// Load scopes
	var scopesList []types.AuthScope
	for _, scope := range strings.Split(token.Scope, " ") {
		for _, authScope := range authScopes {
			if scope == authScope.Name {
				scopesList = append(scopesList, authScope)
			}
		}
	}
	bpd, err := a.Service.GetBasePageData(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failure getting page data", http.StatusInternalServerError))
		return
	}
	var states = types.DeviceAuthStates{
		Pending:  types.DeviceFlowPending,
		Denied:   types.DeviceFlowErrorDenied,
		Expired:  types.DeviceFlowErrorExpired,
		Complete: types.DeviceFlowComplete,
	}
	pageData := types.DeviceAuthPageData{
		BasePageData: *bpd,
		Token:        token,
		States:       states,
		Scopes:       scopesList,
	}

	a.RenderTemplates(ctx, w, r, pageData,
		"templates/device_auth.gohtml")
}

const deviceCodeCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const userCodeCharset = "BCDFGHJKLMNPQRSTVWXZ"

type DeviceFlowUserAuthToken struct {
	AuthToken  string
	DeviceCode string
}

type DeviceFlowStorage struct {
	tokens          map[string]*types.DeviceFlowToken
	authTokens      map[int64]*[]DeviceFlowUserAuthToken
	verificationUrl string
}

func NewDeviceFlowStorage(verificationUrl string) *DeviceFlowStorage {
	return &DeviceFlowStorage{
		tokens:          make(map[string]*types.DeviceFlowToken),
		authTokens:      make(map[int64]*[]DeviceFlowUserAuthToken),
		verificationUrl: verificationUrl,
	}
}

func (s *DeviceFlowStorage) GetUserAuthTokenByDevice(deviceCode string) *types.DeviceFlowToken {
	var dfToken *types.DeviceFlowToken
	for _, token := range s.tokens {
		if token.DeviceCode == deviceCode {
			dfToken = token
		}
	}

	return dfToken
}

func (s *DeviceFlowStorage) SaveUserAuthToken(uid int64, token string, deviceCode string) {
	userToken := DeviceFlowUserAuthToken{
		AuthToken:  token,
		DeviceCode: deviceCode,
	}
	*s.authTokens[uid] = append(*s.authTokens[uid], userToken)
}

func (s *DeviceFlowStorage) GetUserAuthTokens(uid int64) *[]DeviceFlowUserAuthToken {
	return s.authTokens[uid]
}

func (s *DeviceFlowStorage) NewToken(scope string, client *types.ClientApplication) (*types.DeviceFlowToken, error) {
	// Generate the code
	deviceCode := make([]byte, 32)
	for i := range deviceCode {
		deviceCode[i] = deviceCodeCharset[rand.Intn(len(deviceCodeCharset))]
	}

	userCode := make([]byte, 32)
	for i := range userCode {
		userCode[i] = userCodeCharset[rand.Intn(len(userCodeCharset))]
	}

	expiresAt := time.Now()
	expiresAt = expiresAt.Add(900 * time.Second)

	token := types.DeviceFlowToken{
		DeviceCode:              string(deviceCode),
		Scope:                   scope,
		UserCode:                string(userCode),
		VerificationURL:         s.verificationUrl,
		VerificationURLComplete: s.verificationUrl + "?user_code=" + string(userCode),
		ExpiresIn:               900,
		ExpiresAt:               expiresAt,
		Interval:                3,
		FlowState:               types.DeviceFlowPending,
		ClientApplication:       client,
	}

	err := s.Save(&token)
	if err != nil {
		return &token, err
	}

	return &token, nil
}

func (s *DeviceFlowStorage) Save(token *types.DeviceFlowToken) error {
	s.tokens[token.UserCode] = token
	return nil
}

func (s *DeviceFlowStorage) Get(userCode string) (*types.DeviceFlowToken, error) {
	token, found := s.tokens[userCode]
	if !found {
		return nil, errors.New("device code not found")
	}
	if time.Now().After(token.ExpiresAt) {
		return nil, errors.New("device code has expired")
	}
	return token, nil
}

func (s *DeviceFlowStorage) Delete(deviceCode string) {
	delete(s.tokens, deviceCode)
}

func (s *DeviceFlowStorage) Cleanup() {
	for deviceCode, token := range s.tokens {
		if time.Now().After(token.ExpiresAt) {
			s.Delete(deviceCode)
		}
	}
}
