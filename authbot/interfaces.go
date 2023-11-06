package authbot

import "github.com/FlashpointProject/flashpoint-submission-system/types"

type DiscordRoleReader interface {
	GetFlashpointRoleIDsForUser(uid int64) ([]string, error)
	GetFlashpointRoles() ([]types.DiscordRole, error)
}
