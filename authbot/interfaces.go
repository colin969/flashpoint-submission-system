package authbot

import (
	"github.com/FlashpointProject/flashpoint-submission-system/types"
	"time"
)

type DiscordRoleReader interface {
	GetFlashpointRoleIDsForUser(uid int64) ([]string, error)
	GetFlashpointRoles() ([]types.DiscordRole, error)
	GetFlashpointUserInfo(uid int64, roles []types.DiscordRole) (*types.FlashpointDiscordUser, error)
	GetJoinedAtForUser(uid int64) (time.Time, error)
}
