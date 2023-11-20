package clients

import "github.com/FlashpointProject/flashpoint-submission-system/types"

var BotUserMinID = int64(1000)

var ClientApps = []types.ClientApplication{
	{
		UserID:            BotUserMinID + 2,
		UserRoles:         []string{},
		ClientId:          "flashpoint-launcher",
		Name:              "Flashpoint Launcher",
		Scopes:            []string{types.AuthScopeIdentity, types.AuthScopeGameRead, types.AuthScopeGameEdit},
		ClientCredsScopes: []string{},
		OwnerUID:          int64(689080719460663414), // Colin
	},
	{
		UserID:            BotUserMinID + 3,
		UserRoles:         []string{"Curator"},
		ClientId:          "flashpoint-community",
		Name:              "Flashpoint Community",
		Scopes:            []string{types.AuthScopeIdentity},
		RedirectURIs:      []string{"https://fpcomm-dev.colintest.site/auth/callback"},
		ClientCredsScopes: []string{types.AuthScopeIdentity, types.AuthScopeGameRead},
		OwnerUID:          int64(689080719460663414), // Colin
	},
}
