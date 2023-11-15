package types

import "time"

const (
	DeviceFlowPending      = 0
	DeviceFlowErrorExpired = 1
	DeviceFlowErrorDenied  = 2
	DeviceFlowComplete     = 3
)

type DeviceFlowToken struct {
	DeviceCode              string             `json:"device_code"`
	Scope                   string             `json:"-"`
	UserCode                string             `json:"user_code"`
	VerificationURL         string             `json:"verification_uri"`
	VerificationURLComplete string             `json:"verification_uri_complete"`
	ExpiresIn               int64              `json:"expires_in"`
	Interval                int64              `json:"interval"`
	ClientApplication       *ClientApplication `json:"-"`
	ExpiresAt               time.Time          `json:"-"`
	FlowState               int64              `json:"-"`
	AuthToken               map[string]string  `json:"-"` // Explictly add this to responses when suitable
}

type DeviceFlowPollResponse struct {
	Error string `json:"error,omitempty"`
	Token string `json:"access_token,omitempty"`
}

type AuthScope struct {
	Name        string
	Description string
}

const (
	AuthScopeNone                 = ""
	AuthScopeAll                  = "all"
	AuthScopeIdentity             = "identity"
	AuthScopeProfileEdit          = "profile:edit"
	AuthScopeUsersRead            = "users:read"
	AuthScopeSubmissionRead       = "submission:read"
	AuthScopeSubmissionReadFiles  = "submission:read-files"
	AuthScopeSubmissionEdit       = "submission:edit"
	AuthScopeSubmissionUpload     = "submission:upload"
	AuthScopeFlashfreezeRead      = "flashfreeze:read"
	AuthScopeFlashfreezeReadFiles = "flashfreeze:read-files"
	AuthScopeFlashfreezeUpload    = "flashfreeze:upload"
	AuthScopeTagEdit              = "tag:edit"
	AuthScopeGameDataRead         = "game-data:read"
	AuthScopeGameDataEdit         = "game-data:edit"
	AuthScopeGameRead             = "game:read"
	AuthScopeGameEdit             = "game:edit"
	AuthScopeHashCheck            = "hash-check"
)

type ClientApplication struct {
	ClientId string   `json:"client_id"`
	Name     string   `json:"name"`
	Scopes   []string `json:"scopes"`
}

type SessionInfo struct {
	ID        int64  `json:"id"`
	UID       int64  `json:"uid"`
	Scope     string `json:"scope"`
	Client    string `json:"client"`
	ExpiresAt int64  `json:"expires_at"`
	IpAddr    string `json:"ip_addr"`
}
